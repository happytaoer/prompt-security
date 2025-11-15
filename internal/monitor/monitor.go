package monitor

import (
	"log/slog"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/happytaoer/prompt-security/internal/config"
	"github.com/happytaoer/prompt-security/internal/filter"
)

// LogCallback is a function type for logging filtered data
type LogCallback func(originalText, filteredText string, replacements []filter.ReplacementInfo)

// ClipboardWithManager starts monitoring with a config manager for dynamic reload
func ClipboardWithManager(manager *config.Manager, logCallback LogCallback) {
	// Setup JSON logger
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	logger.Info("Starting clipboard monitoring with dynamic config reload...")
	logger.Info("Press Ctrl+C to stop")

	var lastContent string
	for {
		// Get current config from manager
		cfg := manager.Get()

		content, err := clipboard.ReadAll()
		if err != nil {
			logger.Error("Error reading clipboard", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Only process if content has changed
		if content != lastContent && content != "" {
			lastContent = content

			// Filter sensitive data with current config
			filtered, changed, replacementSummary := filter.SensitiveData(content, cfg)

			// If content was filtered, update clipboard
			if changed {
				updateClipboardWithNotification(content, filtered, cfg, replacementSummary, logCallback)
			}
		}

		// Sleep to avoid high CPU usage (use current config's interval)
		time.Sleep(time.Duration(cfg.MonitoringInterval) * time.Millisecond)
	}
}

// updateClipboardWithNotification updates the clipboard with filtered content and shows notifications based on configuration
func updateClipboardWithNotification(originalText, filteredText string, cfg config.Config, summary filter.ReplacementSummary, logCallback LogCallback) {
	// Setup JSON logger
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	if cfg.NotifyOnFilter {
		// Log with structured data including replacements
		if len(summary.Replacements) > 0 {
			logger.Info("Sensitive data detected and filtered",
				"replacements", summary.Replacements)
		} else {
			logger.Info("Sensitive data detected and filtered")
		}
	}

	// Call the log callback if provided
	if logCallback != nil {
		logCallback(originalText, filteredText, summary.Replacements)
	}

	err := clipboard.WriteAll(filteredText)
	if err != nil {
		logger.Error("Error writing to clipboard", "error", err)
	}
}
