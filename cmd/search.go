package cmd

import (
	"errors"
	"f/helper"
	"fmt"

	"github.com/spf13/cobra"
)

func runSearch(args []string, helperFunc func(string, string) ([]string, error)) error {
	if len(args) < 1 {
		return errors.New("not enough arguments")
	}
	query := args[0]
	// Read the files from the directory
	dir, err := helper.GetDirectoryFromArgs(args, 2)
	if err != nil {
		fmt.Println("Error getting directory:", err)
		return err
	}
	results, err := helperFunc(query, dir)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return errors.New("no files found for search criteria")
	}
	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

func runNameSearch(cmd *cobra.Command, args []string) error {
	return runSearch(args, helper.SearchByName)
}

func runContentSearch(cmd *cobra.Command, args []string) error {
	return runSearch(args, helper.SearchByContent)
}

var searchCmd = &cobra.Command{
	Use:   "search <name|content> <query> <directory>",
	Short: "Search for files",
	Long:  "Search for files by name or content",
}

var nameSearchCmd = &cobra.Command{
	Use:   "name <query> <directory>",
	Short: "Search for files by name",
	Long:  "Search for files by name",
	RunE:  runNameSearch,
}

var contentSearchCmd = &cobra.Command{
	Use:   "content <query> <directory>",
	Short: "Search for files by content",
	Long:  "Search for files by content",
	RunE:  runContentSearch,
}

func init() {
	searchCmd.AddCommand(nameSearchCmd)
	searchCmd.AddCommand(contentSearchCmd)
}
