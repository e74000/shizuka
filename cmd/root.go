package cmd

import (
    "github.com/spf13/cobra"
    "os"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
    Use:   "shizuka",
    Short: "Shizuka - A minimalist static site generator",
    Long: `Shizuka is a CLI tool for building static sites from markdown files. 
It supports templating, live development, and project scaffolding with ease.

Examples:
  # Build a static site
  shizuka build

  # Serve the site locally and watch for changes
  shizuka dev

  # Initialize a new project
  shizuka init`,
}

// Execute is the entry point for running the CLI
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
