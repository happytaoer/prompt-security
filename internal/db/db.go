package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// ConfigModel represents the configuration table (GORM model)
type ConfigModel struct {
	ID                      uint   `gorm:"primaryKey;check:id=1"`
	DetectEmails            bool   `gorm:"default:true"`
	DetectPhones            bool   `gorm:"default:true"`
	DetectCreditCards       bool   `gorm:"default:true"`
	DetectSSNs              bool   `gorm:"default:true"`
	DetectIPV4              bool   `gorm:"default:true"`
	CustomEmailPattern      string `gorm:"default:''"`
	CustomPhonePattern      string `gorm:"default:''"`
	CustomCreditCardPattern string `gorm:"default:''"`
	CustomSSNPattern        string `gorm:"default:''"`
	CustomIPV4Pattern       string `gorm:"default:''"`
	EmailReplacement        string `gorm:"default:'security@example.com'"`
	PhoneReplacement        string `gorm:"default:'+1-555-123-4567'"`
	CreditCardReplacement   string `gorm:"default:'XXXX-XXXX-XXXX-XXXX'"`
	SSNReplacement          string `gorm:"default:'XXX-XX-XXXX'"`
	IPV4Replacement         string `gorm:"default:'0.0.0.0'"`
	MonitoringIntervalMs    int    `gorm:"default:500"`
	NotifyOnFilter          bool   `gorm:"default:true"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (ConfigModel) TableName() string {
	return "config"
}

// StringMatchPatternModel represents a string match pattern (GORM model)
type StringMatchPatternModel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Pattern     string `gorm:"not null"`
	Enabled     bool   `gorm:"default:true"`
	Replacement string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (StringMatchPatternModel) TableName() string {
	return "string_match_patterns"
}

// LogEntryModel represents a log entry (GORM model)
type LogEntryModel struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Timestamp    time.Time `gorm:"index:idx_logs_timestamp,sort:desc;default:CURRENT_TIMESTAMP"`
	OriginalText string    `gorm:"not null"`
	FilteredText string    `gorm:"not null"`
	Detections   string    `gorm:"not null"` // JSON string
	CreatedAt    time.Time
}

func (LogEntryModel) TableName() string {
	return "logs"
}

// Initialize initializes the database connection and creates tables if needed
func Initialize() error {
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	db = database

	// Auto migrate tables
	if err := db.AutoMigrate(&ConfigModel{}, &StringMatchPatternModel{}, &LogEntryModel{}); err != nil {
		return fmt.Errorf("failed to migrate tables: %v", err)
	}

	// Insert default config if not exists
	var count int64
	db.Model(&ConfigModel{}).Count(&count)
	if count == 0 {
		defaultConfig := &ConfigModel{ID: 1}
		if err := db.Create(defaultConfig).Error; err != nil {
			return fmt.Errorf("failed to create default config: %v", err)
		}
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
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

// StringMatchPattern represents a string match pattern (API model)
type StringMatchPattern struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Enabled     bool   `json:"enabled"`
	Replacement string `json:"replacement"`
}

// Config represents the application configuration (API model)
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
	var configModel ConfigModel
	if err := db.First(&configModel, 1).Error; err != nil {
		return Config{}, fmt.Errorf("failed to load config: %v", err)
	}

	// Load string match patterns
	patterns, err := LoadStringMatchPatterns()
	if err != nil {
		return Config{}, fmt.Errorf("failed to load string match patterns: %v", err)
	}

	cfg := Config{
		DetectEmails:            configModel.DetectEmails,
		DetectPhones:            configModel.DetectPhones,
		DetectCreditCards:       configModel.DetectCreditCards,
		DetectSSNs:              configModel.DetectSSNs,
		DetectIPV4:              configModel.DetectIPV4,
		CustomEmailPattern:      configModel.CustomEmailPattern,
		CustomPhonePattern:      configModel.CustomPhonePattern,
		CustomCreditCardPattern: configModel.CustomCreditCardPattern,
		CustomSSNPattern:        configModel.CustomSSNPattern,
		CustomIPV4Pattern:       configModel.CustomIPV4Pattern,
		EmailReplacement:        configModel.EmailReplacement,
		PhoneReplacement:        configModel.PhoneReplacement,
		CreditCardReplacement:   configModel.CreditCardReplacement,
		SSNReplacement:          configModel.SSNReplacement,
		IPV4Replacement:         configModel.IPV4Replacement,
		MonitoringInterval:      configModel.MonitoringIntervalMs,
		NotifyOnFilter:          configModel.NotifyOnFilter,
		StringMatchPatterns:     patterns,
	}

	return cfg, nil
}

// SaveConfig saves the configuration to the database
func SaveConfig(cfg Config) error {
	configModel := ConfigModel{
		ID:                      1,
		DetectEmails:            cfg.DetectEmails,
		DetectPhones:            cfg.DetectPhones,
		DetectCreditCards:       cfg.DetectCreditCards,
		DetectSSNs:              cfg.DetectSSNs,
		DetectIPV4:              cfg.DetectIPV4,
		CustomEmailPattern:      cfg.CustomEmailPattern,
		CustomPhonePattern:      cfg.CustomPhonePattern,
		CustomCreditCardPattern: cfg.CustomCreditCardPattern,
		CustomSSNPattern:        cfg.CustomSSNPattern,
		CustomIPV4Pattern:       cfg.CustomIPV4Pattern,
		EmailReplacement:        cfg.EmailReplacement,
		PhoneReplacement:        cfg.PhoneReplacement,
		CreditCardReplacement:   cfg.CreditCardReplacement,
		SSNReplacement:          cfg.SSNReplacement,
		IPV4Replacement:         cfg.IPV4Replacement,
		MonitoringIntervalMs:    cfg.MonitoringInterval,
		NotifyOnFilter:          cfg.NotifyOnFilter,
	}

	return db.Save(&configModel).Error
}

// LoadStringMatchPatterns loads all string match patterns from the database
func LoadStringMatchPatterns() ([]StringMatchPattern, error) {
	var models []StringMatchPatternModel
	if err := db.Order("id").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to query string match patterns: %v", err)
	}

	patterns := make([]StringMatchPattern, len(models))
	for i, m := range models {
		patterns[i] = StringMatchPattern{
			ID:          int(m.ID),
			Name:        m.Name,
			Pattern:     m.Pattern,
			Enabled:     m.Enabled,
			Replacement: m.Replacement,
		}
	}

	return patterns, nil
}

// SaveStringMatchPattern saves or updates a string match pattern
func SaveStringMatchPattern(p StringMatchPattern) error {
	model := StringMatchPatternModel{
		ID:          uint(p.ID),
		Name:        p.Name,
		Pattern:     p.Pattern,
		Enabled:     p.Enabled,
		Replacement: p.Replacement,
	}

	return db.Save(&model).Error
}

// DeleteStringMatchPattern deletes a string match pattern by ID
func DeleteStringMatchPattern(id int) error {
	return db.Delete(&StringMatchPatternModel{}, id).Error
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

// LogEntry represents a filter log entry (API model)
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

	logModel := LogEntryModel{
		Timestamp:    time.Now(),
		OriginalText: originalText,
		FilteredText: filteredText,
		Detections:   string(detectionsJSON),
	}

	return db.Create(&logModel).Error
}

// GetLogs retrieves logs from the database with optional limit
func GetLogs(limit int) ([]LogEntry, error) {
	if limit <= 0 {
		limit = 100 // Default limit
	}

	var models []LogEntryModel
	if err := db.Order("timestamp DESC").Limit(limit).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to query logs: %v", err)
	}

	return convertLogModelsToEntries(models)
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

	var models []LogEntryModel
	if err := db.Order("timestamp DESC").Limit(pageSize).Offset(offset).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to query logs: %v", err)
	}

	return convertLogModelsToEntries(models)
}

// convertLogModelsToEntries converts GORM models to API models
func convertLogModelsToEntries(models []LogEntryModel) ([]LogEntry, error) {
	logs := make([]LogEntry, len(models))
	for i, m := range models {
		var detections []string
		if err := json.Unmarshal([]byte(m.Detections), &detections); err != nil {
			return nil, fmt.Errorf("failed to unmarshal detections: %v", err)
		}

		logs[i] = LogEntry{
			ID:           int(m.ID),
			Timestamp:    m.Timestamp.Format(time.RFC3339),
			OriginalText: m.OriginalText,
			FilteredText: m.FilteredText,
			Detections:   detections,
		}
	}

	return logs, nil
}

// ClearLogs removes all log entries from the database
func ClearLogs() error {
	return db.Where("1 = 1").Delete(&LogEntryModel{}).Error
}

// GetLogCount returns the total number of log entries
func GetLogCount() (int, error) {
	var count int64
	err := db.Model(&LogEntryModel{}).Count(&count).Error
	return int(count), err
}
