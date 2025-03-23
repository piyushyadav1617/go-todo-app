package model

type Todo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Completed   bool    `json:"completed"`
}

type TodoUpdate struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}
