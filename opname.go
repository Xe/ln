package ln

import (
	"context"

	"github.com/Xe/ln/opname"
)

func opnameInEvents(ctx context.Context, e Event) bool {
	if op, ok := opname.Get(ctx); ok {
		e.Data["operation"] = op
	}

	return true
}
