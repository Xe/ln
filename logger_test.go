package ln

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func setup(t *testing.T) (*bytes.Buffer, func()) {
	out := bytes.Buffer{}
	oldFilters := DefaultLogger.Filters
	DefaultLogger.Filters = []Filter{NewWriterFilter(&out, nil)}
	return &out, func() {
		DefaultLogger.Filters = oldFilters
	}
}

func TestSimpleError(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

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

func TestTimeConversion(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	var zeroTime time.Time

	Info(F{"zero": zeroTime})
	data := []string{
		`zero=0001-01-01T00:00:00Z`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

func TestDebug(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	oldPri := DefaultLogger.Pri
	defer func() { DefaultLogger.Pri = oldPri }()


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

func TestFer(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	underTest := foobar{Foo: 1, Bar: "quux"}

	Info(underTest)
	data := []string{
		`foo=1`,
		`bar=quux`,
	}

	for _, line := range data {
		if !bytes.Contains(out.Bytes(), []byte(line)) {
			t.Fatalf("Bytes: %s not in %s", line, out.Bytes())
		}
	}
}

type foobar struct {
	Foo int
	Bar string
}

func (f foobar) F() map[string]interface{} {
	return map[string]interface{} {
		"foo": f.Foo,
		"bar": f.Bar,
	}
}
