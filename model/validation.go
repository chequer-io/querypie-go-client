package model

type Validation struct {
	Result bool     `json:"-" yaml:"result"`
	Reason []string `json:"-" yaml:"reason,omitempty"`
}
