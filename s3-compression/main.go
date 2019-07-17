package main

import (
	"compress/gzip"
	"io"
)

type Compressor interface {
	StreamGzip(r io.Reader) io.Reader
}

type Decompressor interface {
	StreamGunzip(r io.Reader) io.Reader
}

type Gzipper struct {
}

// StreamGzip compresses the
func (g *Gzipper) StreamGzip(s io.Reader) io.Reader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		zip, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		defer zip.Close()

		if err != nil {
			w.CloseWithError(err)
		}

		io.Copy(zip, s)
	}()
	return r
}

func (g *Gzipper) StreamGunzip(r io.Reader) (io.Reader, error) {
	unzip, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return unzip, nil
}

func main() {

}
