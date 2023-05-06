package main

import (
	"flag"
	"github.com/behnambm/assignment/delivery/http"
	"github.com/behnambm/assignment/repo/sqlite"
	"github.com/behnambm/assignment/service/auth"
	"github.com/behnambm/assignment/service/book"
	"github.com/behnambm/assignment/service/user"
)

func main() {
	initDBFlag := flag.Bool("initdb", false, "Create and Seed the database")
	flag.Parse()

	sqliteRepo := sqlite.New("storage.db")

	if *initDBFlag {
		if err := sqlite.CreateTables(sqliteRepo); err != nil {
			panic(err)
		}
		sqlite.SeedTables(sqliteRepo)

		return
	}

	userService := user.New(sqliteRepo)
	bookService := book.New(sqliteRepo)

	jwtAuthService := auth.New("a_secret_key")
	httpServer := http.Server{
		ListenAddr: ":8080",
		AuthSrv:    jwtAuthService,
		UserSrv:    userService,
		BookSrv:    bookService,
	}
	httpServer.Run()
}
