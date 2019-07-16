package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	checkFilesForDuplicateCodes(files)
}

func checkFilesForDuplicateCodes(files []string) error {
	// Channels for goroutine control
	quit := make(chan bool)
	errc := make(chan error)
	done := make(chan error)

	for _, f := range files {
		go func(filename string) {
			fmt.Println(fmt.Sprintf("Processing CSV file: %s", filename))

			// Handle file processing here, if you find a duplicate code emit an error on the errc channel
			// otherwise emit nil on done
			file, err := os.Open(filename)
			if err != nil {
				errc <- err
				return
			}

			reader := csv.NewReader(bufio.NewReader(file))
			for {
				line, err := reader.Read()
				if err == io.EOF {
					done <- nil
					return
				} else if err != nil {
					errc <- err
					return
				}

				err = check(line[1])

				if err != nil {
					errc <- fmt.Errorf("Duplicate code found in: %s", filename)
					return
				}
			}
		}(f)
	}

	count := 0

	for {
		select {
		case err := <-errc:
			fmt.Errorf("Error found: %s", err)
			close(quit)
			return err
		case <-done:
			count++
			if count == len(files) {
				fmt.Println("Done, no duplicates found")
				return nil
			}
		}
	}
}

// check for the code alread existing in the codes map
func check(code string) error {
	lock.RLock()
	_, exists := codes[code]
	lock.RUnlock()

	if exists {
		return errors.New("Duplicate key found")
	} else {
		lock.Lock()
		codes[code] = false
		lock.Unlock()
	}
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
