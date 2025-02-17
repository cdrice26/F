package cmd

import (
    "f/helper"
    "fmt"
    "path/filepath"

    "github.com/spf13/cobra"
)

func runCopy(cmd *cobra.Command, args []string) {
    if len(args) < 2 {
        fmt.Println("Usage: copy <source>... <destination>")
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
            err := helper.Copy(match, dst, false)
            if err != nil {
                fmt.Printf("Error copying %s: %v\n", match, err)
            } else {
                fmt.Printf("Copied %s to %s successfully\n", match, dst)
            }
        }
    }
}

var copyCmd = &cobra.Command{
    Use:   "copy",
    Short: "Copy files, directories, and wildcards",
    Long:  `Copy files, directories, and wildcards from source to destination.`,
    Run:   runCopy,
}




