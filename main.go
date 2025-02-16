package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "f <command> [flags]",
		Short: "An intuitive, consistent file manager for the command line",
		Long:  `A longer description of your application that spans multiple lines and likely contains examples and usage of using your application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from F!")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
