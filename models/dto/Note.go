package dto

type Note struct {
	// ID       int64  `json:"-" sql.field:"id"`
	Name     string `json:"name,omitempty" sql.field:"name"`
	LastName string `json:"last_name,omitempty" sql.field:"last_name"`
	Text     string `json:"text,omitempty" sql.field:"text"`
}
