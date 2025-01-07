package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the output directory",
	Long: `Serve the static files from the output directory (defaults to 'dist').
Useful for previewing your built site locally.`,
	Args: cobra.MaximumNArgs(0),
	Run:  serveFunc,
}

func serveFunc(cmd *cobra.Command, args []string) {
	config := GetConfig()

	if _, err := os.Stat(config.Dst); os.IsNotExist(err) {
		log.Error("dst directory doesn't exist")
		os.Exit(1)
		return
	}

	log.Info("running server on http://localhost:" + config.Port)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(config.Dst))))
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
