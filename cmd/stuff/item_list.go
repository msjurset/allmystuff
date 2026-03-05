package main

import (
	"fmt"
	"os"

	"allmystuff/internal/store"

	"github.com/spf13/cobra"
)

var itemListCmd = &cobra.Command{
	Use:   "list",
	Short: "List items",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()

		q, _ := cmd.Flags().GetString("query")
		tag, _ := cmd.Flags().GetString("tag")
		cond, _ := cmd.Flags().GetString("condition")

		filter := store.ItemFilter{
			Query:     q,
			Tag:       tag,
			Condition: cond,
		}

		if flagJSON {
			data, err := c.ListItemsRaw(filter)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		items, err := c.ListItems(filter)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			fmt.Println("No items found.")
			return nil
		}

		printItemsTable(items)
		return nil
	},
}

func init() {
	itemListCmd.Flags().StringP("query", "q", "", "Search query")
	itemListCmd.Flags().String("tag", "", "Filter by tag")
	itemListCmd.Flags().String("condition", "", "Filter by condition")
	itemCmd.AddCommand(itemListCmd)
}
