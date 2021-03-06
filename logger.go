package ln

import (
	"context"
	"io"
	"log/syslog"
	"os"
	"time"

	"github.com/pkg/errors"
)

// Logger holds the current priority and list of filters
type Logger struct {
	Filters []Filter
}

// AddFilter adds a filter to the beginning of the stack.
func (l *Logger) AddFilter(f Filter) {
	l.Filters = append([]Filter{f}, l.Filters...)
}

// AddFilter adds a filter to the beginning of the default logger stack.
func AddFilter(f Filter) {
	DefaultLogger.AddFilter(f)
}

// DefaultLogger is the default implementation of Logger
var DefaultLogger *Logger

func init() {
	var defaultFilters []Filter

	// Default to STDOUT for logging, but allow LN_OUT to change it.
	var out io.Writer = os.Stdout
	switch os.Getenv("LN_OUT") {
	case "<stderr>":
		out = os.Stderr
	case "<syslog>":
		wr, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_USER, "")
		if err != nil {
			panic(err)
		}

		out = wr
	}

	var formatter Formatter
	switch os.Getenv("LN_FORMATTER") {
	case "text":
		formatter = NewTextFormatter()
	default:
		formatter = JSONFormatter()
	}

	defaultFilters = append(
		defaultFilters,
		NewWriterFilter(out, formatter),
	)

	DefaultLogger = &Logger{
		Filters: defaultFilters,
	}
}

// F is a key-value mapping for structured data.
type F map[string]interface{}

// Extend concatentates one F with one or many Fer instances.
func (f F) Extend(other ...Fer) {
	for _, ff := range other {
		for k, v := range ff.F() {
			f[k] = v
		}
	}
}

// F makes F an Fer
func (f F) F() F {
	return f
}

// Fer allows any type to add fields to the structured logging key->value pairs.
type Fer interface {
	F() F
}

// Event represents an event
type Event struct {
	Time    time.Time
	Data    F
	Message string
}

// Log is the generic logging method.
func (l *Logger) Log(ctx context.Context, xs ...Fer) {
	event := Event{
		Time: time.Now(),
		Data: F{},
	}

	addF := func(bf F) {
		for k, v := range bf {
			event.Data[k] = v
		}
	}

	for _, f := range xs {
		addF(f.F())
	}

	ctxf, ok := FFromContext(ctx)
	if ok {
		addF(ctxf)
	}

	if os.Getenv("LN_DEBUG_ALL_EVENTS") == "1" {
		frame := callersFrame()
		if event.Data == nil {
			event.Data = make(F)
		}
		event.Data["_lineno"] = frame.lineno
		event.Data["_function"] = frame.function
		event.Data["_filename"] = frame.filename
	}

	l.filter(ctx, event)
}

func (l *Logger) filter(ctx context.Context, e Event) {
	for _, f := range l.Filters {
		if !f.Apply(ctx, e) {
			return
		}
	}
}

// Error logs an error and information about the context of said error.
func (l *Logger) Error(ctx context.Context, err error, xs ...Fer) {
	data := F{}
	frame := callersFrame()

	data["_lineno"] = frame.lineno
	data["_function"] = frame.function
	data["_filename"] = frame.filename
	data["err"] = err

	if fer, ok := err.(Fer); ok {
		data.Extend(fer)
	}

	cause := errors.Cause(err)
	if cause != nil && cause.Error() != err.Error() {
		data["cause"] = cause.Error()
	}
	doGoError(err, data)

	xs = append(xs, data)

	l.Log(ctx, xs...)
}

// Fatal logs this set of values, then panics.
func (l *Logger) Fatal(ctx context.Context, xs ...Fer) {
	xs = append(xs, F{"fatal": true})

	l.Log(ctx, xs...)

	panic("ln.Fatal called")
}

// FatalErr combines Fatal and Error.
func (l *Logger) FatalErr(ctx context.Context, err error, xs ...Fer) {
	xs = append(xs, F{"fatal": true})

	data := F{}
	frame := callersFrame()

	data["_lineno"] = frame.lineno
	data["_function"] = frame.function
	data["_filename"] = frame.filename
	data["err"] = err

	cause := errors.Cause(err)
	if cause != nil && cause.Error() != err.Error() {
		data["cause"] = cause.Error()
	}
	doGoError(err, data)

	xs = append(xs, data)
	l.Log(ctx, xs...)

	panic("ln.FatalErr called")
}

// Default Implementation

// Log is the generic logging method.
func Log(ctx context.Context, xs ...Fer) {
	DefaultLogger.Log(ctx, xs...)
}

// Error logs an error and information about the context of said error.
func Error(ctx context.Context, err error, xs ...Fer) {
	DefaultLogger.Error(ctx, err, xs...)
}

// Fatal logs this set of values, then panics.
func Fatal(ctx context.Context, xs ...Fer) {
	DefaultLogger.Fatal(ctx, xs...)
}

// FatalErr combines Fatal and Error.
func FatalErr(ctx context.Context, err error, xs ...Fer) {
	DefaultLogger.FatalErr(ctx, err, xs...)
}
