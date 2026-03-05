package main

import "github.com/spf13/cobra"

var itemCmd = &cobra.Command{
	Use:   "item",
	Short: "Manage items",
}

func init() {
	rootCmd.AddCommand(itemCmd)
}
