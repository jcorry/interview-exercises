package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// csvPath is a string to hold the path to directory containing CSV files
var csvPath string

// codes is the map that will store all valid codes.
var codes map[string]bool

// Lock the map for reading/writing
var lock = sync.RWMutex{}

func main() {
	codes = make(map[string]bool)

	// Set `csvPath` to the arg provided on app startup
	flag.StringVar(&csvPath, "dir", "./csv-files", "The directory in which to look for CSV files")
	flag.Parse()

	fmt.Printf("Dir: %+q\n", csvPath)
	files := getCsvFiles(csvPath)

	duration := time.Duration(500 * time.Millisecond)

	ctx, _ := context.WithTimeout(context.Background(), duration)

	fmt.Println(checkFilesForDuplicateCodes(ctx, files))
}

func checkFilesForDuplicateCodes(ctx context.Context, filenames []string) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, f := range filenames {
		f := f
		g.Go(func() error {
			fmt.Println(fmt.Sprintf("Processing CSV file: %s", f))

			// Handle file processing here, if you find a duplicate code emit an error on the errc channel
			// otherwise emit nil on done
			file, err := os.Open(f)
			if err != nil {
				return err
			}
			defer file.Close()

			reader := csv.NewReader(bufio.NewReader(file))
			lineNum := 1
			for {
				line, err := reader.Read()
				if err == io.EOF {
					return err
				} else if err != nil {
					return err
				}

				err = check(line[1])

				if err != nil {
					return fmt.Errorf("Duplicate code found on line %d of : %s", lineNum, f)
				}
				lineNum++
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			}
			return nil

		})
	}
	return g.Wait()
}

// check for the code alread existing in the codes map
func check(code string) error {
	// Don't check header rows
	if code == "code" {
		return nil
	}
	lock.RLock()
	_, exists := codes[code]
	lock.RUnlock()

	if exists {
		return errors.New("Duplicate key found")
	}

	lock.Lock()
	codes[code] = false
	lock.Unlock()

	return nil
}

// getCsvFiles returns a slice of strings holding the filenames of all of the csv files
// in csvPath
func getCsvFiles(dir string) []string {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".csv") {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
