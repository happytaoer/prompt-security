package main

import (
	"fmt"
	"log"
	"os"

	"github.com/happytaoer/prompt-security/internal/config"
	"github.com/happytaoer/prompt-security/internal/monitor"
	"github.com/happytaoer/prompt-security/internal/web"
	"github.com/spf13/cobra"
)

func main() {
	// Initialize database
	if err := config.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.Close()

	var rootCmd = &cobra.Command{
		Use:   "prompt-security",
		Short: "Monitor clipboard for sensitive data",
		Long:  `A tool that monitors clipboard content and filters sensitive data before it's sent to language models.`,
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetString("port")
			addr := "localhost:" + port

			// Create config manager for dynamic reload
			configManager, err := config.NewManager()
			if err != nil {
				log.Fatalf("Failed to create config manager: %v", err)
			}

			// Create web server with config manager
			webServer := web.NewServer(configManager)

			// Start monitoring in background with dynamic config reload
			go monitor.ClipboardWithManager(configManager, webServer.AddLog)

			// Start web server (blocking)
			if err := webServer.Start(addr); err != nil {
				log.Fatalf("Failed to start web server: %v", err)
			}
		},
	}

	// Add flags (root command controls GUI port)
	rootCmd.PersistentFlags().String("port", "8181", "Port for web server")

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
