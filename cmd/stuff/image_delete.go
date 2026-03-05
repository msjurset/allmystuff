package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var imageDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid image ID: %w", err)
		}

		c := newClient()
		if err := c.DeleteImage(id); err != nil {
			return err
		}

		fmt.Printf("Deleted image %s\n", id)
		return nil
	},
}

func init() {
	imageCmd.AddCommand(imageDeleteCmd)
}
