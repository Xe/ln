package ln

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSimpleError(t *testing.T) {
	out := bytes.Buffer{}
	oldFilters := DefaultLogger.Filters
	DefaultLogger.Filters = []Filter{NewWriterFilter(&out, nil)}
	defer func() {
		DefaultLogger.Filters = oldFilters
	}()
	Info(F{"err": fmt.Errorf("This is an Error!!!")}, "fooey", F{"bar": "foo"})
	data := []string{
		`err="This is an Error!!!"`,
		`fooey`,
		`bar=foo`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

func TestDebug(t *testing.T) {
	oldPri := DefaultLogger.Pri
	defer func() { DefaultLogger.Pri = oldPri }()

	out := bytes.Buffer{}
	oldFilters := DefaultLogger.Filters
	DefaultLogger.Filters = []Filter{NewWriterFilter(&out, nil)}
	defer func() {
		DefaultLogger.Filters = oldFilters
	}()

	// set priority to Debug
	DefaultLogger.Pri = PriDebug
	Debug(F{"err": fmt.Errorf("This is an Error!!!")})

	data := []string{
		`err="This is an Error!!!"`,
		`_lineno=`,
		`_function=ln.TestDebug`,
		`_filename=github.com/apg/ln/logger_test.go`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}
