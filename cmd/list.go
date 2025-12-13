package cmd

import (
	"f/helper"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func runList(cmd *cobra.Command, args []string) {
	// Check if the no-directory-sizes flag is set
	noDirSizes, err := cmd.Flags().GetBool("no-directory-sizes")
	if err != nil {
		fmt.Println("Error getting flag value:", err)
		return
	}
	dir_sizes := !noDirSizes

	// Read the files from the directory
	dir, err := helper.GetDirectoryFromArgs(args)
	if err != nil {
		fmt.Println("Error getting directory:", err)
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading the directory:", err)
		return
	}

	longestFileName := 0
	for _, file := range files {
		if len(file.Name()) > longestFileName {
			longestFileName = len(file.Name())
		}
	}

	formatStr := fmt.Sprintf("%%-%ds %%-10s %%-10s %%-30s\n", longestFileName)
	dataFormatStr := fmt.Sprintf("%%-%ds %%-10s %%-10s %%-30s\n", longestFileName)

	// Display metadata for each file
	fmt.Printf(formatStr, "Name", "Size", "Type", "Modified")
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		// Get file type
		fileType := ""
		fileSize := ""
		if file.IsDir() {
			fileType = "Directory"
			if dir_sizes {
				dirSize, err := helper.GetDirSize(filepath.Join(dir, file.Name()))
				if err != nil {
					fileSize = "Unknown"
				} else {
					fileSize = helper.FormatSize(dirSize)
				}
			} else {
				fileSize = "N/A"
			}
		} else {
			fileType = "File"
			fileSize = helper.FormatSize(fileInfo.Size())
		}

		// Displaying file metadata
		fmt.Printf(dataFormatStr, fileInfo.Name(), fileSize, fileType, fileInfo.ModTime().Format(time.RFC1123))
	}
}

var listCmd = &cobra.Command{
	Use:   "list [directory]",
	Short: "List files in the current directory",
	Long:  `List files in the current directory.`,
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolP("no-directory-sizes", "n", false, "Do not display directory sizes")
}
