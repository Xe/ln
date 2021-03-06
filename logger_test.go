package ln

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"within.website/ln/opname"
)

var ctx context.Context

func setup(t *testing.T) (*bytes.Buffer, func()) {
	ctx = context.Background()

	out := bytes.Buffer{}
	oldFilters := DefaultLogger.Filters
	DefaultLogger.Filters = []Filter{
		NewWriterFilter(&out, nil),
	}
	return &out, func() {
		DefaultLogger.Filters = oldFilters
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

func TestSimpleError(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	Log(ctx, F{"err": fmt.Errorf("This is an Error!!!")}, F{"msg": "fooey", "bar": "foo"})
	data := []string{
		`err="This is an Error!!!"`,
		`fooey`,
		`bar=foo`,
	}

	must(t, out, data)
}

func TestTimeConversion(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	var zeroTime time.Time

	Log(ctx, F{"zero": zeroTime})
	data := []string{
		`zero=0001-01-01T00:00:00Z`,
	}

	must(t, out, data)
}

func TestDebug(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	mctx := opname.With(ctx, "test")

	// set priority to Debug
	Error(mctx, fmt.Errorf("This is an Error!!!"), F{})

	data := []string{
		`err="This is an Error!!!"`,
		`_lineno=`,
		`_function=ln.TestDebug`,
		`ln/logger_test.go`,
		`operation=test`,
	}

	must(t, out, data)
}

func TestFer(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	underTest := foobar{Foo: 1, Bar: "quux"}

	Log(ctx, underTest)
	data := []string{
		`foo=1`,
		`bar=quux`,
	}

	must(t, out, data)
}

type foobar struct {
	Foo int
	Bar string
}

func (f foobar) F() F {
	return F{
		"foo": f.Foo,
		"bar": f.Bar,
	}
}
