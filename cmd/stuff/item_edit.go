package main

import (
	"fmt"
	"os"

	"allmystuff/internal/model"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var itemEditCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Update an item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("invalid item ID: %w", err)
		}

		c := newClient()

		// Fetch existing item to merge with
		existing, err := c.GetItem(id)
		if err != nil {
			return fmt.Errorf("fetching item: %w", err)
		}

		input := model.ItemInput{
			Name:         existing.Name,
			Description:  existing.Description,
			Brand:        existing.Brand,
			Model:        existing.Model,
			SerialNumber: existing.SerialNumber,
			Condition:    existing.Condition,
			Notes:        existing.Notes,
		}

		if existing.PurchaseDate != nil {
			s := existing.PurchaseDate.Format("2006-01-02")
			input.PurchaseDate = &s
		}
		input.PurchasePrice = existing.PurchasePrice
		input.EstimatedValue = existing.EstimatedValue

		// Preserve existing tags
		for _, t := range existing.Tags {
			input.Tags = append(input.Tags, t.Name)
		}

		// Override with provided flags
		if cmd.Flags().Changed("name") {
			input.Name, _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("description") {
			input.Description, _ = cmd.Flags().GetString("description")
		}
		if cmd.Flags().Changed("brand") {
			input.Brand, _ = cmd.Flags().GetString("brand")
		}
		if cmd.Flags().Changed("model") {
			input.Model, _ = cmd.Flags().GetString("model")
		}
		if cmd.Flags().Changed("serial") {
			input.SerialNumber, _ = cmd.Flags().GetString("serial")
		}
		if cmd.Flags().Changed("purchase-date") {
			s, _ := cmd.Flags().GetString("purchase-date")
			input.PurchaseDate = &s
		}
		if cmd.Flags().Changed("purchase-price") {
			v, _ := cmd.Flags().GetFloat64("purchase-price")
			input.PurchasePrice = &v
		}
		if cmd.Flags().Changed("estimated-value") {
			v, _ := cmd.Flags().GetFloat64("estimated-value")
			input.EstimatedValue = &v
		}
		if cmd.Flags().Changed("condition") {
			input.Condition, _ = cmd.Flags().GetString("condition")
		}
		if cmd.Flags().Changed("notes") {
			input.Notes, _ = cmd.Flags().GetString("notes")
		}
		if cmd.Flags().Changed("tag") {
			input.Tags, _ = cmd.Flags().GetStringSlice("tag")
		}

		if flagJSON {
			data, err := c.UpdateItemRaw(id, input)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		item, err := c.UpdateItem(id, input)
		if err != nil {
			return err
		}

		fmt.Printf("Updated item %s\n", item.ID)
		return nil
	},
}

func init() {
	itemEditCmd.Flags().String("name", "", "Item name")
	itemEditCmd.Flags().String("description", "", "Description")
	itemEditCmd.Flags().String("brand", "", "Brand")
	itemEditCmd.Flags().String("model", "", "Model")
	itemEditCmd.Flags().String("serial", "", "Serial number")
	itemEditCmd.Flags().String("purchase-date", "", "Purchase date (YYYY-MM-DD)")
	itemEditCmd.Flags().Float64("purchase-price", 0, "Purchase price")
	itemEditCmd.Flags().Float64("estimated-value", 0, "Estimated value")
	itemEditCmd.Flags().String("condition", "", "Condition")
	itemEditCmd.Flags().String("notes", "", "Notes")
	itemEditCmd.Flags().StringSlice("tag", nil, "Tags (replaces all)")
	itemCmd.AddCommand(itemEditCmd)
}
