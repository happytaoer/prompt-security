package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/happytaoer/prompt-security/internal/config"
	"github.com/happytaoer/prompt-security/internal/db"
	"github.com/happytaoer/prompt-security/internal/filter"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the web server
type Server struct {
	configManager *config.Manager
	logger        *slog.Logger
}

// NewServer creates a new web server instance
func NewServer(manager *config.Manager) *Server {
	return &Server{
		configManager: manager,
		logger:        slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

// AddLog adds a new log entry to the database
func (s *Server) AddLog(originalText, filteredText string, replacements []filter.ReplacementInfo) {
	// Build detections list
	detections := make([]string, 0)
	for _, r := range replacements {
		detections = append(detections, r.Type)
	}

	// Add to database
	if err := db.AddLog(originalText, filteredText, detections); err != nil {
		s.logger.Error("Failed to add log to database", "error", err)
	}
}

// GetConfig returns a copy of the current configuration
func (s *Server) GetConfig() config.Config {
	return s.configManager.Get()
}

// UpdateConfig updates the configuration and notifies all listeners
func (s *Server) UpdateConfig(cfg config.Config) error {
	return s.configManager.Update(cfg)
}

// Start starts the web server
func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()

	// Create a sub-filesystem rooted at the static directory so that
	// visiting http://localhost:8181 serves static/index.html directly.
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to create static filesystem: %w", err)
	}

	// Serve static files from the root path.
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API endpoints
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/logs", s.handleLogs)
	mux.HandleFunc("/api/logs/clear", s.handleClearLogs)

	s.logger.Info("Starting web server", "address", addr)
	fmt.Printf("\n🌐 Web UI available at: http://%s\n\n", addr)

	return http.ListenAndServe(addr, s.corsMiddleware(mux))
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleConfig handles configuration GET and POST requests
func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		cfg := s.GetConfig()
		json.NewEncoder(w).Encode(cfg)

	case http.MethodPost:
		var cfg config.Config
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.UpdateConfig(cfg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleLogs handles log retrieval from database with pagination
func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse pagination parameters
	query := r.URL.Query()
	page := 1
	pageSize := 20

	if pageStr := query.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := query.Get("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	// Get logs from database with pagination
	logs, err := db.GetLogsWithPagination(page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get logs from database", "error", err)
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}

	// Get total count
	totalCount, err := db.GetLogCount()
	if err != nil {
		s.logger.Error("Failed to get log count", "error", err)
		totalCount = 0
	}

	// Calculate total pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	// Prepare response
	response := map[string]interface{}{
		"logs":       logs,
		"page":       page,
		"pageSize":   pageSize,
		"totalCount": totalCount,
		"totalPages": totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleClearLogs handles clearing all logs from database
func (s *Server) handleClearLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear logs from database
	if err := db.ClearLogs(); err != nil {
		s.logger.Error("Failed to clear logs from database", "error", err)
		http.Error(w, "Failed to clear logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
