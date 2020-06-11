package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"

	"gopkg.in/yaml.v1"
)

// addResult adds a result to the results object used in reports
func addResult(res *result, err error) {
	// Decrement thread count
	defer func() {
		remainingThreads--
		atomic.AddInt32(&activeThreads, 1)
	}()

	if res == nil {
		res = &result{}
	}

	if err != nil {
		res.Err = err
	}

	// Add to results slice
	results = append(results, res)

	// Increment progress bar if output is set to plaintext mode
	if output == "plaintext" {
		bar.Increment()
	}
}

// displayReport outputs a report of proxy performance & health
func displayReport(results []*result) {
	success := 0
	fail := 0
	averageTTFB := int64(0)

	// Calculate stats
	for _, v := range results {
		if v == nil {
			continue
		}

		if v.StatusCode == -1 {
			fail++
			continue
		}

		averageTTFB += v.Latency.TTFB
		success++
	}

	if success > 0 {
		averageTTFB /= int64(success)
	}

	switch output {
	case "json":
		// JSON formatted report
		report := &report{Success: success, Fail: fail, AverageTTFB: averageTTFB, Results: results}

		b, err := json.Marshal(report)
		if err != nil {
			panic(err)
		}

		// Write to stdout
		os.Stdout.Write(b)
		return
	case "yaml":
		// YAML formatted report
		report := &report{Success: success, Fail: fail, AverageTTFB: averageTTFB, Results: results}

		b, err := yaml.Marshal(report)
		if err != nil {
			panic(err)
		}

		// Write to stdout
		os.Stdout.Write(b)
		return
	default:
		break
	}

	// Plaintext formatted report
	fmt.Println(fmt.Sprintf("Success rate:      %d/%d", success, success+fail))
	fmt.Println(fmt.Sprintf("Average TTFB:      %dms", averageTTFB))
}
