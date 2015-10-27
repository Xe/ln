package ln

import (
	"fmt"
	"os"
	"time"
)

// Priority represents the importance of an Event.
type Priority int

const (
	// PriEmergency is the emergency priority, and is the highest.
	PriEmergency Priority = iota
	// PriAlert is the alert priority
	PriAlert
	// PriCritical is the critical priority
	PriCritical
	// PriError is the error priority.
	PriError
	// PriWarning is the warning priority
	PriWarning
	// PriNotice is the notice priority
	PriNotice
	// PriInfo is the info priority
	PriInfo
	// PriDebug is the debug priority.
	PriDebug
)

var priStrings = []string{
	"emerg",
	"alert",
	"crit",
	"error",
	"warn",
	"notice",
	"info",
	"debug",
}

func (p Priority) String() string {
	if int(p) < len(priStrings) {
		return priStrings[p]
	}

	return "UNKNOWN"
}

// Logger holds the current priority and list of filters
type Logger struct {
	Pri     Priority
	Filters []Filter
}

// DefaultLogger is the default implementation of Logger
var DefaultLogger *Logger

func init() {
	var defaultFilters []Filter

	defaultFilters = append(defaultFilters, NewWriterFilter(os.Stdout, nil))

	DefaultLogger = &Logger{
		Pri:     PriInfo,
		Filters: defaultFilters,
	}

}

// F is a key-value mapping for structured data.
type F map[string]interface{}

// Event represents an event
type Event struct {
	Pri     Priority
	Time    time.Time
	Data    F
	Message string
}

// Log is the generic logging method.
func (l *Logger) Log(p Priority, xs ...interface{}) {
	if l.Pri < p {
		return // don't log
	}

	var bits []interface{}
	event := Event{Pri: p, Time: time.Now()}

	// Assemble the event
	for _, b := range xs {
		switch b.(type) {
		case F:
			bf := b.(F)
			if event.Data == nil {
				event.Data = bf
			} else {
				for k, v := range bf {
					event.Data[k] = v
				}
			}
		default:
			bits = append(bits, b)
		}
	}

	event.Message = fmt.Sprint(bits...)

	l.filter(event)
}

func (l *Logger) filter(e Event) {
	for _, f := range l.Filters {
		if !f.Apply(e) {
			return
		}
	}
}

// Emergency sets the priority of this event to PriEmergency
func (l *Logger) Emergency(xs ...interface{}) {
	l.Log(PriEmergency, xs...)
}

// Alert sets the priority of this event to PriAlert
func (l *Logger) Alert(xs ...interface{}) {
	l.Log(PriAlert, xs...)
}

// Critical sets the priority of this event to PriCritical
func (l *Logger) Critical(xs ...interface{}) {
	l.Log(PriCritical, xs...)
}

// Error sets the priority of this event to PriError
func (l *Logger) Error(xs ...interface{}) {
	l.Log(PriError, xs...)
}

// Warning sets the priority of this event to PriWarning
func (l *Logger) Warning(xs ...interface{}) {
	l.Log(PriWarning, xs...)
}

// Notice sets the priority of this event to PriNotice
func (l *Logger) Notice(xs ...interface{}) {
	l.Log(PriNotice, xs...)
}

// Info sets the priority of this event to PriInfo
func (l *Logger) Info(xs ...interface{}) {
	l.Log(PriInfo, xs...)
}

// Debug sets the priority of this event to PriDebug
func (l *Logger) Debug(xs ...interface{}) {
	l.Log(PriDebug, xs...)
}

// Default Implementation

// Emergency sets the priority of this event to PriEmergency
func Emergency(xs ...interface{}) {
	DefaultLogger.Log(PriEmergency, xs...)
}

// Alert sets the priority of this event to PriAlert
func Alert(xs ...interface{}) {
	DefaultLogger.Log(PriAlert, xs...)
}

// Critical sets the priority of this event to PriCritical
func Critical(xs ...interface{}) {
	DefaultLogger.Log(PriCritical, xs...)
}

// Error sets the priority of this event to PriError
func Error(xs ...interface{}) {
	DefaultLogger.Log(PriError, xs...)
}

// Warning sets the priority of this event to PriWarning
func Warning(xs ...interface{}) {
	DefaultLogger.Log(PriWarning, xs...)
}

// Notice sets the priority of this event to PriNotice
func Notice(xs ...interface{}) {
	DefaultLogger.Log(PriNotice, xs...)
}

// Info sets the priority of this event to PriInfo
func Info(xs ...interface{}) {
	DefaultLogger.Log(PriInfo, xs...)
}

// Debug sets the priority of this event to PriDebug
func Debug(xs ...interface{}) {
	DefaultLogger.Log(PriDebug, xs...)
}
