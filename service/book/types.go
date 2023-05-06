package book

type Book struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float32 `json:"price"`
	IsPublished bool    `json:"is_published"`
}
