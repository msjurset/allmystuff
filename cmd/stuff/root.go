package main

import (
	"fmt"
	"os"

	"allmystuff/internal/client"
	"allmystuff/internal/secret"

	"github.com/spf13/cobra"
)

var (
	flagURL    string
	flagAPIKey string
	flagJSON   bool
)

var rootCmd = &cobra.Command{
	Use:   "stuff",
	Short: "CLI for the allmystuff inventory API",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "API base URL (env: ALLMYSTUFF_URL, default: http://localhost:8080)")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "API key (env: ALLMYSTUFF_API_KEY)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output raw JSON")
}

func newClient() *client.Client {
	u := flagURL
	if u == "" {
		u = os.Getenv("ALLMYSTUFF_URL")
	}
	if u == "" {
		u = "http://localhost:8080"
	}

	key := flagAPIKey
	if key == "" {
		key = os.Getenv("ALLMYSTUFF_API_KEY")
	}
	if key != "" {
		resolved, err := secret.Resolve(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving API key: %v\n", err)
			os.Exit(1)
		}
		key = resolved
	}

	return client.New(u, key)
}
