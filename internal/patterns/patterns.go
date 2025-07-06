package patterns

import (
	"regexp"

	"github.com/happytaoer/prompt-security/internal/config"
)

// PatternType represents the type of pattern matching to use
type PatternType int

const (
	// RegexPattern indicates a regular expression pattern
	RegexPattern PatternType = iota
	// StringMatchPattern indicates a simple string match pattern
	StringMatchPattern
)

// Pattern represents a pattern to detect sensitive data
type Pattern struct {
	Type    PatternType
	Pattern string
}

// Default patterns for sensitive data detection
const (
	DefaultEmailPatternStr      = `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	DefaultPhonePatternStr      = `(\+\d{1,3}[\s-]?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}`
	DefaultCreditCardPatternStr = `\b(?:\d{4}[- ]?){3}\d{4}\b`
	DefaultSSNPatternStr        = `\b\d{3}-\d{2}-\d{4}\b`
	DefaultIPV4PatternStr       = `\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`
)

var (
	// Default compiled patterns
	defaultEmailPattern      = regexp.MustCompile(DefaultEmailPatternStr)
	defaultPhonePattern      = regexp.MustCompile(DefaultPhonePatternStr)
	defaultCreditCardPattern = regexp.MustCompile(DefaultCreditCardPatternStr)
	defaultSSNPattern        = regexp.MustCompile(DefaultSSNPatternStr)
	defaultIPV4Pattern       = regexp.MustCompile(DefaultIPV4PatternStr)
)

// GetEmailPattern returns the appropriate email pattern based on configuration
func GetEmailPattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomEmailPattern != "" {
		// Try to compile custom pattern, fallback to default if it fails
		pattern, err := regexp.Compile(cfg.CustomEmailPattern)
		if err == nil {
			return pattern
		}
	}
	return defaultEmailPattern
}

// GetPhonePattern returns the appropriate phone pattern based on configuration
func GetPhonePattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomPhonePattern != "" {
		// Try to compile custom pattern, fallback to default if it fails
		pattern, err := regexp.Compile(cfg.CustomPhonePattern)
		if err == nil {
			return pattern
		}
	}
	return defaultPhonePattern
}

// GetCreditCardPattern returns the appropriate credit card pattern based on configuration
func GetCreditCardPattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomCreditCardPattern != "" {
		// Try to compile custom pattern, fallback to default if it fails
		pattern, err := regexp.Compile(cfg.CustomCreditCardPattern)
		if err == nil {
			return pattern
		}
	}
	return defaultCreditCardPattern
}

// GetSSNPattern returns the appropriate SSN pattern based on configuration
func GetSSNPattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomSSNPattern != "" {
		// Try to compile custom pattern, fallback to default if it fails
		pattern, err := regexp.Compile(cfg.CustomSSNPattern)
		if err == nil {
			return pattern
		}
	}
	return defaultSSNPattern
}

// GetIPV4Pattern returns the appropriate IPv4 pattern based on configuration
func GetIPV4Pattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomIPV4Pattern != "" {
		// Try to compile custom pattern, fallback to default if it fails
		pattern, err := regexp.Compile(cfg.CustomIPV4Pattern)
		if err == nil {
			return pattern
		}
	}
	return defaultIPV4Pattern
}
