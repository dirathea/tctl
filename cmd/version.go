package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version    = ""
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
