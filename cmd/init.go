package cmd

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

//go:embed embed/*.*
var initFS embed.FS

var forceFlag bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [src]",
	Short: "Scaffold a new project",
	Long: `Create a new project directory with a basic structure and configuration file.
Defaults to 'site' for the source directory.`,
	Args: cobra.MaximumNArgs(2),
	Run:  initFunc,
}

func initFunc(cmd *cobra.Command, args []string) {
	var src, dst string

	if len(args) == 2 {
		src, dst = args[0], args[1]
	} else if len(args) == 0 {
		config := GetConfig()
		src, dst = config.Src, config.Dst
	} else {
		cmd.PrintErrln("Usage: build [src] [dst] (provide both or neither)")
		return
	}

	if err := WriteConfig(Config{
		Src:  src,
		Dst:  dst,
		Port: DegaultPort,
	}); err != nil {
		cmd.PrintErrln("failed to create config", "err", err)
		return
	}

	cmd.Printf("Scaffolding new project in %s...\n", src)
	err := createScaffold(src)
	if err != nil {
		// TODO: Clean up if we fuck up?
		cmd.Println("Error scaffolding project:", err)
		return
	} else {
		cmd.Println("Project created successfully!")
	}

	cmd.Printf("Building (%s) -> %s\n", src, dst)

	if err := buildSite(src, dst, nil); err != nil {
		cmd.PrintErrln("Failed to build site:", err)
		return
	}

	cmd.Println("Built successfully")
}

func createScaffold(src string) error {
	// Check if the target directory exists
	if _, err := os.Stat(src); err == nil {
		// Directory exists
		if !forceFlag {
			return fmt.Errorf("directory %s already exists; use --force to overwrite", src)
		}
	}

	// Create a temporary cleanup function
	var createdDirs []string
	var createdFiles []string

	cleanup := func() {
		for _, file := range createdFiles {
			_ = os.Remove(file) // Attempt to remove files
		}
		for i := len(createdDirs) - 1; i >= 0; i-- { // Reverse order for directories
			_ = os.Remove(createdDirs[i])
		}
	}

	// Remove the directory if it exists and force is enabled
	if forceFlag {
		if err := os.RemoveAll(src); err != nil {
			return fmt.Errorf("failed to remove existing directory %s: %w", src, err)
		}
	}

	// Create the root directory
	if err := os.MkdirAll(src, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", src, err)
	}
	createdDirs = append(createdDirs, src)

	// Create subdirectories
	dirs := []string{
		"content",
		"content/posts",
		"static",
		"templates",
	}
	for _, dir := range dirs {
		fullPath := filepath.Join(src, dir)
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			cleanup()
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
		createdDirs = append(createdDirs, fullPath)
	}

	// Copy embedded files
	files := map[string]string{
		"embed/index.md":   "content/index.md",
		"embed/post_1.md":  "content/posts/1.md",
		"embed/post_2.md":  "content/posts/2.md",
		"embed/post_3.md":  "content/posts/3.md",
		"embed/styles.css": "static/styles.css",
		"embed/index.tmpl": "templates/index.tmpl",
		"embed/post.tmpl":  "templates/post.tmpl",
	}
	for id, path := range files {
		fullPath := filepath.Join(src, path)
		content, err := initFS.ReadFile(id)
		if err != nil {
			cleanup()
			return fmt.Errorf("failed to read embedded file %s: %w", id, err)
		}

		// Write the file
		if err := os.WriteFile(fullPath, content, os.ModePerm); err != nil {
			cleanup()
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
		createdFiles = append(createdFiles, fullPath)
	}

	return nil
}

func init() {
	initCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Overwrite existing files")
	rootCmd.AddCommand(initCmd)
}
