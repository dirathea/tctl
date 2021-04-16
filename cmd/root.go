package cmd

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var (
	profileLoop int
	rootCmd     = &cobra.Command{
		Use:   "tctl [URL]",
		Short: "TableCheck Control / tctl for homework assignment",
		Long:  "A simple CLI to access a URL and print all the content to stdout",
		Args:  cobra.MinimumNArgs(1),
		RunE:  rootRun,
	}
)

func getMean(durations []int) float64 {
	total := 0.0
	for _, duration := range durations {
		total += float64(duration)
	}

	return total / float64(len(durations))
}

func getMedian(durations sort.IntSlice) float64 {
	durations.Sort()
	isEven := len(durations)%2 == 0
	midIndex := len(durations) / 2
	if isEven {
		return float64(durations[midIndex])
	}

	return (float64(durations[midIndex]) + float64(durations[midIndex+1])) / 2
}

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

	type RunResult struct {
		Duration   time.Duration
		Success    bool
		StatusCode int
	}

	results := []RunResult{}

	runFunc := func(client *resty.Client, url string) {
		resp, err := client.R().Get(url)
		result := RunResult{
			Duration:   resp.Time(),
			Success:    err == nil,
			StatusCode: resp.StatusCode(),
		}

		results = append(results, result)
	}

	var wg sync.WaitGroup

	for i := 0; i < profileLoop; i++ {
		wg.Add(1)
		go runFunc(client, url)
	}

	wg.Wait()

	return nil
}

func init() {
	rootCmd.Flags().IntVarP(&profileLoop, "profile", "p", 0, "Do profiling to the URL by accessing it several times and calculate response duration statistic")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
