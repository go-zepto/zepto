package fields

type FieldOptions = map[string]interface{}

type Field struct {
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Options FieldOptions `json:"options"`
}
