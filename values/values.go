package values

import (
	"fmt"
	"strings"
	"time"
)

// DefaultTimestampFormat ...
var DefaultTimestampFormat = "2006-01-02 15:04:05"

// TimestampValue ...
type TimestampValue struct {
	Format  *string
	Default time.Time
	ts      *time.Time
}

func (ts *TimestampValue) getFormat() string {
	if ts.Format != nil {
		return *ts.Format
	}
	return DefaultTimestampFormat
}

// Set ...
func (ts *TimestampValue) Set(value string) error {
	t, err := time.Parse(ts.getFormat(), value)
	if err != nil {
		return fmt.Errorf("%s cannot be parsed as a timestamp (expected format: %s)", value, ts.getFormat())
	}
	ts.ts = &t
	return nil
}

func (ts TimestampValue) String() string {
	if ts.ts != nil {
		return ts.ts.Format(ts.getFormat())
	}
	return ts.Default.Format(ts.getFormat())
}

// EnumValue ...
type EnumValue struct {
	Enum      []string
	Default   string
	AllowNone bool
	set       bool
	selected  string
}

// Set ...
func (e *EnumValue) Set(value string) error {
	e.set = true
	value = strings.TrimSpace(strings.ToLower(value))
	for _, enum := range e.Enum {
		if strings.ToLower(enum) == value {
			e.selected = value
			return nil
		}
	}
	if !e.AllowNone {
		return fmt.Errorf("Unknown option: \"%s\". Allowed is one of %s", value, strings.Join(e.Enum, ", "))
	}
	return nil
}

func (e EnumValue) String() string {
	if !e.set {
		return e.Default
	}
	return e.selected
}

// EnumListValue ...
type EnumListValue struct {
	Enum       []string
	Default    []string
	AllowEmpty bool
	selected   []string
	set        bool
}

// Parse ...
func (e EnumListValue) Parse(value string) []string {
	split := strings.Split(value, ",")
	var valid []string
	for _, s := range split {
		if len(strings.TrimSpace(s)) < 1 {
			continue
		}
		valid = append(valid, s)
	}
	return valid
}

// Set ...
func (e *EnumListValue) Set(value string) error {
	e.set = true
	var validEnums []string
	enums := strings.Split(strings.ToLower(value), ",")
	for _, enum := range enums {
		enum = strings.TrimSpace(strings.ToLower(enum))
		if enum == "" {
			continue
		}
		valid := false
		for _, availableEnum := range e.Enum {
			if enum == strings.TrimSpace(strings.ToLower(availableEnum)) {
				valid = true
			}
		}
		if !valid {
			return fmt.Errorf("Unknown option: \"%s\". Allowed values are %s", enum, strings.Join(e.Enum, ", "))
		}
		validEnums = append(validEnums, enum)
	}

	if !e.AllowEmpty && len(validEnums) < 1 {
		return fmt.Errorf("Must specify at least one of: %s", strings.Join(e.Enum, ", "))
	}
	e.selected = validEnums
	return nil
}

func (e EnumListValue) serialize(values []string) string {
	for i := range values {
		values[i] = strings.ToLower(strings.TrimSpace(values[i]))
	}
	return strings.Join(values, ",")
}

func (e EnumListValue) String() string {
	if !e.set {
		return e.serialize(e.Default)
	}
	return e.serialize(e.selected)
}
