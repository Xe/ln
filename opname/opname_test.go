package opname

import (
	"context"
	"testing"
)

func TestOpname(t *testing.T) {
	ctx := context.Background()
	ctxBase := With(ctx, "base")
	ctxExtended := With(ctxBase, "extended")

	_, ok := Get(ctx)
	if ok {
		t.Fatal("shouldn't get an operation name")
	}
	
	val, ok := Get(ctxBase)
	if !ok {
		t.Fatal("ctxBase should have an operation name")
	}
	
	if val != "base" {
		t.Fatalf("wanted base, got: %v", val)
	}
	
	val, ok = Get(ctxExtended)
	if !ok {
		t.Fatal("ctxExtended should have an operation name")
	}
	
	if val != "base.extended" {
		t.Fatalf("wanted base.extended, got: %v", val)
	}
}