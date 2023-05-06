package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Repo struct {
	db *sql.DB
}

func New(dsn string) *Repo {
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		panic(err)
	}

	if pingErr := conn.Ping(); pingErr != nil {
		panic(pingErr)
	}
	return &Repo{
		db: conn,
	}
}

func (r Repo) GetUserFromEmail(email string) (User, error) {
	row := r.db.QueryRow(`SELECT * FROM user WHERE email = ?`, email)
	if row.Err() != nil {
		return User{}, row.Err()
	}
	user := User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r Repo) GetUserFromID(id int) (User, error) {
	row := r.db.QueryRow(`SELECT * FROM user WHERE id = ?`, id)
	if row.Err() != nil {
		return User{}, row.Err()
	}
	user := User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r Repo) GetBook(id int) (Book, error) {
	row := r.db.QueryRow("SELECT id, name, price, is_published FROM book WHERE id = ?", id)
	if row.Err() != nil {
		return Book{}, row.Err()
	}
	var book Book
	scanErr := row.Scan(&book.ID, &book.Name, &book.Price, &book.IsPublished)
	if scanErr != nil {
		return Book{}, scanErr
	}
	return book, nil
}

func (r Repo) BookList() ([]Book, error) {
	rows, err := r.db.Query("SELECT id, name, price, is_published FROM book")
	if err != nil {
		return nil, err
	}
	var bookList []Book
	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.ID, &book.Name, &book.Price, &book.IsPublished)
		if err != nil {
			log.Println("BOOK LIST ERR", err)
		}
		bookList = append(bookList, book)
	}
	return bookList, nil
}

func (r Repo) PublishedBookList() ([]Book, error) {
	rows, err := r.db.Query("SELECT id, name, price, is_published FROM book WHERE is_published = true")
	if err != nil {
		return nil, err
	}
	var bookList []Book
	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.ID, &book.Name, &book.Price, &book.IsPublished)
		if err != nil {
			log.Println("PUBLISHED BOOK LIST ERR", err)
		}
		bookList = append(bookList, book)
	}
	return bookList, nil
}

func (r Repo) CreateBook(name string, price float32, isPublished bool) (Book, error) {
	res, execErr := r.db.Exec("INSERT INTO book (name, price, is_published) VALUES (?, ?, ?)", name, price, isPublished)
	if execErr != nil {
		return Book{}, execErr
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Book{}, err
	}
	return Book{
		ID:          int(id),
		Name:        name,
		Price:       price,
		IsPublished: isPublished,
	}, nil
}

func (r Repo) UpdateBook(id int, name string, price float32, isPublished bool) (Book, error) {
	res, execErr := r.db.Exec(
		"UPDATE book SET name = ?, price = ?, is_published = ? WHERE id = ?",
		name, price, isPublished, id,
	)
	if execErr != nil {
		return Book{}, execErr
	}
	affected, affectedErr := res.RowsAffected()
	if affectedErr != nil {
		return Book{}, affectedErr
	}
	if affected == 0 {
		return Book{}, fmt.Errorf("couldn't update the book")
	}
	return Book{
		ID:          int(id),
		Name:        name,
		Price:       price,
		IsPublished: isPublished,
	}, nil
}

func (r Repo) DeleteBook(id int) error {
	res, execErr := r.db.Exec("DELETE FROM book WHERE id = ?", id)
	if execErr != nil {
		return execErr
	}
	affected, affectedErr := res.RowsAffected()
	if affectedErr != nil {
		return affectedErr
	}
	if affected == 0 {
		return fmt.Errorf("no rows affected")
	} else if affected > 1 {
		log.Println("MORE THAT ONE ROWS AFFECTED - book id ", id)
	}
	return nil
}
