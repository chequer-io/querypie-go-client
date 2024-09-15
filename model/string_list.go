package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringList struct {
	Values []string
}

func (s StringList) Ellipsis(show int) string {
	if len(s.Values) <= show {
		return s.String()
	}
	var extractedList StringList
	extractedList.Values = s.Values[:show]
	extracted := extractedList.String()
	extracted = extracted[:len(extracted)-1]
	return fmt.Sprintf("%s,+%d]", extracted, len(s.Values)-show)
}

// MarshalJSON implements the json.Marshaller interface
func (s StringList) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func (s *StringList) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.Values); err != nil {
		return err
	}
	return nil
}

// Value Implement the Valuer interface for "database/sql/driver"
func (s StringList) Value() (driver.Value, error) {
	return json.Marshal(s.Values)
}

// Scan Implement the Scanner interface for "database/sql/driver"
func (s *StringList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s.Values)
}

// String implements the fmt.Stringer interface
// NOTE(JK): It should be a pointer receiver method, or
// it does not work on fmt.Printf("%s", s).
func (s StringList) String() string {
	_json, _ := json.Marshal(s.Values)
	return string(_json)
}
