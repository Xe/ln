package ex

import (
	"net"
	"net/http"
	"time"

	"within.website/ln"
	"within.website/ln/opname"
)

func HTTPLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		f := ln.F{
			"remote_ip":       host,
			"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
			"path":            r.URL.Path,
		}
		ctx := ln.WithF(r.Context(), f)
		st := time.Now()
		ctx = opname.With(ctx, "http")

		next.ServeHTTP(w, r.WithContext(ctx))

		af := time.Now()
		f["request_duration"] = af.Sub(st)

		ws, ok := w.(interface {
			Status() int
		})
		if ok {
			f["status"] = ws.Status()
		}

		ln.Log(ctx, f)
	})
}
