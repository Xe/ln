package ln

import (
	"errors"
	"testing"
)

func TestRegressionErrorPutsStuffInF(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	f := F{
		"this": "is_real",
	}

	l := new(Logger)
	l.Error(ctx, errors.New("boo"), f)

	Log(ctx, f)
	data := []string{`this=is_real`}

	must(t, out, data)
}
