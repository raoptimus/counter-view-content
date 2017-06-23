package main

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type CompressResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w CompressResponseWriter) Write(b []byte) (int, error) {
	s := string(b)
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\r", " ", -1)
	s = strings.Replace(s, "  ", " ", -1)
	return w.Writer.Write([]byte(s))
}

func newCompressHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ae := r.Header.Get("Accept-Encoding")
		switch {
		case strings.Contains(ae, "deflate"):
			{
				w.Header().Set("Content-Encoding", "deflate")
				w.Header().Set("Vary", "Accept-Encoding")
				gz, _ := flate.NewWriter(w, 7)
				defer gz.Close()
				gzr := CompressResponseWriter{Writer: gz, ResponseWriter: w}
				fn(gzr, r)
			}
		case strings.Contains(ae, "gzip"):
			{
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Set("Vary", "Accept-Encoding")
				gz, _ := gzip.NewWriterLevel(w, 7)
				defer gz.Close()
				gzr := CompressResponseWriter{Writer: gz, ResponseWriter: w}
				fn(gzr, r)
			}
		default:
			{
				fn(w, r)
				return
			}
		}
	}
}
