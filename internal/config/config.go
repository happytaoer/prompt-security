package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// StringMatchPattern represents a simple string match pattern for sensitive data
type StringMatchPattern struct {
	Name        string `json:"name"`        // Name of the pattern (for identification)
	Pattern     string `json:"pattern"`     // The exact string to match
	Enabled     bool   `json:"enabled"`     // Whether this pattern is enabled
	Replacement string `json:"replacement"` // What to replace matches with
}

// Config represents the application configuration
type Config struct {
	// Patterns to detect
	DetectEmails      bool `json:"detect_emails"`
	DetectPhones      bool `json:"detect_phones"`
	DetectCreditCards bool `json:"detect_credit_cards"`
	DetectSSNs        bool `json:"detect_ssns"`

	// String match patterns
	StringMatchPatterns []StringMatchPattern `json:"string_match_patterns"`

	// Replacements
	EmailReplacement      string `json:"email_replacement"`
	PhoneReplacement      string `json:"phone_replacement"`
	CreditCardReplacement string `json:"credit_card_replacement"`
	SSNReplacement        string `json:"ssn_replacement"`
	APIKeyReplacement     string `json:"api_key_replacement"`

	// Monitoring settings
	MonitoringInterval int  `json:"monitoring_interval_ms"` // in milliseconds
	NotifyOnFilter     bool `json:"notify_on_filter"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		DetectEmails:      true,
		DetectPhones:      true,
		DetectCreditCards: true,
		DetectSSNs:        true,

		StringMatchPatterns: []StringMatchPattern{},

		EmailReplacement:      "security@example.com",
		PhoneReplacement:      "+1-555-123-4567",
		CreditCardReplacement: "XXXX-XXXX-XXXX-XXXX",
		SSNReplacement:        "XXX-XX-XXXX",

		MonitoringInterval: 500,
		NotifyOnFilter:     true,
	}
}

// Load loads configuration from file or creates default if not exists
func Load() (Config, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return DefaultConfig(), err
	}

	configPath := filepath.Join(configDir, "config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		config := DefaultConfig()
		if err := Save(config); err != nil {
			return config, fmt.Errorf("failed to create default config: %v", err)
		}
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// Save saves the configuration to file
func Save(config Config) error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Show displays the current configuration
func Show(config Config) {
	fmt.Println("Current Configuration:")
	fmt.Println("---------------------")
	fmt.Printf("Detect Emails: %v", config.DetectEmails)
	if config.DetectEmails {
		fmt.Printf(" (Replacement: %s)", config.EmailReplacement)
	}
	fmt.Println()

	fmt.Printf("Detect Phone Numbers: %v", config.DetectPhones)
	if config.DetectPhones {
		fmt.Printf(" (Replacement: %s)", config.PhoneReplacement)
	}
	fmt.Println()

	fmt.Printf("Detect Credit Cards: %v", config.DetectCreditCards)
	if config.DetectCreditCards {
		fmt.Printf(" (Replacement: %s)", config.CreditCardReplacement)
	}
	fmt.Println()

	fmt.Printf("Detect SSNs: %v", config.DetectSSNs)
	if config.DetectSSNs {
		fmt.Printf(" (Replacement: %s)", config.SSNReplacement)
	}
	fmt.Println()

	fmt.Printf("Monitoring Interval: %d ms\n", config.MonitoringInterval)
	fmt.Printf("Notify on Filter: %v\n", config.NotifyOnFilter)
}

// getConfigDir returns the configuration directory
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".prompt-security")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return configDir, nil
}
