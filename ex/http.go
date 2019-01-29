package ex

import (
	"net"
	"net/http"
	"time"

	"within.website/ln"
)

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// HTTPLog automagically logs HTTP traffic.
func HTTPLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		f := ln.F{
			"remote_ip": host,
			"path":      r.URL.Path,
		}

		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			f["x_forwarded_for"] = xff
		}

		ctx := ln.WithF(r.Context(), f)
		st := time.Now()
		srw := &statusResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(srw, r.WithContext(ctx))

		f["request_duration"] = time.Since(st)

		if srw.status != 0 {
			f["status"] = srw.status
		}

		ln.Log(ctx, f)
	})
}
