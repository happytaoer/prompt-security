package patterns

import (
	"regexp"
	"sync"

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

// PatternCache caches compiled regular expressions to avoid recompilation
type PatternCache struct {
	mu       sync.RWMutex
	patterns map[string]*regexp.Regexp
}

// globalCache is the global pattern cache instance
var globalCache = &PatternCache{
	patterns: make(map[string]*regexp.Regexp),
}

// Get retrieves a compiled pattern from cache or compiles and caches it
func (pc *PatternCache) Get(key string, patternStr string) (*regexp.Regexp, error) {
	// Fast path: read lock for cache hit
	pc.mu.RLock()
	if pattern, ok := pc.patterns[key]; ok {
		pc.mu.RUnlock()
		return pattern, nil
	}
	pc.mu.RUnlock()

	// Slow path: compile and cache
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, err
	}

	pc.mu.Lock()
	pc.patterns[key] = pattern
	pc.mu.Unlock()

	return pattern, nil
}

// Clear removes all cached patterns (useful for testing or config reload)
func (pc *PatternCache) Clear() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.patterns = make(map[string]*regexp.Regexp)
}

// GetEmailPattern returns the appropriate email pattern based on configuration
func GetEmailPattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomEmailPattern != "" {
		// Try to get from cache or compile custom pattern, fallback to default if it fails
		pattern, err := globalCache.Get("email", cfg.CustomEmailPattern)
		if err == nil {
			return pattern
		}
	}
	return defaultEmailPattern
}

// GetPhonePattern returns the appropriate phone pattern based on configuration
func GetPhonePattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomPhonePattern != "" {
		// Try to get from cache or compile custom pattern, fallback to default if it fails
		pattern, err := globalCache.Get("phone", cfg.CustomPhonePattern)
		if err == nil {
			return pattern
		}
	}
	return defaultPhonePattern
}

// GetCreditCardPattern returns the appropriate credit card pattern based on configuration
func GetCreditCardPattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomCreditCardPattern != "" {
		// Try to get from cache or compile custom pattern, fallback to default if it fails
		pattern, err := globalCache.Get("creditCard", cfg.CustomCreditCardPattern)
		if err == nil {
			return pattern
		}
	}
	return defaultCreditCardPattern
}

// GetSSNPattern returns the appropriate SSN pattern based on configuration
func GetSSNPattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomSSNPattern != "" {
		// Try to get from cache or compile custom pattern, fallback to default if it fails
		pattern, err := globalCache.Get("ssn", cfg.CustomSSNPattern)
		if err == nil {
			return pattern
		}
	}
	return defaultSSNPattern
}

// GetIPV4Pattern returns the appropriate IPv4 pattern based on configuration
func GetIPV4Pattern(cfg *config.Config) *regexp.Regexp {
	if cfg != nil && cfg.CustomIPV4Pattern != "" {
		// Try to get from cache or compile custom pattern, fallback to default if it fails
		pattern, err := globalCache.Get("ipv4", cfg.CustomIPV4Pattern)
		if err == nil {
			return pattern
		}
	}
	return defaultIPV4Pattern
}
