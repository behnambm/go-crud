package param

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BookCreateRequest struct {
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	IsPublished bool    `json:"is_published,omitempty"`
}

type BookUpdateRequest struct {
	BookCreateRequest
	IsPublished bool `json:"is_published"`
}

type MinimalBookResponse struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type FullBookResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	IsPublished bool    `json:"is_published"`
}
