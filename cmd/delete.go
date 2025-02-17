package cmd

import (
	"f/helper"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func runDelete(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: delete <source>...")
		return
	}

	srcs := args

	for _, src := range srcs {
		// Expand the source path to handle wildcards
		matches, err := filepath.Glob(src)
		if err != nil {
			fmt.Printf("Error processing source path: %v\n", err)
			continue
		}

		if len(matches) == 0 {
			fmt.Printf("No files matched the source pattern: %s\n", src)
			continue
		}

		for _, match := range matches {
			err := helper.Delete(match)
			if err != nil {
				fmt.Printf("Error deleting %s: %v\n", match, err)
			} else {
				fmt.Printf("Deleted %s successfully\n", match)
			}
		}
	}
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete files, directories, and wildcards",
	Long:  `Delete files, directories, and wildcards from source to destination.`,
	Run:   runDelete,
}
