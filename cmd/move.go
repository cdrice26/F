package cmd

import (
	"f/helper"
	"fmt"
    "path/filepath"

    "github.com/spf13/cobra"
)

func runMove(cmd *cobra.Command, args []string) {
    if len(args) < 2 {
        fmt.Println("Usage: move <source>... <destination>")
        return
    }

    dst := args[len(args)-1]
    srcs := args[:len(args)-1]

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
            err := helper.Copy(match, dst, true)
            if err != nil {
                fmt.Printf("Error moving %s: %v\n", match, err)
            } else {
                fmt.Printf("Moved %s to %s successfully\n", match, dst)
            }
        }
    }
}

var moveCmd = &cobra.Command{
    Use:   "move",
    Short: "Move a file",
    Long:  `Move a file from source to destination. The filename should not be included in the destination path.`,
    Run:   runMove,
}




