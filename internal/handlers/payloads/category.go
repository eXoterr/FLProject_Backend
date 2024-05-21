package payloads

type Category struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name" validate:"required"`
}
