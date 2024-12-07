package cmd

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [dst]",
	Short: "Serve the output directory",
	Long: `Serve the static files from the output directory (defaults to 'dist').
Useful for previewing your built site locally.`,
	Args: cobra.MaximumNArgs(1),
	Run:  serveFunc,
}

func serveFunc(cmd *cobra.Command, args []string) {
	var dst string

	if len(args) > 1 {
		cmd.PrintErrln("Usage: serve [dst]")
		os.Exit(1)
		return
	}

	config := GetConfig()
	if len(args) == 1 {
		dst = args[0]
	} else {
		dst = config.Dst
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		cmd.PrintErrln("Destination directory does not exist")
		os.Exit(1)
		return
	}

	port := portFlag
	if port == "" {
		port = config.Port
	}

	cmd.Printf("Running on http://localhost:%s\n", port)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(dst))))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		cmd.PrintErrln("Failed to start server:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
