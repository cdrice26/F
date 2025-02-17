package cmd

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
)

func runRename(cmd *cobra.Command, args []string) {
    if len(args) != 2 {
        fmt.Println("Usage: rename <source> <newname>")
        return
    }

    src := args[0]
    newName := args[1]

    srcDir := filepath.Dir(src)
    dst := filepath.Join(srcDir, newName)

    err := os.Rename(src, dst)
    if err != nil {
        fmt.Printf("Error renaming %s to %s: %v\n", src, newName, err)
    } else {
        fmt.Printf("Renamed %s to %s successfully\n", src, newName)
    }
}

var renameCmd = &cobra.Command{
    Use:   "rename",
    Short: "Rename a file or directory",
    Long:  `Rename a file or directory to a new name within the same directory.`,
    Run:   runRename,
}