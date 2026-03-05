package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"allmystuff/internal/model"
)

func printItemsTable(items []model.Item) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tBRAND\tCONDITION\tTAGS")
	for _, it := range items {
		tags := make([]string, len(it.Tags))
		for i, t := range it.Tags {
			tags[i] = t.Name
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", it.ID, it.Name, it.Brand, it.Condition, strings.Join(tags, ", "))
	}
	w.Flush()
}

func printItemDetail(it *model.Item) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", it.ID)
	fmt.Fprintf(w, "Name:\t%s\n", it.Name)
	if it.Description != "" {
		fmt.Fprintf(w, "Description:\t%s\n", it.Description)
	}
	if it.Brand != "" {
		fmt.Fprintf(w, "Brand:\t%s\n", it.Brand)
	}
	if it.Model != "" {
		fmt.Fprintf(w, "Model:\t%s\n", it.Model)
	}
	if it.SerialNumber != "" {
		fmt.Fprintf(w, "Serial:\t%s\n", it.SerialNumber)
	}
	if it.PurchaseDate != nil {
		fmt.Fprintf(w, "Purchase Date:\t%s\n", it.PurchaseDate.Format("2006-01-02"))
	}
	if it.PurchasePrice != nil {
		fmt.Fprintf(w, "Purchase Price:\t$%.2f\n", *it.PurchasePrice)
	}
	if it.EstimatedValue != nil {
		fmt.Fprintf(w, "Estimated Value:\t$%.2f\n", *it.EstimatedValue)
	}
	if it.Condition != "" {
		fmt.Fprintf(w, "Condition:\t%s\n", it.Condition)
	}
	if it.Notes != "" {
		fmt.Fprintf(w, "Notes:\t%s\n", it.Notes)
	}
	fmt.Fprintf(w, "Created:\t%s\n", it.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Updated:\t%s\n", it.UpdatedAt.Format("2006-01-02 15:04:05"))

	if len(it.Tags) > 0 {
		tags := make([]string, len(it.Tags))
		for i, t := range it.Tags {
			tags[i] = t.Name
		}
		fmt.Fprintf(w, "Tags:\t%s\n", strings.Join(tags, ", "))
	}
	if len(it.Images) > 0 {
		fmt.Fprintf(w, "Images:\t%d\n", len(it.Images))
		for _, img := range it.Images {
			fmt.Fprintf(w, "  %s\t%s\n", img.ID, img.Filename)
		}
	}
	w.Flush()
}

func printTagsTable(tags []model.Tag) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME")
	for _, t := range tags {
		fmt.Fprintf(w, "%d\t%s\n", t.ID, t.Name)
	}
	w.Flush()
}

