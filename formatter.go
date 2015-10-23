package saw

import (
	"bytes"
	"fmt"
	"time"
)

var (

	// DefaultTimeFormat represents the way in which time will be formatted by default
	DefaultTimeFormat = time.RFC3339
)

// Formatter defines the formatting of events
type Formatter interface {
	Format(Event) ([]byte, error)
}

// DefaultFormatter is the default way in which to format events
var DefaultFormatter Formatter

func init() {
	DefaultFormatter = NewTextFormatter()
}

type TextFormatter struct {
	TimeFormat string
}

// NewTextFormatter returns a Formatter that outputs as text.
func NewTextFormatter() Formatter {
	return &TextFormatter{TimeFormat: DefaultTimeFormat}
}

func (t *TextFormatter) Format(e Event) ([]byte, error) {
	var writer bytes.Buffer

	writer.WriteString(e.Time.Format(t.TimeFormat))
	writer.WriteByte(' ')

	writer.WriteString("level=")
	writer.WriteString(e.Pri.String())

	for k, v := range e.Data {
		writer.WriteByte(' ')
		if shouldQuote(k) {
			writer.WriteString(fmt.Sprintf("%q", k))
		} else {
			writer.WriteString(k)
		}

		writer.WriteByte('=')

		switch v.(type) {
		case string:
			vs, _ := v.(string)
			if shouldQuote(vs) {
				writer.WriteString(fmt.Sprintf("%q", vs))
			} else {
				writer.WriteString(vs)
			}
		case error:
			tmperr, _ := v.(error)
			es := tmperr.Error()

			if shouldQuote(es) {
				writer.WriteString(fmt.Sprintf("%q", es))
			} else {
				writer.WriteString(es)
			}
		default:
			fmt.Fprint(&writer, v)
		}
	}

	if len(e.Message) > 0 {
		writer.WriteByte(' ')
		writer.WriteString(e.Message)
	}

	writer.WriteByte('\n')

	return writer.Bytes(), nil
}

func shouldQuote(s string) bool {
	for _, b := range s {
		if !((b >= 'A' && b <= 'Z') ||
			(b >= 'a' && b <= 'z') ||
			(b >= '0' && b <= '9') ||
			(b == '-' || b == '.' || b == '#' ||
				b == '/')) {
			return true
		}
	}
	return false
}
