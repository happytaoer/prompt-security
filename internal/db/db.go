package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Initialize initializes the database connection and creates tables if needed
func Initialize() error {
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}

	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	db = database

	// Create tables
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return db
}

// getDBPath returns the path to the SQLite database file
func getDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".prompt-security")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return filepath.Join(configDir, "config.db"), nil
}

// createTables creates the necessary database tables
func createTables() error {
	// Configuration table - stores a single row with all config
	configTableSQL := `
	CREATE TABLE IF NOT EXISTS config (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		detect_emails INTEGER NOT NULL DEFAULT 1,
		detect_phones INTEGER NOT NULL DEFAULT 1,
		detect_credit_cards INTEGER NOT NULL DEFAULT 1,
		detect_ssns INTEGER NOT NULL DEFAULT 1,
		detect_ipv4 INTEGER NOT NULL DEFAULT 1,
		
		custom_email_pattern TEXT DEFAULT '',
		custom_phone_pattern TEXT DEFAULT '',
		custom_credit_card_pattern TEXT DEFAULT '',
		custom_ssn_pattern TEXT DEFAULT '',
		custom_ipv4_pattern TEXT DEFAULT '',
		
		email_replacement TEXT DEFAULT 'security@example.com',
		phone_replacement TEXT DEFAULT '+1-555-123-4567',
		credit_card_replacement TEXT DEFAULT 'XXXX-XXXX-XXXX-XXXX',
		ssn_replacement TEXT DEFAULT 'XXX-XX-XXXX',
		ipv4_replacement TEXT DEFAULT '0.0.0.0',
		
		monitoring_interval_ms INTEGER NOT NULL DEFAULT 500,
		notify_on_filter INTEGER NOT NULL DEFAULT 1,
		
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(configTableSQL); err != nil {
		return fmt.Errorf("failed to create config table: %v", err)
	}

	// String match patterns table
	stringPatternsTableSQL := `
	CREATE TABLE IF NOT EXISTS string_match_patterns (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		pattern TEXT NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 1,
		replacement TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(stringPatternsTableSQL); err != nil {
		return fmt.Errorf("failed to create string_match_patterns table: %v", err)
	}

	// Logs table
	logsTableSQL := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		original_text TEXT NOT NULL,
		filtered_text TEXT NOT NULL,
		detections TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp DESC);`

	if _, err := db.Exec(logsTableSQL); err != nil {
		return fmt.Errorf("failed to create logs table: %v", err)
	}

	// Insert default config if not exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM config").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check config existence: %v", err)
	}

	if count == 0 {
		insertDefaultSQL := `
		INSERT INTO config (id) VALUES (1);`
		if _, err := db.Exec(insertDefaultSQL); err != nil {
			return fmt.Errorf("failed to insert default config: %v", err)
		}
	}

	return nil
}

// StringMatchPattern represents a string match pattern
type StringMatchPattern struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Enabled     bool   `json:"enabled"`
	Replacement string `json:"replacement"`
}

// Config represents the application configuration
type Config struct {
	DetectEmails      bool `json:"detect_emails"`
	DetectPhones      bool `json:"detect_phones"`
	DetectCreditCards bool `json:"detect_credit_cards"`
	DetectSSNs        bool `json:"detect_ssns"`
	DetectIPV4        bool `json:"detect_ipv4"`

	StringMatchPatterns []StringMatchPattern `json:"string_match_patterns"`

	CustomEmailPattern      string `json:"custom_email_pattern"`
	CustomPhonePattern      string `json:"custom_phone_pattern"`
	CustomCreditCardPattern string `json:"custom_credit_card_pattern"`
	CustomSSNPattern        string `json:"custom_ssn_pattern"`
	CustomIPV4Pattern       string `json:"custom_ipv4_pattern"`

	EmailReplacement      string `json:"email_replacement"`
	PhoneReplacement      string `json:"phone_replacement"`
	CreditCardReplacement string `json:"credit_card_replacement"`
	SSNReplacement        string `json:"ssn_replacement"`
	IPV4Replacement       string `json:"ipv4_replacement"`

	MonitoringInterval int  `json:"monitoring_interval_ms"`
	NotifyOnFilter     bool `json:"notify_on_filter"`
}

// LoadConfig loads the configuration from the database
func LoadConfig() (Config, error) {
	var cfg Config

	row := db.QueryRow(`
		SELECT 
			detect_emails, detect_phones, detect_credit_cards, detect_ssns, detect_ipv4,
			custom_email_pattern, custom_phone_pattern, custom_credit_card_pattern, 
			custom_ssn_pattern, custom_ipv4_pattern,
			email_replacement, phone_replacement, credit_card_replacement, 
			ssn_replacement, ipv4_replacement,
			monitoring_interval_ms, notify_on_filter
		FROM config WHERE id = 1
	`)

	var detectEmails, detectPhones, detectCreditCards, detectSSNs, detectIPV4 int
	var notifyOnFilter int

	err := row.Scan(
		&detectEmails, &detectPhones, &detectCreditCards, &detectSSNs, &detectIPV4,
		&cfg.CustomEmailPattern, &cfg.CustomPhonePattern, &cfg.CustomCreditCardPattern,
		&cfg.CustomSSNPattern, &cfg.CustomIPV4Pattern,
		&cfg.EmailReplacement, &cfg.PhoneReplacement, &cfg.CreditCardReplacement,
		&cfg.SSNReplacement, &cfg.IPV4Replacement,
		&cfg.MonitoringInterval, &notifyOnFilter,
	)

	if err != nil {
		return cfg, fmt.Errorf("failed to load config: %v", err)
	}

	// Convert integers to booleans
	cfg.DetectEmails = detectEmails == 1
	cfg.DetectPhones = detectPhones == 1
	cfg.DetectCreditCards = detectCreditCards == 1
	cfg.DetectSSNs = detectSSNs == 1
	cfg.DetectIPV4 = detectIPV4 == 1
	cfg.NotifyOnFilter = notifyOnFilter == 1

	// Load string match patterns
	patterns, err := LoadStringMatchPatterns()
	if err != nil {
		return cfg, fmt.Errorf("failed to load string match patterns: %v", err)
	}
	cfg.StringMatchPatterns = patterns

	return cfg, nil
}

// SaveConfig saves the configuration to the database
func SaveConfig(cfg Config) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Convert booleans to integers
	detectEmails := 0
	if cfg.DetectEmails {
		detectEmails = 1
	}
	detectPhones := 0
	if cfg.DetectPhones {
		detectPhones = 1
	}
	detectCreditCards := 0
	if cfg.DetectCreditCards {
		detectCreditCards = 1
	}
	detectSSNs := 0
	if cfg.DetectSSNs {
		detectSSNs = 1
	}
	detectIPV4 := 0
	if cfg.DetectIPV4 {
		detectIPV4 = 1
	}
	notifyOnFilter := 0
	if cfg.NotifyOnFilter {
		notifyOnFilter = 1
	}

	_, err = tx.Exec(`
		UPDATE config SET
			detect_emails = ?,
			detect_phones = ?,
			detect_credit_cards = ?,
			detect_ssns = ?,
			detect_ipv4 = ?,
			custom_email_pattern = ?,
			custom_phone_pattern = ?,
			custom_credit_card_pattern = ?,
			custom_ssn_pattern = ?,
			custom_ipv4_pattern = ?,
			email_replacement = ?,
			phone_replacement = ?,
			credit_card_replacement = ?,
			ssn_replacement = ?,
			ipv4_replacement = ?,
			monitoring_interval_ms = ?,
			notify_on_filter = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = 1
	`,
		detectEmails, detectPhones, detectCreditCards, detectSSNs, detectIPV4,
		cfg.CustomEmailPattern, cfg.CustomPhonePattern, cfg.CustomCreditCardPattern,
		cfg.CustomSSNPattern, cfg.CustomIPV4Pattern,
		cfg.EmailReplacement, cfg.PhoneReplacement, cfg.CreditCardReplacement,
		cfg.SSNReplacement, cfg.IPV4Replacement,
		cfg.MonitoringInterval, notifyOnFilter,
	)

	if err != nil {
		return fmt.Errorf("failed to update config: %v", err)
	}

	return tx.Commit()
}

// LoadStringMatchPatterns loads all string match patterns from the database
func LoadStringMatchPatterns() ([]StringMatchPattern, error) {
	rows, err := db.Query(`
		SELECT id, name, pattern, enabled, replacement
		FROM string_match_patterns
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query string match patterns: %v", err)
	}
	defer rows.Close()

	var patterns []StringMatchPattern
	for rows.Next() {
		var p StringMatchPattern
		var enabled int

		err := rows.Scan(&p.ID, &p.Name, &p.Pattern, &enabled, &p.Replacement)
		if err != nil {
			return nil, fmt.Errorf("failed to scan string match pattern: %v", err)
		}

		p.Enabled = enabled == 1
		patterns = append(patterns, p)
	}

	return patterns, nil
}

// SaveStringMatchPattern saves or updates a string match pattern
func SaveStringMatchPattern(p StringMatchPattern) error {
	enabled := 0
	if p.Enabled {
		enabled = 1
	}

	if p.ID == 0 {
		// Insert new pattern
		_, err := db.Exec(`
			INSERT INTO string_match_patterns (name, pattern, enabled, replacement)
			VALUES (?, ?, ?, ?)
		`, p.Name, p.Pattern, enabled, p.Replacement)
		return err
	}

	// Update existing pattern
	_, err := db.Exec(`
		UPDATE string_match_patterns
		SET name = ?, pattern = ?, enabled = ?, replacement = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, p.Name, p.Pattern, enabled, p.Replacement, p.ID)
	return err
}

// DeleteStringMatchPattern deletes a string match pattern by ID
func DeleteStringMatchPattern(id int) error {
	_, err := db.Exec("DELETE FROM string_match_patterns WHERE id = ?", id)
	return err
}

// MarshalConfig converts Config to JSON (for compatibility)
func MarshalConfig(cfg Config) ([]byte, error) {
	return json.Marshal(cfg)
}

// UnmarshalConfig converts JSON to Config (for compatibility)
func UnmarshalConfig(data []byte) (Config, error) {
	var cfg Config
	err := json.Unmarshal(data, &cfg)
	return cfg, err
}

// LogEntry represents a filter log entry
type LogEntry struct {
	ID           int      `json:"id"`
	Timestamp    string   `json:"timestamp"`
	OriginalText string   `json:"original"`
	FilteredText string   `json:"filtered"`
	Detections   []string `json:"detections"`
}

// AddLog adds a new log entry to the database
func AddLog(originalText, filteredText string, detections []string) error {
	detectionsJSON, err := json.Marshal(detections)
	if err != nil {
		return fmt.Errorf("failed to marshal detections: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO logs (original_text, filtered_text, detections)
		VALUES (?, ?, ?)
	`, originalText, filteredText, string(detectionsJSON))

	return err
}

// GetLogs retrieves logs from the database with optional limit
func GetLogs(limit int) ([]LogEntry, error) {
	if limit <= 0 {
		limit = 100 // Default limit
	}

	rows, err := db.Query(`
		SELECT id, timestamp, original_text, filtered_text, detections
		FROM logs
		ORDER BY timestamp DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %v", err)
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		var detectionsJSON string

		err := rows.Scan(&log.ID, &log.Timestamp, &log.OriginalText, &log.FilteredText, &detectionsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log entry: %v", err)
		}

		// Unmarshal detections
		if err := json.Unmarshal([]byte(detectionsJSON), &log.Detections); err != nil {
			return nil, fmt.Errorf("failed to unmarshal detections: %v", err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// GetLogsWithPagination retrieves logs with pagination support
func GetLogsWithPagination(page, pageSize int) ([]LogEntry, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	offset := (page - 1) * pageSize

	rows, err := db.Query(`
		SELECT id, timestamp, original_text, filtered_text, detections
		FROM logs
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %v", err)
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		var detectionsJSON string

		err := rows.Scan(&log.ID, &log.Timestamp, &log.OriginalText, &log.FilteredText, &detectionsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log entry: %v", err)
		}

		// Unmarshal detections
		if err := json.Unmarshal([]byte(detectionsJSON), &log.Detections); err != nil {
			return nil, fmt.Errorf("failed to unmarshal detections: %v", err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// ClearLogs removes all log entries from the database
func ClearLogs() error {
	_, err := db.Exec("DELETE FROM logs")
	return err
}

// GetLogCount returns the total number of log entries
func GetLogCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM logs").Scan(&count)
	return count, err
}
