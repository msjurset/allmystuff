package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var imageAddCmd = &cobra.Command{
	Use:   "add <item-id> <file>",
	Short: "Upload an image for an item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		itemID, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid item ID: %w", err)
		}

		c := newClient()
		img, err := c.UploadImage(itemID, args[1])
		if err != nil {
			return err
		}

		fmt.Printf("Uploaded image %s (%s)\n", img.ID, img.Filename)
		return nil
	},
}

func init() {
	imageCmd.AddCommand(imageAddCmd)
}
