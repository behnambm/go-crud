package http

import (
	_ "github.com/behnambm/go-crud/delivery/http/docs"

	"github.com/behnambm/go-crud/delivery/http/middleware"
	"github.com/behnambm/go-crud/param"
	"github.com/behnambm/go-crud/service/auth"
	"github.com/behnambm/go-crud/service/book"
	"github.com/behnambm/go-crud/service/user"
	"github.com/behnambm/go-crud/utils/hash"
	httpUtils "github.com/behnambm/go-crud/utils/http"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
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

// @title						Go CRUD API SPEC
// @version					1.0
// @description				This document will provide information about using this API
// @contact.name				Behnam Mohammadzadeh
// @contact.url				https://blog.behnambm.ir/
// @contact.email				behnam.mohamadzadeh21@gmail.com
// @host						http://localhost:8080
// @BasePath					/swagger
//
// @securityDefinitions.basic	BasicAuth
func (s Server) Run() {
	e := echo.New()
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())

	authRoute := e.Group("/auth")
	authRoute.POST("/login", s.Login)

	bookRoute := e.Group("/book", middleware.Auth(s.UserSrv, s.AuthSrv))
	bookRoute.GET("/", s.BookList)
	bookRoute.GET("/:id", s.GetBook)
	bookRoute.POST("/", s.CreateBook, middleware.LoginRequired())
	bookRoute.PUT("/:id", s.UpdateBook, middleware.LoginRequired(), middleware.AdminRequired())
	bookRoute.DELETE("/:id", s.DeleteBook, middleware.LoginRequired(), middleware.AdminRequired())

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(s.ListenAddr))
}

// Login godoc
//
//	@Summary		Login the user
//	@Description	Using this route you can authenticate and get the JWT token if provided credentials are valid
//	@Tags			Auth
//	@Param			request	body	param.LoginRequest	true	"query params"
//	@Accept			application/json
//	@Accept			text/xml
//	@Produce		json
//	@success		200	{object}	param.LoginOKResponse		"Token field contains JWT token"
//	@failure		400	{object}	param.BadRequestResponse{}	"invalid data"
//	@failure		403	{object}	param.BadRequestResponse{}	"invalid credentials"
//	@failure		500	{object}	param.BadRequestResponse{}	"server error"
//	@Router			/auth/login [post]
//	@Security		BasicAuth
func (s Server) Login(c echo.Context) error {
	loginRequest := param.LoginRequest{}
	if err := c.Bind(&loginRequest); err != nil {
		log.Println("LOGIN HANDLER JSON BIND ERR", err)
		return c.JSON(http.StatusBadRequest, param.BadRequestResponse{Error: "invalid data"})
	}

	userData, err := s.UserSrv.GetUserFromEmail(loginRequest.Email)
	if err != nil {
		log.Println("LOGIN HANDLER FETCH USER ERR", err)
		return c.JSON(http.StatusForbidden, param.BadRequestResponse{Error: "invalid credentials"})
	}

	hashedPassword, hashErr := hash.String(loginRequest.Password)

	if hashErr != nil {
		log.Println("LOGIN HANDLER HASH GENERATION ERR", hashErr)
		return c.JSON(http.StatusBadRequest, param.BadRequestResponse{Error: "invalid data"})
	}
	if userData.Password != hashedPassword {
		return c.JSON(http.StatusForbidden, param.BadRequestResponse{Error: "invalid credentials"})
	}

	jwt, jwtErr := s.AuthSrv.GenerateJWT(userData.ID)

	if jwtErr != nil {
		log.Println("LOGIN HANDLER Auth GENERATE ERR", jwtErr)
		return c.JSON(http.StatusInternalServerError, param.BadRequestResponse{Error: "server error"})
	}
	return c.JSON(http.StatusOK, param.LoginOKResponse{Token: jwt})
}

func (s Server) GetBook(c echo.Context) error {
	id, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	bookFromDB, err := s.BookSrv.GetBook(id)
	if err != nil {
		log.Println("GET BOOK HANDLER - GET BOOK SERVICE ERR", err)
		return c.JSON(http.StatusNotFound, echo.Map{"error": "couldn't get book"})
	}
	if !httpUtils.IsAuthenticated(c) || !httpUtils.IsAdmin(c) {
		if !bookFromDB.IsPublished {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "book not found"})
		}
		responseBook := param.MinimalBookResponse{
			ID:    bookFromDB.ID,
			Name:  bookFromDB.Name,
			Price: bookFromDB.Price,
		}
		return c.JSON(http.StatusOK, echo.Map{"book": responseBook})
	}

	responseBook := param.FullBookResponse{
		ID:          bookFromDB.ID,
		Name:        bookFromDB.Name,
		Price:       bookFromDB.Price,
		IsPublished: bookFromDB.IsPublished,
	}
	return c.JSON(http.StatusOK, echo.Map{"book": responseBook})
}

func (s Server) BookList(c echo.Context) error {
	if !httpUtils.IsAuthenticated(c) || !httpUtils.IsAdmin(c) {
		books, err := s.BookSrv.PublishedBookList()
		if err != nil {
			log.Println("PUBLISHED BOOK LIST HANDLER ERR", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "couldn't get list of books"})
		}
		return c.JSON(http.StatusOK, echo.Map{"books": books})
	}

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

	if !httpUtils.IsAdmin(c) {
		bookCreateRequest.IsPublished = false
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
