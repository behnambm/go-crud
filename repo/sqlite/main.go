package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/behnambm/assignment/utils/hash"
	"log"
)
import _ "github.com/mattn/go-sqlite3"

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

// CreateTables will be used to create the tables that we need.
// a better solution would be using migrator
func CreateTables(repo *Repo) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) UNIQUE,
	    password VARCHAR(255) NOT NULL
	)
	`
	_, err := repo.db.Exec(userTable)
	if err != nil {
		return err
	}
	bookTable := `
	CREATE TABLE IF NOT EXISTS book (
	    id INTEGER PRIMARY KEY AUTOINCREMENT, 
	    name VARCHAR(255) NOT NULL UNIQUE,
	   	price REAL,
	    is_published BOOLEAN NOT NULL
	)
	`
	_, err = repo.db.Exec(bookTable)
	if err != nil {
		return fmt.Errorf("book table err %w", err)
	}
	log.Println("DATABASE CREATED")
	return nil
}

func SeedTables(repo *Repo) {
	userEmail := "test@gmail.com"
	passwordHash, hashErr := hash.String("123")
	if hashErr != nil {
		panic(hashErr)
	}
	_, err := repo.db.Exec("INSERT INTO user (email, password) VALUES (?, ?)", userEmail, passwordHash)
	if err != nil {
		panic(err)
	}
	_, err = repo.db.Exec("INSERT INTO book (name, price, is_published) VALUES ('Go development book', 50.99, 1)")
	if err != nil {
		panic(err)
	}
	_, err = repo.db.Exec("INSERT INTO book (name, price, is_published) VALUES ('Go development book Vol.2', 50.99, 0)")
	if err != nil {
		panic(err)
	}

	log.Println("DATABASE SEED SUCCESSFUL")
}

func (r Repo) GetUserFromEmail(email string) (User, error) {
	row := r.db.QueryRow(`SELECT * FROM user WHERE email = ?`, email)
	if row.Err() != nil {
		return User{}, row.Err()
	}
	user := User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
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
	if err := row.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		return User{}, err
	}
	return user, nil
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
