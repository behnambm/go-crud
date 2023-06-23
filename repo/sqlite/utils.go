package sqlite

import (
	"fmt"
	"github.com/behnambm/go-crud/utils/hash"
	"log"
)

// CreateTables will be used to create the tables that we need.
// a better solution would be using migrator
func CreateTables(repo *Repo) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS user (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) UNIQUE,
	    password VARCHAR(255) NOT NULL,
	    is_admin BOOLEAN NOT NULL
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
	_, err := repo.db.Exec("INSERT INTO user (email, password, is_admin) VALUES (?, ?, 1)", userEmail, passwordHash)
	if err != nil {
		panic(err)
	}
	userEmail = "test2@gmail.com"
	passwordHash, hashErr = hash.String("123")
	if hashErr != nil {
		panic(hashErr)
	}
	_, err = repo.db.Exec("INSERT INTO user (email, password, is_admin) VALUES (?, ?, 0)", userEmail, passwordHash)
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
