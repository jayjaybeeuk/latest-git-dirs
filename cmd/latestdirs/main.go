package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"

	"github.com/jamesburton/latestdirs/internal/model"
	"github.com/jamesburton/latestdirs/internal/output"
	"github.com/jamesburton/latestdirs/internal/scanner"
)

func main() {
	top := flag.Int("top", 20, "number of results")
	depth := flag.Int("max-depth", 10, "maximum depth")
	jsonOutput := flag.Bool("json", false, "json output")
	workers := flag.Int("workers", runtime.NumCPU(), "worker count")

	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	excludes := map[string]bool{
		".git":         true,
		"node_modules": true,
		"bin":          true,
		"obj":          true,
	}

	workerCount := *workers
	if workerCount <= 0 {
		workerCount = runtime.NumCPU()
		if workerCount <= 0 {
			workerCount = 1
		}
	}

	jobs := make(chan string)
	results := make(chan model.Result)

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for dir := range jobs {
				latest, err := scanner.LatestModified(dir, *depth, excludes)
				if err != nil {
					continue
				}
				if latest.IsZero() {
					continue
				}

				results <- model.Result{
					Path:      dir,
					Timestamp: latest,
					Unix:      latest.Unix(),
					Source:    "filesystem",
				}
			}
		}()
	}

	go func() {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if excludes[entry.Name()] {
				continue
			}

			path := filepath.Join(root, entry.Name())
			jobs <- path
		}

		close(jobs)
		wg.Wait()
		close(results)
	}()

	var collected []model.Result

	for r := range results {
		collected = append(collected, r)
	}

	sort.Slice(collected, func(i, j int) bool {
		return collected[i].Timestamp.After(collected[j].Timestamp)
	})

	if len(collected) > *top {
		collected = collected[:*top]
	}

	if *jsonOutput {
		_ = output.PrintJSON(collected)
	} else {
		output.PrintTable(collected)
	}
}
