package param

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BookCreateRequest struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type BookUpdateRequest struct {
	BookCreateRequest
}
