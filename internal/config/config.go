package config

import (
	"github.com/happytaoer/prompt-security/internal/db"
)

// Re-export types from db package for backward compatibility
type StringMatchPattern = db.StringMatchPattern
type Config = db.Config

// Initialize initializes the database
func Initialize() error {
	return db.Initialize()
}

// Close closes the database connection
func Close() error {
	return db.Close()
}

// Load loads configuration from the database
func Load() (Config, error) {
	return db.LoadConfig()
}

// Save saves the configuration to the database
func Save(config Config) error {
	return db.SaveConfig(config)
}
