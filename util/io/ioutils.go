// Author: ZHU HAIHUA
// Date: 8/15/16
package ioutils

import (
	"compress/gzip"
	"io"
)

func ZipReader(encoding string, reader io.ReadCloser) (io.ReadCloser, error) {
	var r io.ReadCloser
	var err error
	switch encoding {
	case "gzip":
		r, err = gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
	default:
		r = reader
	}
	return r, nil
}

func ZipWriter(encoding string, writer io.Writer) io.Writer {
	var w io.Writer
	switch encoding {
	case "gzip":
		w = gzip.NewWriter(writer)
	default:
		w = writer
	}
	return w
}
