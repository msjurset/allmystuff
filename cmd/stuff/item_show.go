package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var itemShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show item details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid item ID: %w", err)
		}

		c := newClient()

		if flagJSON {
			data, err := c.GetItemRaw(id)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		item, err := c.GetItem(id)
		if err != nil {
			return err
		}

		printItemDetail(item)
		return nil
	},
}

func init() {
	itemCmd.AddCommand(itemShowCmd)
}
