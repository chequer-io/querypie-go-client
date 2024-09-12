package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringList struct {
	values []string
}

func (s StringList) Ellipsis() string {
	if len(s.values) < 4 {
		return s.String()
	}
	var extractedList StringList
	extractedList.values = s.values[:3]
	extracted := extractedList.String()
	extracted = extracted[:len(extracted)-1]
	return fmt.Sprintf("%s,+%d]", extracted, len(s.values)-3)
}

// MarshalJSON implements the json.Marshaller interface
func (s StringList) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.values)
}

func (s *StringList) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.values); err != nil {
		return err
	}
	return nil
}

// Value Implement the Valuer interface for "database/sql/driver"
func (s StringList) Value() (driver.Value, error) {
	return json.Marshal(s.values)
}

// Scan Implement the Scanner interface for "database/sql/driver"
func (s *StringList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s.values)
}

// String implements the fmt.Stringer interface
// NOTE(JK): It should be a pointer receiver method, or
// it does not work on fmt.Printf("%s", s).
func (s StringList) String() string {
	_json, _ := json.Marshal(s.values)
	return string(_json)
}
