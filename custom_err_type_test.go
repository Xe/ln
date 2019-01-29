package ln

import (
	"context"
	"fmt"
	"testing"
)

type customError struct {
	statusCode int
	message    string
}

func (e customError) Error() string {
	return fmt.Sprintf("customError: got status code %d and message %q", e.statusCode, e.message)
}

func (e customError) F() F {
	return F{
		"err_status_code": e.statusCode,
		"err_message":     e.message,
	}
}

func TestCustomErrorType(t *testing.T) {
	out, teardown := setup(t)
	defer teardown()

	ctx := context.Background()

	err := &customError{
		statusCode: 420,
		message:    "success",
	}

	Error(ctx, err)
	data := []string{
		`err_status_code=420`,
		`err_message=success`,
	}

	must(t, out, data)
}
