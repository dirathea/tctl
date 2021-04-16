package cmd

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/go-resty/resty/v2"
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

	for i := 0; i < profileLoop; i++ {
		wg.Add(1)
		go runFunc(client, url)
	}

	wg.Wait()

	// Sort ascending
	sort.Stable(results)

	fmt.Printf("Min %v\n", results[0].Duration)
	fmt.Printf("Max %v\n", results[len(results)-1].Duration)
	fmt.Printf("Mean %v\n", results.Mean())
	fmt.Printf("Median %v\n", results.Median())
	fmt.Printf("Total Success %v %%\n", results.PercentSuccess())
	fmt.Printf("All Failed status code %v\n", results.AllErrorStatusCode())

	// Sort based on response size
	sort.SliceStable(results, func(x, y int) bool {
		return results[x].Size < results[y].Size
	})

	fmt.Printf("Min size %v\n", results[0].Size)
	fmt.Printf("Max size %v\n", results[len(results)-1].Size)

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
