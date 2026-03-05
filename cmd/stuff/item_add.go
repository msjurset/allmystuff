package main

import (
	"fmt"
	"os"

	"allmystuff/internal/model"

	"github.com/spf13/cobra"
)

var itemAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new item",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}

		desc, _ := cmd.Flags().GetString("description")
		brand, _ := cmd.Flags().GetString("brand")
		mdl, _ := cmd.Flags().GetString("model")
		serial, _ := cmd.Flags().GetString("serial")
		purchDate, _ := cmd.Flags().GetString("purchase-date")
		purchPrice, _ := cmd.Flags().GetFloat64("purchase-price")
		estValue, _ := cmd.Flags().GetFloat64("estimated-value")
		condition, _ := cmd.Flags().GetString("condition")
		notes, _ := cmd.Flags().GetString("notes")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		input := model.ItemInput{
			Name:         name,
			Description:  desc,
			Brand:        brand,
			Model:        mdl,
			SerialNumber: serial,
			Condition:    condition,
			Notes:        notes,
			Tags:         tags,
		}

		if cmd.Flags().Changed("purchase-date") {
			input.PurchaseDate = &purchDate
		}
		if cmd.Flags().Changed("purchase-price") {
			input.PurchasePrice = &purchPrice
		}
		if cmd.Flags().Changed("estimated-value") {
			input.EstimatedValue = &estValue
		}

		if flagJSON {
			data, err := c.CreateItemRaw(input)
			if err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, string(data))
			return nil
		}

		item, err := c.CreateItem(input)
		if err != nil {
			return err
		}

		fmt.Printf("Created item %s\n", item.ID)
		return nil
	},
}

func init() {
	itemAddCmd.Flags().String("name", "", "Item name (required)")
	itemAddCmd.Flags().String("description", "", "Description")
	itemAddCmd.Flags().String("brand", "", "Brand")
	itemAddCmd.Flags().String("model", "", "Model")
	itemAddCmd.Flags().String("serial", "", "Serial number")
	itemAddCmd.Flags().String("purchase-date", "", "Purchase date (YYYY-MM-DD)")
	itemAddCmd.Flags().Float64("purchase-price", 0, "Purchase price")
	itemAddCmd.Flags().Float64("estimated-value", 0, "Estimated value")
	itemAddCmd.Flags().String("condition", "", "Condition")
	itemAddCmd.Flags().String("notes", "", "Notes")
	itemAddCmd.Flags().StringSlice("tag", nil, "Tags (repeatable)")
	itemCmd.AddCommand(itemAddCmd)
}
