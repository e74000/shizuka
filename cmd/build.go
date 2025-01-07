package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site",
	Long:  `Builds site from the src directory to the dst directory, as specified in shizuka_conf.json`,
	Args:  cobra.MaximumNArgs(0),
	Run:   buildFunc,
}

func buildFunc(cmd *cobra.Command, args []string) {
	config := GetConfig()

	if !exists(config.Src) {
		log.Error("source directory doesn't exist", "directory", config.Src)
	}

	if err := buildSite(config.Src, config.Dst, makeOpts(config)); err != nil {
		log.Error("failed to build site", "error", err)
		return
	}

	log.Info("built site successfully")
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
