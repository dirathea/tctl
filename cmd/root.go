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
		Short: "TableCheck CLI / tctl for homework assignment",
		Long:  "A simple CLI to access a URL and print all the content to stdout",
		Args:  cobra.MinimumNArgs(1),
		RunE:  rootRun,
	}
)

type RunResult struct {
	Duration   time.Duration
	Success    bool
	StatusCode int
	Size       int64
}

type RunResultSlice []RunResult

func (sl RunResultSlice) Len() int {
	return len(sl)
}

func (rs RunResultSlice) Less(i, j int) bool {
	return rs[i].Duration.Nanoseconds() < rs[j].Duration.Nanoseconds()
}

func (rs RunResultSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs RunResultSlice) Mean() time.Duration {
	sort.Stable(rs)
	total, _ := time.ParseDuration("0s")
	for _, r := range rs {
		total += r.Duration
	}

	return time.Duration(total.Nanoseconds() / int64(len(rs)))
}

func (rs RunResultSlice) Median() time.Duration {
	sort.Stable(rs)
	midIndex := len(rs) / 2
	if rs.isEven() {
		return time.Duration((rs[midIndex-1].Duration + rs[midIndex].Duration) / 2)
	}

	return rs[midIndex].Duration
}

func (rs RunResultSlice) isEven() bool {
	return len(rs)%2 == 0
}

func (rs RunResultSlice) PercentSuccess() float64 {
	totalSuccess := 0
	for _, result := range rs {
		if result.Success {
			totalSuccess += 1
		}
	}

	return (float64(totalSuccess) / float64(len(rs))) * 100
}

func (sl RunResultSlice) AllErrorStatusCode() []int {
	allCodes := map[int]bool{}

	for _, result := range sl {
		if !result.Success {
			allCodes[result.StatusCode] = true
		}
	}

	results := []int{}

	for c := range allCodes {
		results = append(results, c)
	}

	return results
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
