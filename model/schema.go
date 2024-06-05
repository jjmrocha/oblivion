package model

type Field struct {
	Name     string   `json:"field"`
	Type     DataType `json:"type"`
	Required bool     `json:"not-null"`
	Indexed  bool     `json:"indexed"`
}
