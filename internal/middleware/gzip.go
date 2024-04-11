package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/carry0987/FileTree-API/internal/utils"
)

// GzipResponseWriter is a custom http.ResponseWriter that compresses responses with GZip.
type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Compresses HTTP responses for clients that support it.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if utils.IsWebSocket(r) || !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := GzipResponseWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(gzw, r)
	})
}
