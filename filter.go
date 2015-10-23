package ln

import "io"

// Filter interface for defining chain filters
type Filter interface {
	Apply(Event) bool
	Run()
	Close()
}

// FilterFunc allows simple functions to implement the Filter interface
type FilterFunc func(e Event) bool

// Apply implements the Filter interface
func (ff FilterFunc) Apply(e Event) bool {
	return ff(e)
}

// Run implements the Filter interface
func (ff *FilterFunc) Run() {}

// Close implements the Filter interface
func (ff *FilterFunc) Close() {}

// WriterFilter implements a filter, which arbitrarily writes to an io.Writer
type WriterFilter struct {
	Out       io.Writer
	Formatter Formatter
}

// NewWriterFilter creates a filter to add to the chain
func NewWriterFilter(out io.Writer, format Formatter) *WriterFilter {
	if format == nil {
		format = DefaultFormatter
	}
	return &WriterFilter{
		Out:       out,
		Formatter: format,
	}
}

// Apply implements the Filter interface
func (w *WriterFilter) Apply(e Event) bool {
	output, err := w.Formatter.Format(e)
	if err == nil {
		w.Out.Write(output)
	}

	return true
}

// Apply implements the Filter interface
func (w *WriterFilter) Run() {}

// Apply implements the Filter interface
func (w *WriterFilter) Close() {}
