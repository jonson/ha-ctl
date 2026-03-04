package output

// Formatter formats data for output.
type Formatter interface {
	Format(data interface{}) (string, error)
}

// New returns a formatter for the given format name ("json" or "text").
func New(format string) Formatter {
	switch format {
	case "text":
		return &TextFormatter{}
	default:
		return &JSONFormatter{}
	}
}
