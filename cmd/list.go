package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func getDirectoryFromArgs(args []string) (string, error) {
	if len(args) == 0 {
		return os.Getwd()
	} else {
		return args[0], nil
	}
}

func formatSize(size int64) string {
	if size >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(size)/float64(1024*1024*1024))
	} else if size >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(1024*1024))
	} else if size >= 1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/float64(1024))
	}
	return fmt.Sprintf("%d Bytes", size)
}

func runList(cmd *cobra.Command, args []string) {

	// Read the files from the directory
	dir, err := getDirectoryFromArgs(args)
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
	fmt.Println("-----------------------------------------------------------------------------------------")
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		// Get file type
		fileType := ""
		if file.IsDir() {
			fileType = "Directory"
		} else {
			fileType = "File"
		}

		// Displaying file metadata
		fmt.Printf(dataFormatStr, fileInfo.Name(), formatSize(fileInfo.Size()), fileType, fileInfo.ModTime().Format(time.RFC1123))
	}
}

var listCmd = &cobra.Command{
	Use:   "list [directory]",
	Short: "List files in the current directory",
	Long:  `List files in the current directory.`,
	Run:   runList,
}
