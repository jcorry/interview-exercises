package main

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func ReaderToByte(r io.Reader) []byte {
	b, _ := ioutil.ReadAll(r)
	return b
}

func TestGzipper_StreamGzip(t *testing.T) {
	type args struct {
		s io.Reader
	}
	tests := []struct {
		name string
		args args
		want io.Reader
	}{
		{
			"success",
			args{
				strings.NewReader("My test compression string"),
			},
			strings.NewReader("My test compression string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gzipper{}

			got := g.StreamGzip(tt.args.s)
			unzipped, err := g.StreamGunzip(got)

			if err != nil {
				t.Errorf("Error unzipping data")
			}

			bytes, _ := ioutil.ReadAll(tt.args.s)

			if string(ReaderToByte(unzipped)) != string(bytes) {
				t.Errorf("Unzipped data does not match zipped data")
			}
		})
	}
}
