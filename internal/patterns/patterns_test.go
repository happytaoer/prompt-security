package patterns

import (
	"regexp"
	"sync"
	"testing"

	"github.com/happytaoer/prompt-security/internal/config"
)

// TestPatternCache_Get tests the basic cache functionality
func TestPatternCache_Get(t *testing.T) {
	cache := &PatternCache{
		patterns: make(map[string]*regexp.Regexp),
	}

	// Test first compilation (cache miss)
	pattern1, err := cache.Get("test1", `\d+`)
	if err != nil {
		t.Fatalf("Failed to compile pattern: %v", err)
	}
	if pattern1 == nil {
		t.Fatal("Expected non-nil pattern")
	}

	// Test cache hit (should return same instance)
	pattern2, err := cache.Get("test1", `\d+`)
	if err != nil {
		t.Fatalf("Failed to get cached pattern: %v", err)
	}
	if pattern1 != pattern2 {
		t.Error("Expected same pattern instance from cache")
	}

	// Test different key
	pattern3, err := cache.Get("test2", `[a-z]+`)
	if err != nil {
		t.Fatalf("Failed to compile different pattern: %v", err)
	}
	if pattern1 == pattern3 {
		t.Error("Expected different pattern instances for different keys")
	}
}

// TestPatternCache_InvalidPattern tests error handling
func TestPatternCache_InvalidPattern(t *testing.T) {
	cache := &PatternCache{
		patterns: make(map[string]*regexp.Regexp),
	}

	// Test invalid regex
	_, err := cache.Get("invalid", `[invalid`)
	if err == nil {
		t.Error("Expected error for invalid regex pattern")
	}
}

// TestPatternCache_Clear tests cache clearing
func TestPatternCache_Clear(t *testing.T) {
	cache := &PatternCache{
		patterns: make(map[string]*regexp.Regexp),
	}

	// Add some patterns
	cache.Get("test1", `\d+`)
	cache.Get("test2", `[a-z]+`)

	if len(cache.patterns) != 2 {
		t.Errorf("Expected 2 cached patterns, got %d", len(cache.patterns))
	}

	// Clear cache
	cache.Clear()

	if len(cache.patterns) != 0 {
		t.Errorf("Expected 0 cached patterns after clear, got %d", len(cache.patterns))
	}
}

// TestPatternCache_Concurrent tests thread safety
func TestPatternCache_Concurrent(t *testing.T) {
	cache := &PatternCache{
		patterns: make(map[string]*regexp.Regexp),
	}

	var wg sync.WaitGroup
	concurrency := 100

	// Concurrent reads and writes
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := "pattern"
			pattern := `\d+`
			_, err := cache.Get(key, pattern)
			if err != nil {
				t.Errorf("Goroutine %d: Failed to get pattern: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Should only have one cached pattern despite concurrent access
	if len(cache.patterns) != 1 {
		t.Errorf("Expected 1 cached pattern, got %d", len(cache.patterns))
	}
}

// TestGetEmailPattern_WithCache tests email pattern caching
func TestGetEmailPattern_WithCache(t *testing.T) {
	// Clear global cache before test
	globalCache.Clear()

	cfg := &config.Config{
		CustomEmailPattern: `[a-zA-Z0-9]+@test\.com`,
	}

	// First call should compile and cache
	pattern1 := GetEmailPattern(cfg)
	if pattern1 == nil {
		t.Fatal("Expected non-nil pattern")
	}

	// Second call should return cached pattern
	pattern2 := GetEmailPattern(cfg)
	if pattern1 != pattern2 {
		t.Error("Expected same pattern instance from cache")
	}

	// Verify cache contains the pattern
	if len(globalCache.patterns) != 1 {
		t.Errorf("Expected 1 cached pattern, got %d", len(globalCache.patterns))
	}
}

// TestGetEmailPattern_Default tests default pattern fallback
func TestGetEmailPattern_Default(t *testing.T) {
	globalCache.Clear()

	// Test with nil config
	pattern1 := GetEmailPattern(nil)
	if pattern1 != defaultEmailPattern {
		t.Error("Expected default email pattern for nil config")
	}

	// Test with empty custom pattern
	cfg := &config.Config{
		CustomEmailPattern: "",
	}
	pattern2 := GetEmailPattern(cfg)
	if pattern2 != defaultEmailPattern {
		t.Error("Expected default email pattern for empty custom pattern")
	}

	// Cache should be empty
	if len(globalCache.patterns) != 0 {
		t.Errorf("Expected 0 cached patterns, got %d", len(globalCache.patterns))
	}
}

// TestGetEmailPattern_InvalidCustom tests fallback on invalid custom pattern
func TestGetEmailPattern_InvalidCustom(t *testing.T) {
	globalCache.Clear()

	cfg := &config.Config{
		CustomEmailPattern: `[invalid`,
	}

	pattern := GetEmailPattern(cfg)
	if pattern != defaultEmailPattern {
		t.Error("Expected default email pattern for invalid custom pattern")
	}
}

// TestAllPatternGetters tests all pattern getter functions
func TestAllPatternGetters(t *testing.T) {
	globalCache.Clear()

	cfg := &config.Config{
		CustomEmailPattern:      `custom@email`,
		CustomPhonePattern:      `\d{11}`,
		CustomCreditCardPattern: `\d{16}`,
		CustomSSNPattern:        `\d{9}`,
		CustomIPV4Pattern:       `\d+\.\d+\.\d+\.\d+`,
	}

	tests := []struct {
		name   string
		getter func(*config.Config) *regexp.Regexp
	}{
		{"Email", GetEmailPattern},
		{"Phone", GetPhonePattern},
		{"CreditCard", GetCreditCardPattern},
		{"SSN", GetSSNPattern},
		{"IPv4", GetIPV4Pattern},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern1 := tt.getter(cfg)
			pattern2 := tt.getter(cfg)

			if pattern1 == nil {
				t.Error("Expected non-nil pattern")
			}

			if pattern1 != pattern2 {
				t.Error("Expected same pattern instance from cache")
			}
		})
	}

	// Should have 5 cached patterns
	if len(globalCache.patterns) != 5 {
		t.Errorf("Expected 5 cached patterns, got %d", len(globalCache.patterns))
	}
}

// BenchmarkGetEmailPattern_WithoutCache benchmarks pattern compilation without cache
func BenchmarkGetEmailPattern_WithoutCache(b *testing.B) {
	customPattern := `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = regexp.Compile(customPattern)
	}
}

// BenchmarkGetEmailPattern_WithCache benchmarks pattern retrieval with cache
func BenchmarkGetEmailPattern_WithCache(b *testing.B) {
	globalCache.Clear()

	cfg := &config.Config{
		CustomEmailPattern: `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEmailPattern(cfg)
	}
}

// BenchmarkPatternCache_ConcurrentAccess benchmarks concurrent cache access
func BenchmarkPatternCache_ConcurrentAccess(b *testing.B) {
	cache := &PatternCache{
		patterns: make(map[string]*regexp.Regexp),
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cache.Get("test", `\d+`)
		}
	})
}
