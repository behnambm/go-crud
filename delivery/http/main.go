package http

import (
	"github.com/behnambm/assignment/delivery/http/middleware"
	"github.com/behnambm/assignment/param"
	"github.com/behnambm/assignment/service/auth"
	"github.com/behnambm/assignment/service/book"
	"github.com/behnambm/assignment/service/user"
	"github.com/behnambm/assignment/utils/hash"
	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	ListenAddr string
	AuthSrv    auth.Service
	UserSrv    user.Service
	BookSrv    book.Service
}

func (s Server) Run() {
	e := echo.New()
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())

	authRoute := e.Group("/auth")
	authRoute.POST("/login", s.Login)

	bookRoute := e.Group("/book")
	bookRoute.GET("/", s.BookList)
	bookRoute.POST("/", s.CreateBook, middleware.Auth(s.UserSrv, s.AuthSrv))
	bookRoute.PUT("/:id", s.UpdateBook, middleware.Auth(s.UserSrv, s.AuthSrv))
	bookRoute.DELETE("/:id", s.DeleteBook, middleware.Auth(s.UserSrv, s.AuthSrv))

	e.Logger.Fatal(e.Start(s.ListenAddr))
}

func (s Server) Login(c echo.Context) error {
	loginRequest := param.LoginRequest{}
	if err := c.Bind(&loginRequest); err != nil {
		log.Println("LOGIN HANDLER JSON BIND ERR", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data"})
	}

	userData, err := s.UserSrv.GetUserFromEmail(loginRequest.Email)
	if err != nil {
		log.Println("LOGIN HANDLER FETCH USER ERR", err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "invalid credentials"})
	}

	hashedPassword, hashErr := hash.String(loginRequest.Password)

	if hashErr != nil {
		log.Println("LOGIN HANDLER HASH GENERATION ERR", hashErr)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data"})
	}
	if userData.Password != hashedPassword {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "invalid credentials"})
	}

	jwt, jwtErr := s.AuthSrv.GenerateJWT(userData.ID)

	if jwtErr != nil {
		log.Println("LOGIN HANDLER Auth GENERATE ERR", jwtErr)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": jwt})
}

func (s Server) BookList(c echo.Context) error {
	books, err := s.BookSrv.BookList()
	if err != nil {
		log.Println("BOOK LIST HANDLER ERR", err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "couldn't get list of books"})
	}
	return c.JSON(http.StatusOK, echo.Map{"books": books})
}

func (s Server) CreateBook(c echo.Context) error {
	bookCreateRequest := param.BookCreateRequest{}
	if bindErr := c.Bind(&bookCreateRequest); bindErr != nil {
		log.Println("CREATE BOOK HANDLER BIND ERR", bindErr)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data"})
	}

	createdBook, createErr := s.BookSrv.CreateBook(bookCreateRequest)
	if createErr != nil {
		log.Println("CREATE BOOK HANDLER CREATE SERVICE ERR", createErr)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "couldn't create book"})
	}

	return c.JSON(http.StatusOK, echo.Map{"book": createdBook})
}

func (s Server) UpdateBook(c echo.Context) error {
	id := c.Param("id")
	bookId, convErr := strconv.Atoi(id)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	bookUpdateRequest := param.BookUpdateRequest{}
	if bindErr := c.Bind(&bookUpdateRequest); bindErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data"})
	}

	updatedBook, updateErr := s.BookSrv.UpdateBook(bookId, bookUpdateRequest)
	if updateErr != nil {
		log.Println("UPDATE BOOK HANDLER UPDATE ERR", updateErr)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "couldn't update the book"})
	}

	return c.JSON(http.StatusOK, echo.Map{"book": updatedBook})
}

func (s Server) DeleteBook(c echo.Context) error {
	id := c.Param("id")
	bookId, convErr := strconv.Atoi(id)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	if deleteErr := s.BookSrv.DeleteBook(bookId); deleteErr != nil {
		log.Println("DELETE BOOK HANDLER DELETE ERR", deleteErr)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "couldn't delete the book"})
	}

	return c.NoContent(http.StatusNoContent)
}
