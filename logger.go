package ln

import (
	"os"
	"time"

	"github.com/pkg/errors"
)

// Logger holds the current priority and list of filters
type Logger struct {
	Filters []Filter
}

// DefaultLogger is the default implementation of Logger
var DefaultLogger *Logger

func init() {
	var defaultFilters []Filter

	// Default to STDOUT for logging, but allow LN_OUT to change it.
	out := os.Stdout
	if os.Getenv("LN_OUT") == "<stderr>" {
		out = os.Stderr
	}

	defaultFilters = append(defaultFilters, NewWriterFilter(out, nil))

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
func (l *Logger) Log(xs ...Fer) {
	event := Event{Time: time.Now()}

	addF := func(bf F) {
		if event.Data == nil {
			event.Data = bf
		} else {
			for k, v := range bf {
				event.Data[k] = v
			}
		}
	}

	for _, f := range xs {
		addF(f.F())
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

	l.filter(event)
}

func (l *Logger) filter(e Event) {
	for _, f := range l.Filters {
		if !f.Apply(e) {
			return
		}
	}
}

// Error logs an error and information about the context of said error.
func (l *Logger) Error(err error, xs ...Fer) {
	data := F{}
	frame := callersFrame()

	data["_lineno"] = frame.lineno
	data["_function"] = frame.function
	data["_filename"] = frame.filename
	data["err"] = err

	cause := errors.Cause(err)
	if cause != nil {
		data["cause"] = cause.Error()
	}

	xs = append(xs, data)

	l.Log(xs...)
}

// Fatal logs this set of values, then exits with status code 1.
func (l *Logger) Fatal(xs ...Fer) {
	xs = append(xs, F{"fatal": true})

	l.Log(xs...)

	os.Exit(1)
}

// FatalErr combines Fatal and Error.
func (l *Logger) FatalErr(err error, xs ...Fer) {
	xs = append(xs, F{"fatal": true})

	data := F{}
	frame := callersFrame()

	data["_lineno"] = frame.lineno
	data["_function"] = frame.function
	data["_filename"] = frame.filename
	data["err"] = err

	cause := errors.Cause(err)
	if cause != nil {
		data["cause"] = cause.Error()
	}

	xs = append(xs, data)
	l.Log(xs...)

	os.Exit(1)
}

// Default Implementation

// Log is the generic logging method.
func Log(xs ...Fer) {
	DefaultLogger.Log(xs...)
}

// Error logs an error and information about the context of said error.
func Error(err error, xs ...Fer) {
	DefaultLogger.Error(err, xs...)
}

// Fatal logs this set of values, then exits with status code 1.
func Fatal(xs ...Fer) {
	DefaultLogger.Fatal(xs...)
}

// FatalErr combines Fatal and Error.
func FatalErr(err error, xs ...Fer) {
	DefaultLogger.FatalErr(err, xs...)
}
