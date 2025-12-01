package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringArray represents an array of strings stored as PostgreSQL text[].
//
// swagger:model
type StringArray []string

// Scan converts the PostgreSQL text[] format ("{a,b,c}") into a Go []string.
func (s *StringArray) Scan(value any) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	var str string

	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("unsupported type for StringArray: %T", value)
	}

	str = strings.Trim(str, "{}")
	if str == "" {
		*s = []string{}
		return nil
	}

	*s = strings.Split(str, ",")
	return nil
}

// Value converts a Go []string into a PostgreSQL text[] format ("{a,b,c}").
func (s StringArray) Value() (driver.Value, error) {
	return "{" + strings.Join(s, ",") + "}", nil
}
