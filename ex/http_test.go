package ex

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"within.website/ln"
)

func setup() (context.Context, *bytes.Buffer, func()) {
	ctx, cancel := context.WithCancel(context.Background())

	out := bytes.Buffer{}
	oldFilters := ln.DefaultLogger.Filters
	ln.DefaultLogger.Filters = []ln.Filter{
		ln.NewWriterFilter(&out, nil),
	}
	return ctx, &out, func() {
		cancel()
		ln.DefaultLogger.Filters = oldFilters
	}
}

func must(t *testing.T, out *bytes.Buffer, data []string) {
	t.Helper()

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

func TestHTTPLog(t *testing.T) {
	ctx, out, teardown := setup()
	defer teardown()

	ts := httptest.NewServer(HTTPLog(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})))

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusTeapot {
		t.Fatalf("wanted status code %d, got: %d", http.StatusTeapot, resp.StatusCode)
	}

	time.Sleep(time.Millisecond)

	data := []string{
		`status=418`,
		`remote_ip=127.0.0.1`,
		`path=/`,
		`request_duration`,
	}

	must(t, out, data)
}
