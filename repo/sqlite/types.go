package sqlite

type User struct {
	ID       int
	Email    string
	Password string
	IsAdmin  bool
}

type Book struct {
	ID          int
	Name        string
	Price       float32
	IsPublished bool
}
