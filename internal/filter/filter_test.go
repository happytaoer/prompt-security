package filter

import (
	"strings"
	"testing"

	"github.com/happytaoer/prompt-security/internal/config"
)

// TestSensitiveData_Email tests email filtering
func TestSensitiveData_Email(t *testing.T) {
	cfg := config.Config{
		DetectEmails:     true,
		EmailReplacement: "[EMAIL]",
	}

	tests := []struct {
		name              string
		input             string
		expectChanged     bool
		expectContains    string
		expectNotContains string
	}{
		{
			name:              "Single email",
			input:             "Contact me at user@example.com",
			expectChanged:     true,
			expectContains:    "[EMAIL]",
			expectNotContains: "user@example.com",
		},
		{
			name:              "Multiple emails",
			input:             "Email john@test.com or jane@company.org",
			expectChanged:     true,
			expectContains:    "[EMAIL]",
			expectNotContains: "john@test.com",
		},
		{
			name:          "No email",
			input:         "This is just plain text",
			expectChanged: false,
		},
		{
			name:              "Email with numbers",
			input:             "Contact user123@example456.com",
			expectChanged:     true,
			expectContains:    "[EMAIL]",
			expectNotContains: "user123@example456.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, changed, summary := SensitiveData(tt.input, cfg)

			if changed != tt.expectChanged {
				t.Errorf("Expected changed=%v, got %v", tt.expectChanged, changed)
			}

			if tt.expectContains != "" && !strings.Contains(filtered, tt.expectContains) {
				t.Errorf("Expected filtered text to contain '%s', got: %s", tt.expectContains, filtered)
			}

			if tt.expectNotContains != "" && strings.Contains(filtered, tt.expectNotContains) {
				t.Errorf("Expected filtered text NOT to contain '%s', got: %s", tt.expectNotContains, filtered)
			}

			if tt.expectChanged && len(summary.Replacements) == 0 {
				t.Error("Expected replacements in summary when changed=true")
			}

			if tt.expectChanged {
				for _, r := range summary.Replacements {
					if r.Type != SensitiveTypeEmail {
						t.Errorf("Expected replacement type '%s', got '%s'", SensitiveTypeEmail, r.Type)
					}
				}
			}
		})
	}
}

// TestSensitiveData_Phone tests phone number filtering
func TestSensitiveData_Phone(t *testing.T) {
	cfg := config.Config{
		DetectPhones:     true,
		PhoneReplacement: "[PHONE]",
	}

	tests := []struct {
		name              string
		input             string
		expectChanged     bool
		expectNotContains string
	}{
		{
			name:              "US phone with dashes",
			input:             "Call me at 123-456-7890",
			expectChanged:     true,
			expectNotContains: "123-456-7890",
		},
		{
			name:              "Phone with parentheses",
			input:             "Phone: (555) 123-4567",
			expectChanged:     true,
			expectNotContains: "(555) 123-4567",
		},
		{
			name:          "No phone",
			input:         "Just some text",
			expectChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, changed, summary := SensitiveData(tt.input, cfg)

			if changed != tt.expectChanged {
				t.Errorf("Expected changed=%v, got %v", tt.expectChanged, changed)
			}

			if tt.expectNotContains != "" && strings.Contains(filtered, tt.expectNotContains) {
				t.Errorf("Expected filtered text NOT to contain '%s'", tt.expectNotContains)
			}

			if tt.expectChanged {
				for _, r := range summary.Replacements {
					if r.Type != SensitiveTypePhone {
						t.Errorf("Expected replacement type '%s', got '%s'", SensitiveTypePhone, r.Type)
					}
				}
			}
		})
	}
}

// TestSensitiveData_CreditCard tests credit card filtering
func TestSensitiveData_CreditCard(t *testing.T) {
	cfg := config.Config{
		DetectCreditCards:     true,
		CreditCardReplacement: "[CARD]",
	}

	tests := []struct {
		name              string
		input             string
		expectChanged     bool
		expectNotContains string
	}{
		{
			name:              "Card with spaces",
			input:             "Card: 1234 5678 9012 3456",
			expectChanged:     true,
			expectNotContains: "1234 5678 9012 3456",
		},
		{
			name:              "Card with dashes",
			input:             "Card: 1234-5678-9012-3456",
			expectChanged:     true,
			expectNotContains: "1234-5678-9012-3456",
		},
		{
			name:              "Card without separators",
			input:             "Card: 1234567890123456",
			expectChanged:     true,
			expectNotContains: "1234567890123456",
		},
		{
			name:          "No card",
			input:         "Just text",
			expectChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, changed, _ := SensitiveData(tt.input, cfg)

			if changed != tt.expectChanged {
				t.Errorf("Expected changed=%v, got %v", tt.expectChanged, changed)
			}

			if tt.expectNotContains != "" && strings.Contains(filtered, tt.expectNotContains) {
				t.Errorf("Expected filtered text NOT to contain '%s'", tt.expectNotContains)
			}
		})
	}
}

// TestSensitiveData_SSN tests SSN filtering
func TestSensitiveData_SSN(t *testing.T) {
	cfg := config.Config{
		DetectSSNs:     true,
		SSNReplacement: "[SSN]",
	}

	tests := []struct {
		name              string
		input             string
		expectChanged     bool
		expectNotContains string
	}{
		{
			name:              "Valid SSN",
			input:             "SSN: 123-45-6789",
			expectChanged:     true,
			expectNotContains: "123-45-6789",
		},
		{
			name:          "Invalid SSN format",
			input:         "SSN: 12345678",
			expectChanged: false,
		},
		{
			name:          "No SSN",
			input:         "Just text",
			expectChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, changed, _ := SensitiveData(tt.input, cfg)

			if changed != tt.expectChanged {
				t.Errorf("Expected changed=%v, got %v", tt.expectChanged, changed)
			}

			if tt.expectNotContains != "" && strings.Contains(filtered, tt.expectNotContains) {
				t.Errorf("Expected filtered text NOT to contain '%s'", tt.expectNotContains)
			}
		})
	}
}

// TestSensitiveData_IPv4 tests IPv4 address filtering
func TestSensitiveData_IPv4(t *testing.T) {
	cfg := config.Config{
		DetectIPV4:      true,
		IPV4Replacement: "[IP]",
	}

	tests := []struct {
		name              string
		input             string
		expectChanged     bool
		expectNotContains string
	}{
		{
			name:              "Valid IPv4",
			input:             "Server IP: 192.168.1.1",
			expectChanged:     true,
			expectNotContains: "192.168.1.1",
		},
		{
			name:              "Multiple IPs",
			input:             "Connect to 10.0.0.1 or 172.16.0.1",
			expectChanged:     true,
			expectNotContains: "10.0.0.1",
		},
		{
			name:          "Invalid IP",
			input:         "IP: 999.999.999.999",
			expectChanged: false,
		},
		{
			name:          "No IP",
			input:         "Just text",
			expectChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, changed, _ := SensitiveData(tt.input, cfg)

			if changed != tt.expectChanged {
				t.Errorf("Expected changed=%v, got %v", tt.expectChanged, changed)
			}

			if tt.expectNotContains != "" && strings.Contains(filtered, tt.expectNotContains) {
				t.Errorf("Expected filtered text NOT to contain '%s'", tt.expectNotContains)
			}
		})
	}
}

// TestSensitiveData_StringMatch tests custom string pattern matching
func TestSensitiveData_StringMatch(t *testing.T) {
	cfg := config.Config{
		StringMatchPatterns: []config.StringMatchPattern{
			{
				Name:        "company_name",
				Pattern:     "Acme Corporation",
				Enabled:     true,
				Replacement: "[COMPANY]",
			},
			{
				Name:        "project_name",
				Pattern:     "Project Phoenix",
				Enabled:     true,
				Replacement: "[PROJECT]",
			},
			{
				Name:        "disabled_pattern",
				Pattern:     "Secret",
				Enabled:     false,
				Replacement: "[REDACTED]",
			},
		},
	}

	tests := []struct {
		name              string
		input             string
		expectChanged     bool
		expectContains    string
		expectNotContains string
	}{
		{
			name:              "Match company name",
			input:             "I work at Acme Corporation",
			expectChanged:     true,
			expectContains:    "[COMPANY]",
			expectNotContains: "Acme Corporation",
		},
		{
			name:              "Match project name",
			input:             "Working on Project Phoenix",
			expectChanged:     true,
			expectContains:    "[PROJECT]",
			expectNotContains: "Project Phoenix",
		},
		{
			name:           "Disabled pattern should not match",
			input:          "This is Secret information",
			expectChanged:  false,
			expectContains: "Secret",
		},
		{
			name:          "No match",
			input:         "Just plain text",
			expectChanged: false,
		},
		{
			name:              "Multiple matches",
			input:             "Acme Corporation is working on Project Phoenix",
			expectChanged:     true,
			expectContains:    "[COMPANY]",
			expectNotContains: "Acme Corporation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered, changed, summary := SensitiveData(tt.input, cfg)

			if changed != tt.expectChanged {
				t.Errorf("Expected changed=%v, got %v", tt.expectChanged, changed)
			}

			if tt.expectContains != "" && !strings.Contains(filtered, tt.expectContains) {
				t.Errorf("Expected filtered text to contain '%s', got: %s", tt.expectContains, filtered)
			}

			if tt.expectNotContains != "" && strings.Contains(filtered, tt.expectNotContains) {
				t.Errorf("Expected filtered text NOT to contain '%s', got: %s", tt.expectNotContains, filtered)
			}

			if tt.expectChanged && len(summary.Replacements) == 0 {
				t.Error("Expected replacements in summary when changed=true")
			}
		})
	}
}

// TestSensitiveData_MultipleTypes tests filtering multiple types at once
func TestSensitiveData_MultipleTypes(t *testing.T) {
	cfg := config.Config{
		DetectEmails:          true,
		DetectPhones:          true,
		DetectCreditCards:     true,
		DetectSSNs:            true,
		DetectIPV4:            true,
		EmailReplacement:      "[EMAIL]",
		PhoneReplacement:      "[PHONE]",
		CreditCardReplacement: "[CARD]",
		SSNReplacement:        "[SSN]",
		IPV4Replacement:       "[IP]",
	}

	input := `
		Contact: user@example.com
		Phone: 123-456-7890
		Card: 1234-5678-9012-3456
		SSN: 123-45-6789
		Server: 192.168.1.1
	`

	filtered, changed, summary := SensitiveData(input, cfg)

	if !changed {
		t.Error("Expected text to be changed")
	}

	// Check all sensitive data is replaced
	if strings.Contains(filtered, "user@example.com") {
		t.Error("Email should be filtered")
	}
	if strings.Contains(filtered, "123-456-7890") {
		t.Error("Phone should be filtered")
	}
	if strings.Contains(filtered, "1234-5678-9012-3456") {
		t.Error("Credit card should be filtered")
	}
	if strings.Contains(filtered, "123-45-6789") {
		t.Error("SSN should be filtered")
	}
	if strings.Contains(filtered, "192.168.1.1") {
		t.Error("IP should be filtered")
	}

	// Check replacements are present
	if !strings.Contains(filtered, "[EMAIL]") {
		t.Error("Expected [EMAIL] replacement")
	}
	if !strings.Contains(filtered, "[PHONE]") {
		t.Error("Expected [PHONE] replacement")
	}
	if !strings.Contains(filtered, "[CARD]") {
		t.Error("Expected [CARD] replacement")
	}
	if !strings.Contains(filtered, "[SSN]") {
		t.Error("Expected [SSN] replacement")
	}
	if !strings.Contains(filtered, "[IP]") {
		t.Error("Expected [IP] replacement")
	}

	// Check summary contains all types
	if len(summary.Replacements) < 5 {
		t.Errorf("Expected at least 5 replacements, got %d", len(summary.Replacements))
	}
}

// TestSensitiveData_NoDetection tests when all detection is disabled
func TestSensitiveData_NoDetection(t *testing.T) {
	cfg := config.Config{
		DetectEmails:      false,
		DetectPhones:      false,
		DetectCreditCards: false,
		DetectSSNs:        false,
		DetectIPV4:        false,
	}

	input := "user@example.com 123-456-7890 1234-5678-9012-3456"
	filtered, changed, summary := SensitiveData(input, cfg)

	if changed {
		t.Error("Expected no changes when all detection is disabled")
	}

	if filtered != input {
		t.Error("Expected filtered text to be identical to input")
	}

	if len(summary.Replacements) != 0 {
		t.Error("Expected no replacements when all detection is disabled")
	}
}

// TestSensitiveData_EmptyInput tests empty input
func TestSensitiveData_EmptyInput(t *testing.T) {
	cfg := config.Config{
		DetectEmails:     true,
		EmailReplacement: "[EMAIL]",
	}

	filtered, changed, summary := SensitiveData("", cfg)

	if changed {
		t.Error("Expected no changes for empty input")
	}

	if filtered != "" {
		t.Error("Expected empty output for empty input")
	}

	if len(summary.Replacements) != 0 {
		t.Error("Expected no replacements for empty input")
	}
}

// TestSensitiveData_DuplicateMatches tests handling of duplicate matches
func TestSensitiveData_DuplicateMatches(t *testing.T) {
	cfg := config.Config{
		DetectEmails:     true,
		EmailReplacement: "[EMAIL]",
	}

	input := "Email user@test.com and also user@test.com again"
	filtered, changed, summary := SensitiveData(input, cfg)

	if !changed {
		t.Error("Expected text to be changed")
	}

	// Both instances should be replaced
	emailCount := strings.Count(filtered, "[EMAIL]")
	if emailCount != 2 {
		t.Errorf("Expected 2 email replacements, got %d", emailCount)
	}

	// Summary should record both matches
	if len(summary.Replacements) != 2 {
		t.Errorf("Expected 2 replacements in summary, got %d", len(summary.Replacements))
	}
}

// TestReplacementInfo tests ReplacementInfo structure
func TestReplacementInfo(t *testing.T) {
	info := ReplacementInfo{
		Type:        SensitiveTypeEmail,
		Original:    "test@example.com",
		Replacement: "[EMAIL]",
	}

	if info.Type != SensitiveTypeEmail {
		t.Errorf("Expected Type=%s, got %s", SensitiveTypeEmail, info.Type)
	}
	if info.Original != "test@example.com" {
		t.Errorf("Expected Original='test@example.com', got '%s'", info.Original)
	}
	if info.Replacement != "[EMAIL]" {
		t.Errorf("Expected Replacement='[EMAIL]', got '%s'", info.Replacement)
	}
}

// TestReplacementSummary tests ReplacementSummary structure
func TestReplacementSummary(t *testing.T) {
	summary := ReplacementSummary{
		Replacements: []ReplacementInfo{
			{Type: SensitiveTypeEmail, Original: "a@b.com", Replacement: "[EMAIL]"},
			{Type: SensitiveTypePhone, Original: "123-456-7890", Replacement: "[PHONE]"},
		},
	}

	if len(summary.Replacements) != 2 {
		t.Errorf("Expected 2 replacements, got %d", len(summary.Replacements))
	}

	if summary.Replacements[0].Type != SensitiveTypeEmail {
		t.Error("First replacement should be email type")
	}
	if summary.Replacements[1].Type != SensitiveTypePhone {
		t.Error("Second replacement should be phone type")
	}
}

// BenchmarkSensitiveData_Email benchmarks email filtering
func BenchmarkSensitiveData_Email(b *testing.B) {
	cfg := config.Config{
		DetectEmails:     true,
		EmailReplacement: "[EMAIL]",
	}
	input := "Contact me at user@example.com for more information"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SensitiveData(input, cfg)
	}
}

// BenchmarkSensitiveData_MultipleTypes benchmarks filtering multiple types
func BenchmarkSensitiveData_MultipleTypes(b *testing.B) {
	cfg := config.Config{
		DetectEmails:          true,
		DetectPhones:          true,
		DetectCreditCards:     true,
		DetectSSNs:            true,
		DetectIPV4:            true,
		EmailReplacement:      "[EMAIL]",
		PhoneReplacement:      "[PHONE]",
		CreditCardReplacement: "[CARD]",
		SSNReplacement:        "[SSN]",
		IPV4Replacement:       "[IP]",
	}
	input := "Email: user@test.com, Phone: 123-456-7890, Card: 1234-5678-9012-3456, SSN: 123-45-6789, IP: 192.168.1.1"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SensitiveData(input, cfg)
	}
}

// BenchmarkSensitiveData_NoMatch benchmarks when no sensitive data is found
func BenchmarkSensitiveData_NoMatch(b *testing.B) {
	cfg := config.Config{
		DetectEmails:     true,
		EmailReplacement: "[EMAIL]",
	}
	input := "This is just plain text without any sensitive information"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SensitiveData(input, cfg)
	}
}

// BenchmarkSensitiveData_StringMatch benchmarks string pattern matching
func BenchmarkSensitiveData_StringMatch(b *testing.B) {
	cfg := config.Config{
		StringMatchPatterns: []config.StringMatchPattern{
			{Name: "company", Pattern: "Acme Corp", Enabled: true, Replacement: "[COMPANY]"},
			{Name: "project", Pattern: "Project X", Enabled: true, Replacement: "[PROJECT]"},
		},
	}
	input := "Working at Acme Corp on Project X"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SensitiveData(input, cfg)
	}
}
