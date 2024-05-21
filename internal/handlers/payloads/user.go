package payloads

type UserPayload interface {
}

type User struct {
	Email    string `json:"email,omitempty" validate:"email,required"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=254"`
}

type JWTUser struct {
	ID   int    `json:"id,omitempty"`
	Role string `json:"role,omitempty"`
}

type JWTTokenPair struct {
	RefreshToken string `json:"refresh,omitempty"`
	AccessToken  string `json:"token,omitempty"`
}

type JWTRefresh struct {
	RefreshToken string `json:"token,omitempty" validate:"required"`
}
