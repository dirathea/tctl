package cmd

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	profileLoop int
	rootCmd     = &cobra.Command{
		Use:   "tctl [URL]",
		Short: "TableCheck CLI / tctl for homework assignment",
		Long:  "A simple CLI to access a URL and print all the content to stdout",
		Args:  cobra.MinimumNArgs(1),
		RunE:  rootRun,
	}
)

func rootRun(cmd *cobra.Command, args []string) (err error) {
	url := args[0]
	client := resty.New()
	if profileLoop == 0 {
		// No profile flag added.
		resp, err := client.R().Get(url)
		if err != nil {
			return err
		}
		fmt.Println(resp)
		return nil
	}

	results := RunResultSlice{}
	var wg sync.WaitGroup

	runFunc := func(client *resty.Client, url string) {
		resp, _ := client.R().Get(url)
		result := RunResult{
			Duration:   resp.Time(),
			Success:    resp.IsSuccess(),
			StatusCode: resp.StatusCode(),
			Size:       resp.Size(),
		}

		results = append(results, result)
		wg.Done()
	}

	// Do multiple requests simultaneously
	for i := 0; i < profileLoop; i++ {
		wg.Add(1)
		go runFunc(client, url)
	}

	wg.Wait()

	// Sort ascending
	sort.Stable(results)
	tableData := [][]string{
		{"Min", fmt.Sprint(results[0].Duration)},
		{"Max", fmt.Sprint(results[len(results)-1].Duration)},
		{"Mean", fmt.Sprint(results.Mean())},
		{"Median", fmt.Sprint(results.Median())},
		{"Success Percentage", fmt.Sprint(results.PercentSuccess())},
		{"All Failed status code", fmt.Sprint(results.AllErrorStatusCode())},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.AppendBulk(tableData)

	// Sort based on response size
	sort.SliceStable(results, func(x, y int) bool {
		return results[x].Size < results[y].Size
	})

	sizeData := [][]string{
		{"Min Response size in bytes", fmt.Sprint(results[0].Size)},
		{"Max Response size in bytes", fmt.Sprint(results[len(results)-1].Size)},
	}

	table.AppendBulk(sizeData)
	table.Render()
	return nil
}

func init() {
	rootCmd.Flags().IntVarP(&profileLoop, "profile", "p", 0, "Do profiling to the URL by accessing it several times and calculate response duration statistic")
}

func Execute(version string) {
	// Assign supplied version to vars for versionCmd
	Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
