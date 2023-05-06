package user

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"is_admin"`
}
