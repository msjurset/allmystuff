package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var itemDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid item ID: %w", err)
		}

		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			fmt.Printf("Delete item %s? [y/N] ", id)
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if !strings.EqualFold(strings.TrimSpace(scanner.Text()), "y") {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		c := newClient()
		if err := c.DeleteItem(id); err != nil {
			return err
		}

		fmt.Printf("Deleted item %s\n", id)
		return nil
	},
}

func init() {
	itemDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation")
	itemCmd.AddCommand(itemDeleteCmd)
}
