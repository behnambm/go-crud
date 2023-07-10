package http

import (
	_ "github.com/behnambm/go-crud/delivery/http/docs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"time"

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
	ListenAddr      string
	AuthSrv         auth.Service
	UserSrv         user.Service
	BookSrv         book.Service
	Metrics         *Metrics
	MetricsRegistry *prometheus.Registry
}

type Metrics struct {
	endpointHit *prometheus.CounterVec
	duration    *prometheus.HistogramVec
}

func NewMetrics(reg *prometheus.Registry, namespace string) *Metrics {
	m := Metrics{
		endpointHit: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "endpoint_hit",
				Help:      "Number of requests made against endpoints",
			},
			[]string{"type"},
		),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_duration_seconds",
			Help:      "Duration of the request.",
			Buckets:   []float64{0.1, 0.15, 0.2, 0.25, 0.3, 0.5, 1},
		}, []string{"status", "method"}),
	}

	reg.MustRegister(m.endpointHit, m.duration)

	return &m
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

	go func() {
		echoServer := echo.New()
		promHandler := promhttp.HandlerFor(s.MetricsRegistry, promhttp.HandlerOpts{Registry: s.MetricsRegistry})
		echoServer.GET("/metrics", echo.WrapHandler(promHandler)) // adds route to serve gathered metrics
		echoServer.Start(":9000")
	}()

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
	start := time.Now()
	statusCode := 200
	s.Metrics.endpointHit.With(prometheus.Labels{"type": "auth_login"}).Inc()
	defer func() {
		s.Metrics.duration.With(prometheus.Labels{"method": "GET", "status": strconv.Itoa(statusCode)}).Observe(time.Since(start).Seconds())
	}()

	loginRequest := param.LoginRequest{}
	if err := c.Bind(&loginRequest); err != nil {
		log.Println("LOGIN HANDLER JSON BIND ERR", err)
		statusCode = 400
		return c.JSON(http.StatusBadRequest, param.BadRequestResponse{Error: "invalid data"})
	}

	userData, err := s.UserSrv.GetUserFromEmail(loginRequest.Email)
	if err != nil {
		log.Println("LOGIN HANDLER FETCH USER ERR", err)
		statusCode = 403
		return c.JSON(http.StatusForbidden, param.BadRequestResponse{Error: "invalid credentials"})
	}

	hashedPassword, hashErr := hash.String(loginRequest.Password)

	if hashErr != nil {
		log.Println("LOGIN HANDLER HASH GENERATION ERR", hashErr)
		statusCode = 400
		return c.JSON(http.StatusBadRequest, param.BadRequestResponse{Error: "invalid data"})
	}
	if userData.Password != hashedPassword {
		statusCode = 403
		return c.JSON(http.StatusForbidden, param.BadRequestResponse{Error: "invalid credentials"})
	}

	jwt, jwtErr := s.AuthSrv.GenerateJWT(userData.ID)

	if jwtErr != nil {
		log.Println("LOGIN HANDLER Auth GENERATE ERR", jwtErr)
		statusCode = 500
		return c.JSON(http.StatusInternalServerError, param.BadRequestResponse{Error: "server error"})
	}
	return c.JSON(http.StatusOK, param.LoginOKResponse{Token: jwt})
}

func (s Server) GetBook(c echo.Context) error {
	start := time.Now()
	statusCode := 200
	s.Metrics.endpointHit.With(prometheus.Labels{"type": "book_get_single"}).Inc()
	defer func() {
		s.Metrics.duration.With(prometheus.Labels{"method": "GET", "status": strconv.Itoa(statusCode)}).Observe(time.Since(start).Seconds())
	}()

	id, convErr := strconv.Atoi(c.Param("id"))
	if convErr != nil {
		statusCode = 400
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	bookFromDB, err := s.BookSrv.GetBook(id)
	if err != nil {
		log.Println("GET BOOK HANDLER - GET BOOK SERVICE ERR", err)
		statusCode = 404
		return c.JSON(http.StatusNotFound, echo.Map{"error": "couldn't get book"})
	}
	if !httpUtils.IsAuthenticated(c) || !httpUtils.IsAdmin(c) {
		if !bookFromDB.IsPublished {
			statusCode = 404
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
	start := time.Now()
	statusCode := 200
	s.Metrics.endpointHit.With(prometheus.Labels{"type": "book_get_list"}).Inc()
	defer func() {
		s.Metrics.duration.With(prometheus.Labels{"method": "GET", "status": strconv.Itoa(statusCode)}).Observe(time.Since(start).Seconds())
	}()

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	if !httpUtils.IsAuthenticated(c) || !httpUtils.IsAdmin(c) {
		books, err := s.BookSrv.PublishedBookList()
		if err != nil {
			log.Println("PUBLISHED BOOK LIST HANDLER ERR", err)
			statusCode = 404
			return c.JSON(http.StatusNotFound, echo.Map{"error": "couldn't get list of books"})
		}
		return c.JSON(http.StatusOK, echo.Map{"books": books})
	}

	books, err := s.BookSrv.BookList()
	if err != nil {
		log.Println("BOOK LIST HANDLER ERR", err)
		statusCode = 404
		return c.JSON(http.StatusNotFound, echo.Map{"error": "couldn't get list of books"})
	}
	return c.JSON(http.StatusOK, echo.Map{"books": books})
}

func (s Server) CreateBook(c echo.Context) error {
	start := time.Now()
	statusCode := 200
	s.Metrics.endpointHit.With(prometheus.Labels{"type": "book_create"}).Inc()
	defer func() {
		s.Metrics.duration.With(prometheus.Labels{"method": "POST", "status": strconv.Itoa(statusCode)}).Observe(time.Since(start).Seconds())
	}()

	bookCreateRequest := param.BookCreateRequest{}
	if bindErr := c.Bind(&bookCreateRequest); bindErr != nil {
		log.Println("CREATE BOOK HANDLER BIND ERR", bindErr)
		statusCode = 400
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data"})
	}

	if !httpUtils.IsAdmin(c) {
		bookCreateRequest.IsPublished = false
	}
	createdBook, createErr := s.BookSrv.CreateBook(bookCreateRequest)
	if createErr != nil {
		log.Println("CREATE BOOK HANDLER CREATE SERVICE ERR", createErr)
		statusCode = 500
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "couldn't create book"})
	}

	return c.JSON(http.StatusOK, echo.Map{"book": createdBook})
}

func (s Server) UpdateBook(c echo.Context) error {
	start := time.Now()
	statusCode := 200
	s.Metrics.endpointHit.With(prometheus.Labels{"type": "book_update"}).Inc()
	defer func() {
		s.Metrics.duration.With(prometheus.Labels{"method": "PUT", "status": strconv.Itoa(statusCode)}).Observe(time.Since(start).Seconds())
	}()

	id := c.Param("id")
	bookId, convErr := strconv.Atoi(id)
	if convErr != nil {
		statusCode = 400
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	bookUpdateRequest := param.BookUpdateRequest{}
	if bindErr := c.Bind(&bookUpdateRequest); bindErr != nil {
		statusCode = 400
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid data"})
	}

	updatedBook, updateErr := s.BookSrv.UpdateBook(bookId, bookUpdateRequest)
	if updateErr != nil {
		log.Println("UPDATE BOOK HANDLER UPDATE ERR", updateErr)
		statusCode = 500
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "couldn't update the book"})
	}

	return c.JSON(http.StatusOK, echo.Map{"book": updatedBook})
}

func (s Server) DeleteBook(c echo.Context) error {
	start := time.Now()
	statusCode := 204
	s.Metrics.endpointHit.With(prometheus.Labels{"type": "book_delete"}).Inc()
	defer func() {
		s.Metrics.duration.With(prometheus.Labels{"method": "DELETE", "status": strconv.Itoa(statusCode)}).Observe(time.Since(start).Seconds())
	}()

	id := c.Param("id")
	bookId, convErr := strconv.Atoi(id)
	if convErr != nil {
		statusCode = 400
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	if deleteErr := s.BookSrv.DeleteBook(bookId); deleteErr != nil {
		log.Println("DELETE BOOK HANDLER DELETE ERR", deleteErr)
		statusCode = 500
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "couldn't delete the book"})
	}

	return c.NoContent(http.StatusNoContent)
}
