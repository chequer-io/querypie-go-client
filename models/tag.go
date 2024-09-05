package models

type Tag struct {
	Key          string `json:"key"`
	Synchronized bool   `json:"synchronized"`
	Value        string `json:"value"`
}
