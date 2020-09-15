package middleware

import (
	"net/http"
	"runtime/debug"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func (mware *HTTPMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				mware.log.Log(
					"err", err,
					"trace", debug.Stack(),
				)
			}
		}()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		mware.log.Log(
			"status", wrapped.status,
			"method", r.Method,
			"path", r.URL.EscapedPath(),
			"duration", time.Since(start),
		)
	})
}
