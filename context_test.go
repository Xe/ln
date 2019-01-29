package ln

import (
	"context"
	"testing"
)

func TestContextStorage(t *testing.T) {
	var val interface{} = "bar"
	ctx := context.Background()
	ctx = WithF(ctx, F{"foo": val})

	f, ok := FFromContext(ctx)
	if !ok {
		t.Fatal("expected F to be in context but it wasn't")
	}

	if cmp, ok := f["foo"]; ok {
		if cmp != val {
			t.Fatalf("expected %v from context, got: %v", val, cmp)
		}
	}
}
