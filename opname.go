package ln

import (
	"context"
)

func opnameInEvents(ctx context.Context, e Event) bool {
	if op, ok := opname.Get(ctx); ok {
		d.Data["operation"] = op
	}
	
	return true
}