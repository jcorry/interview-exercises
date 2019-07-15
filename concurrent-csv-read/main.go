package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var csvPath string

func main() {
	flag.StringVar(&csvPath, "dir", "./csv-files", "The directory in which to look for CSV files")
	flag.Parse()

	fmt.Printf("Dir: %+q\n", csvPath)
	for _, f := range getCsvFiles(csvPath) {
		fmt.Println(fmt.Sprintf("CSV file: %s", f))
	}
}

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
