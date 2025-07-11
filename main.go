package main

import (
	"fmt"
	"log"
	"os"

	"github.com/happytaoer/prompt-security/internal/config"
	"github.com/happytaoer/prompt-security/internal/monitor"
	"github.com/spf13/cobra"
)

func main() {
	// Load configuration
	cfg, err := config.Load()

	if err != nil {
		log.Printf("Warning: Using default configuration: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "prompt-security",
		Short: "Monitor clipboard for sensitive data",
		Long:  `A tool that monitors clipboard content and filters sensitive data before it's sent to language models.`,
	}

	// Monitor command
	var monitorCmd = &cobra.Command{
		Use:   "monitor",
		Short: "Start monitoring clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			// Check if interval flag is set, and override config if it is
			if interval, _ := cmd.Flags().GetInt("interval"); interval > 0 {
				cfg.MonitoringInterval = interval
			}
			monitor.Clipboard(cfg)
		},
	}

	// Add flags to monitor command
	monitorCmd.Flags().Int("interval", 0, "Override monitoring interval in milliseconds")

	// Config command
	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			config.Show(cfg)
		},
	}

	// Add commands
	rootCmd.AddCommand(monitorCmd, configCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
