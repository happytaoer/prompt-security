package monitor

import (
	"log/slog"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/happytaoer/prompt-security/internal/config"
	"github.com/happytaoer/prompt-security/internal/filter"
)

// Clipboard starts monitoring the clipboard for sensitive data
func Clipboard(cfg config.Config) {
	// Setup JSON logger
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)

	logger.Info("Starting clipboard monitoring...")
	logger.Info("Press Ctrl+C to stop")

	var lastContent string
	for {
		content, err := clipboard.ReadAll()
		if err != nil {
			logger.Error("Error reading clipboard", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Only process if content has changed
		if content != lastContent && content != "" {
			lastContent = content

			// Filter sensitive data
			filtered, changed, replacementSummary := filter.SensitiveData(content, cfg)

			// If content was filtered, update clipboard
			if changed {
				updateClipboardWithNotification(filtered, cfg, replacementSummary)
			}
		}

		// Sleep to avoid high CPU usage
		time.Sleep(time.Duration(cfg.MonitoringInterval) * time.Millisecond)
	}
}

// updateClipboardWithNotification updates the clipboard with filtered content and shows notifications based on configuration
func updateClipboardWithNotification(content string, cfg config.Config, summary filter.ReplacementSummary) {
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

	err := clipboard.WriteAll(content)
	if err != nil {
		logger.Error("Error writing to clipboard", "error", err)
	}
}
