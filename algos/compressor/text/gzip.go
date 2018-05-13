package textcompressor

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

const (
	NoCompression      = gzip.NoCompression
	BestSpeed          = gzip.BestSpeed
	BestCompression    = gzip.BestCompression
	DefaultCompression = gzip.DefaultCompression
	HuffmanOnly        = gzip.HuffmanOnly
)

type GzipCompressor struct {
	Level int
}

func (c GzipCompressor) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, c.Level)
	if err != nil {
		return nil, err
	}
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c GzipCompressor) Decompress(data []byte) ([]byte, error) {
	rdata := bytes.NewReader(data)
	reader, err := gzip.NewReader(rdata)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reader)
}
