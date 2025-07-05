package filter

import (
	"regexp"
	"strings"

	"github.com/happytaoer/prompt-security/internal/config"
	"github.com/happytaoer/prompt-security/internal/patterns"
)

// Sensitive data type constants
const (
	SensitiveTypeEmail      = "email"
	SensitiveTypePhone      = "phone"
	SensitiveTypeCreditCard = "credit_card"
	SensitiveTypeSSN        = "ssn"
	SensitiveTypeAPIKey     = "api_key"
)

// ReplacementInfo stores information about a single sensitive data replacement
type ReplacementInfo struct {
	Type        string // Type of sensitive data (email, phone, etc.)
	Original    string // Original sensitive data
	Replacement string // What it was replaced with
}

// ReplacementSummary contains all replacements made during filtering
type ReplacementSummary struct {
	Replacements []ReplacementInfo
}

// SensitiveData filters sensitive data from text and returns the filtered text,
// a boolean indicating whether any changes were made, and a summary of replacements
func SensitiveData(text string, cfg config.Config) (string, bool, ReplacementSummary) {
	original := text
	summary := ReplacementSummary{}

	// Helper function to find and replace sensitive data with regex
	findAndReplaceRegex := func(pattern *regexp.Regexp, replacement string, dataType string) {
		matches := pattern.FindAllString(text, -1)
		for _, match := range matches {
			summary.Replacements = append(summary.Replacements, ReplacementInfo{
				Type:        dataType,
				Original:    match,
				Replacement: replacement,
			})
		}
		text = pattern.ReplaceAllString(text, replacement)
	}

	// Helper function to find and replace sensitive data with string match
	findAndReplaceString := func(pattern string, replacement string, dataType string) {
		if strings.Contains(text, pattern) {
			summary.Replacements = append(summary.Replacements, ReplacementInfo{
				Type:        dataType,
				Original:    pattern,
				Replacement: replacement,
			})
			text = strings.ReplaceAll(text, pattern, replacement)
		}
	}

	// Filter emails
	if cfg.DetectEmails {
		findAndReplaceRegex(patterns.EmailPattern, cfg.EmailReplacement, SensitiveTypeEmail)
	}

	// Filter phone numbers
	if cfg.DetectPhones {
		findAndReplaceRegex(patterns.PhonePattern, cfg.PhoneReplacement, SensitiveTypePhone)
	}

	// Filter credit card numbers
	if cfg.DetectCreditCards {
		findAndReplaceRegex(patterns.CreditCardPattern, cfg.CreditCardReplacement, SensitiveTypeCreditCard)
	}

	// Filter SSNs
	if cfg.DetectSSNs {
		findAndReplaceRegex(patterns.SSNPattern, cfg.SSNReplacement, SensitiveTypeSSN)
	}

	// Filter string match patterns
	for _, stringPattern := range cfg.StringMatchPatterns {
		if stringPattern.Enabled {
			findAndReplaceString(stringPattern.Pattern, stringPattern.Replacement, stringPattern.Name)
		}
	}

	return text, text != original, summary
}
