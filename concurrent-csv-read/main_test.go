package main

import "testing"

func Test_GetCsvFiles(t *testing.T) {

	tests := []struct {
		name      string
		path      string
		fileCount int
		wantErr   error
	}{
		{
			"Success",
			"testdata",
			1,
			nil,
		},
		{
			"Dir not found",
			"badpath",
			0,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := getCsvFiles(tt.path)

			if len(files) != tt.fileCount {
				t.Errorf("Unexpected file count; want %d; got %d", tt.fileCount, len(files))
			}
		})
	}

}
