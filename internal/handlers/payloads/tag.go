package payloads

type Tag struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name" validate:"required"`
}
