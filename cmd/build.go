package cmd

import (
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [src] [dst]",
	Short: "Build the site",
	Long: `Build the static site from the source directory to the destination directory.
Defaults are 'site' for src and 'dist' for dst.`,
	Args: cobra.MaximumNArgs(2),
	Run:  buildFunc,
}

func buildFunc(cmd *cobra.Command, args []string) {
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

	cmd.Printf("Building (%s) -> %s\n", src, dst)

	if err := buildSite(src, dst, nil); err != nil {
		cmd.PrintErrln("Failed to build site:", err)
		return
	}

	cmd.Println("Built successfully")
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
