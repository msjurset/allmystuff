package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()

		if flagJSON {
			data, err := c.ListTagsRaw()
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		tags, err := c.ListTags()
		if err != nil {
			return err
		}

		if len(tags) == 0 {
			fmt.Println("No tags found.")
			return nil
		}

		printTagsTable(tags)
		return nil
	},
}

func init() {
	tagCmd.AddCommand(tagListCmd)
	rootCmd.AddCommand(tagCmd)
}
