package main

import (
	"flag"
	"github.com/behnambm/go-crud/delivery/http"
	"github.com/behnambm/go-crud/repo/sqlite"
	"github.com/behnambm/go-crud/service/auth"
	"github.com/behnambm/go-crud/service/book"
	"github.com/behnambm/go-crud/service/user"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/swaggo/echo-swagger"
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

	reg := prometheus.NewRegistry()
	metrics := http.NewMetrics(reg, "book_app")

	httpServer := http.Server{
		ListenAddr:      ":8080",
		AuthSrv:         jwtAuthService,
		UserSrv:         userService,
		BookSrv:         bookService,
		Metrics:         metrics,
		MetricsRegistry: reg,
	}
	httpServer.Run()
}
