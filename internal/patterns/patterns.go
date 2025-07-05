package patterns

import (
	"regexp"
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

// Patterns for sensitive data detection
const (
	EmailPatternStr      = `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	PhonePatternStr      = `(\+\d{1,3}[\s-]?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}`
	CreditCardPatternStr = `\b(?:\d{4}[- ]?){3}\d{4}\b`
	SSNPatternStr        = `\b\d{3}-\d{2}-\d{4}\b`
	IPV4PatternStr       = `\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`
)

var (
	// EmailPattern matches common email address formats
	EmailPattern = regexp.MustCompile(EmailPatternStr)

	// PhonePattern matches common phone number formats
	PhonePattern = regexp.MustCompile(PhonePatternStr)

	// CreditCardPattern matches credit card number formats (16 digits in groups of 4)
	CreditCardPattern = regexp.MustCompile(CreditCardPatternStr)

	// SSNPattern matches US Social Security Number format
	SSNPattern = regexp.MustCompile(SSNPatternStr)

	// IPV4Pattern matches IPv4 addresses like 192.168.1.1
	IPV4Pattern = regexp.MustCompile(IPV4PatternStr)
)
