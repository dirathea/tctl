package cmd

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tctl [URL]",
	Short: "TableCheck Control / tctl for homework assignment",
	Long:  "A simple CLI to access a URL and print all the content to stdout",
	Args:  cobra.MinimumNArgs(1),
	RunE:  rootRun,
}

func rootRun(cmd *cobra.Command, args []string) (err error) {
	client := resty.New()
	resp, err := client.R().Get(args[0])
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
