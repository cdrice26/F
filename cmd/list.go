package cmd

import (
	"f/helper"
	"fmt"
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

	// Check if the tree flag is set
	isTree, err := cmd.Flags().GetBool("tree")
	if err != nil {
		fmt.Println("Error getting flag value:", err)
		return
	}

	// Check if the hidden flag is set
	includeHidden, err := cmd.Flags().GetBool("hidden")
	if err != nil {
		fmt.Println("Error getting flag value:", err)
		return
	}

	// Read the files from the directory
	dir, err := helper.GetDirectoryFromArgs(args)
	if err != nil {
		fmt.Println("Error getting directory:", err)
		return
	}

	files := []helper.Entry{}
	if isTree {
		files, err = helper.GetDirectoryTree(dir, includeHidden)
		if err != nil {
			fmt.Println("Error reading the directory:", err)
			return
		}
	} else {
		files, err = helper.GetFileListing(dir, includeHidden)
		if err != nil {
			fmt.Println("Error reading the directory:", err)
			return
		}
	}

	longestFileName := 0
	for _, file := range files {
		if len(file.DirEntry.Name()) > longestFileName {
			longestFileName = len(file.DirEntry.Name())
		}
	}

	formatStr := ""
	if isTree {
		formatStr = fmt.Sprintf("%%-10s %%-10s %%-30s\n")
	} else {
		formatStr = fmt.Sprintf("%%-%ds %%-10s %%-10s %%-30s\n", longestFileName)
	}
	dataFormatStr := ""
	if isTree {
		dataFormatStr = fmt.Sprintf("%%-10s %%-10s %%-30s %%s\n")
	} else {
		dataFormatStr = fmt.Sprintf("%%-%ds %%-10s %%-10s %%-30s\n", longestFileName)
	}

	// Display metadata for each file
	if isTree {
		fmt.Printf(formatStr, "Size", "Type", "Modified")
	} else {
		fmt.Printf(formatStr, "Name", "Size", "Type", "Modified")
	}
	for _, file := range files {
		fileInfo, err := file.DirEntry.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		// Get file type
		fileType := ""
		fileSize := ""
		if file.DirEntry.IsDir() {
			fileType = "Directory"
			if dir_sizes {
				dirSize, err := helper.GetDirSize(filepath.Join(dir, file.DirEntry.Name()))
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
		if isTree {
			fmt.Printf(dataFormatStr, fileSize, fileType, fileInfo.ModTime().Format(time.RFC1123), file.Path)
		} else {
			fmt.Printf(dataFormatStr, file.DirEntry.Name(), fileSize, fileType, fileInfo.ModTime().Format(time.RFC1123))
		}
	}
}

var listCmd = &cobra.Command{
	Use:   "list [directory]",
	Short: "List files in the specified directory",
	Long:  `List files in the specified directory. If no directory is specified, the current directory is used.`,
	Run:   runList,
}

func init() {
	listCmd.Flags().BoolP("no-directory-sizes", "n", false, "Do not display directory sizes")
	listCmd.Flags().BoolP("tree", "t", false, "Display directory tree")
	listCmd.Flags().BoolP("hidden", "a", false, "Include hidden files and directories")
}
