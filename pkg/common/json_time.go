package common

import (
	"time"
)

// JSONTimeDefaultLayout default format layout
const JSONTimeDefaultLayout = "2006-01-02T15:04:05.000000Z"

var jsonTimeDefaultMarshalFormat = JSONTimeDefaultLayout
var jsonTimeDefaultUnmarshalFormat = JSONTimeDefaultLayout

// JSONTime time.Time for JSON
type JSONTime struct {
	time.Time
}

// GetDefaultMarshalFormat get default format
func (t *JSONTime) GetDefaultMarshalFormat() string {
	return jsonTimeDefaultMarshalFormat
}

// GetDefaultUnmarshalFormat get default format
func (t *JSONTime) GetDefaultUnmarshalFormat() string {
	return jsonTimeDefaultUnmarshalFormat
}

// SetDefaultMarshalFormat set default format
func (t *JSONTime) SetDefaultMarshalFormat(format string) {
	jsonTimeDefaultMarshalFormat = format
}

// SetDefaultUnmarshalFormat set default format
func (t *JSONTime) SetDefaultUnmarshalFormat(format string) {
	jsonTimeDefaultUnmarshalFormat = format
}

// SetDefaultMarshalFormatAsDefault set default format as default
func (t *JSONTime) SetDefaultMarshalFormatAsDefault() {
	jsonTimeDefaultMarshalFormat = JSONTimeDefaultLayout
}

// SetDefaultUnmarshalFormatAsDefault set default format as default
func (t *JSONTime) SetDefaultUnmarshalFormatAsDefault() {
	jsonTimeDefaultUnmarshalFormat = JSONTimeDefaultLayout
}

// UnmarshalJSON for JSON
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	parsed, err := time.Parse(`"`+jsonTimeDefaultUnmarshalFormat+`"`, string(data))
	t.Time = parsed

	return err
}

// MarshalJSON for JSON
func (t *JSONTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	s := t.UTC().Format(jsonTimeDefaultMarshalFormat)
	return []byte(`"` + s + `"`), nil
}
