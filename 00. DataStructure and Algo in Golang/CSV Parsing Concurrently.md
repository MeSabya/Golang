## Parallel CSV parsing in Go, extracting specific columns and writing filtered results concurrently.

- Efficient parallel processing using goroutines
- Worker pool for parsing chunks of CSV
- Channel-based communication for concurrent writing
- Easy to modify: just change selected columns or filter logic

```go
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

const (
	inputFile        = "input.csv"
	outputFile       = "output.csv"
	numWorkers       = 4
	chunkSize        = 1000
)

var (
	columnsToExtract = []string{"Name", "Age"} // target columns
)

func main() {
	in, err := os.Open(inputFile)
	checkErr(err)
	defer in.Close()

	reader := csv.NewReader(in)

	header, err := reader.Read()
	checkErr(err)

	colIndices := findColumnIndices(header, columnsToExtract)

	// Channels
	linesCh := make(chan [][]string, numWorkers)
	resultCh := make(chan []map[string]string, numWorkers)

	var wg sync.WaitGroup

	// Writer goroutine
	wg.Add(1)
	go writer(resultCh, &wg)

	// Worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(linesCh, resultCh, colIndices, &wg)
	}

	// Read and distribute lines
	go func() {
		chunk := make([][]string, 0, chunkSize)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			checkErr(err)
			chunk = append(chunk, record)
			if len(chunk) >= chunkSize {
				linesCh <- chunk
				chunk = make([][]string, 0, chunkSize)
			}
		}
		if len(chunk) > 0 {
			linesCh <- chunk
		}
		close(linesCh)
	}()

	wg.Wait()
}

func findColumnIndices(header, targets []string) map[string]int {
	indexMap := make(map[string]int)
	for _, target := range targets {
		found := false
		for i, col := range header {
			if col == target {
				indexMap[target] = i
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("Column %s not found in header", target))
		}
	}
	return indexMap
}

func worker(linesCh <-chan [][]string, resultCh chan<- []map[string]string, colIndices map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	for lines := range linesCh {
		var results []map[string]string
		for _, row := range lines {
			if passesFilter(row, colIndices) {
				result := make(map[string]string)
				for col, idx := range colIndices {
					result[col] = row[idx]
				}
				results = append(results, result)
			}
		}
		if len(results) > 0 {
			resultCh <- results
		}
	}
}

func passesFilter(row []string, colIndices map[string]int) bool {
	// Example: filter rows where Age > 30
	ageStr := row[colIndices["Age"]]
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return false
	}
	return age > 30
}

func writer(resultCh <-chan []map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	out, err := os.Create(outputFile)
	checkErr(err)
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()

	// Write header
	writer.Write(columnsToExtract)

	for results := range resultCh {
		for _, row := range results {
			record := make([]string, len(columnsToExtract))
			for i, col := range columnsToExtract {
				record[i] = row[col]
			}
			writer.Write(record)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
```

### Example CSV (input.csv):

```csv
ID,Name,Age,Country
1,Alice,25,USA
2,Bob,35,Canada
3,Carol,40,India
```

### Output (output.csv):
```csv
Name,Age
Bob,35
Carol,40
```

