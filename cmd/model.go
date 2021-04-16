package cmd

import (
	"sort"
	"time"
)

type RunResult struct {
	Duration   time.Duration
	Success    bool
	StatusCode int
	Size       int64
}

type RunResultSlice []RunResult

func (rs RunResultSlice) Len() int {
	return len(rs)
}

func (rs RunResultSlice) Less(i, j int) bool {
	return rs[i].Duration.Nanoseconds() < rs[j].Duration.Nanoseconds()
}

func (rs RunResultSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Mean calculates mean from result slice
func (rs RunResultSlice) Mean() time.Duration {
	total, _ := time.ParseDuration("0s")
	for _, r := range rs {
		total += r.Duration
	}

	return time.Duration(total.Nanoseconds() / int64(len(rs)))
}

// Median calculate median from result slice
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

// PercentSuccess calculate percentage of success result
func (rs RunResultSlice) PercentSuccess() float64 {
	totalSuccess := 0
	for _, result := range rs {
		if result.Success {
			totalSuccess += 1
		}
	}

	return (float64(totalSuccess) / float64(len(rs))) * 100
}

// AllErrorStatusCode returns all status code that not success
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
