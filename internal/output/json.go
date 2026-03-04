package output

import "encoding/json"

// JSONFormatter outputs data as indented JSON.
type JSONFormatter struct{}

// Format marshals data to indented JSON.
func (f *JSONFormatter) Format(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
