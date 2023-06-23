package http

import (
	"encoding/json"
	"github.com/behnambm/go-crud/repo/sqlite"
	"github.com/behnambm/go-crud/service/auth"
	"github.com/behnambm/go-crud/service/book"
	"github.com/behnambm/go-crud/service/user"
	"github.com/labstack/echo/v4"
	"io"
	netHttp "net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var adminUserJWT string

func setupTest() (Server, func()) {
	sqliteRepo := sqlite.New("storage_test.db")
	sqlite.CreateTables(sqliteRepo)
	sqlite.SeedTables(sqliteRepo)

	userService := user.New(sqliteRepo)
	bookService := book.New(sqliteRepo)

	jwtAuthService := auth.New("a_secret_key")
	httpServer := Server{
		ListenAddr: ":8080",
		AuthSrv:    jwtAuthService,
		UserSrv:    userService,
		BookSrv:    bookService,
	}

	tearDownFunc := func() {
		os.Remove("storage_test.db")
	}

	return httpServer, tearDownFunc
}

func TestAll(t *testing.T) {
	t.Run("Successful Login", SuccessfulLogin)
	t.Run("Login With Invalid Data", InvalidLogin)
	t.Run("Successful Get Book", SuccessfulGetBook)
}

func SuccessfulLogin(t *testing.T) {
	httpServer, tearDownFunc := setupTest()
	defer tearDownFunc()

	loginData := `{"email": "test@gmail.com", "password": "123"}`
	req := httptest.NewRequest(netHttp.MethodPost, "/auth/login/", strings.NewReader(loginData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, resp)

	loginErr := httpServer.Login(c)
	if loginErr != nil {
		t.Fatal(loginErr)
	}

	var jsonResponse map[string]interface{}
	data, _ := io.ReadAll(resp.Body)
	jsonErr := json.Unmarshal(data, &jsonResponse)
	if jsonErr != nil {
		return
	}
	if resp.Code != netHttp.StatusOK {
		t.Fatal("could not login,", "got", resp.Code, "code")
	}
	if jsonResponse["token"] == "" {
		t.Fatal("could not get jwt token")
	}
	adminUserJWT, _ = jsonResponse["token"].(string)
}

func InvalidLogin(t *testing.T) {
	httpServer, tearDownFunc := setupTest()
	defer tearDownFunc()

	loginData := `{"email": "test12@gmail.com", "password": "12388"}`
	req := httptest.NewRequest(netHttp.MethodPost, "/auth/login/", strings.NewReader(loginData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, resp)

	loginErr := httpServer.Login(c)
	if loginErr != nil {
		t.Fatal(loginErr)
	}

	var jsonResponse map[string]interface{}
	data, _ := io.ReadAll(resp.Body)
	jsonErr := json.Unmarshal(data, &jsonResponse)
	if jsonErr != nil {
		return
	}
	if resp.Code != netHttp.StatusForbidden {
		t.Fatal("logged in,", "got", resp.Code, "code")
	}
	if jsonResponse["error"] == "" {
		t.Fatal("could get jwt token")
	}
	t.Log(jsonResponse["error"])
	if jsonResponse["error"] != "invalid credentials" {
		t.Fatal("invalid response message")
	}
}

func SuccessfulGetBook(t *testing.T) {
	httpServer, tearDownFunc := setupTest()
	defer tearDownFunc()

	req := httptest.NewRequest(netHttp.MethodGet, "/book/", nil)
	req.Header.Set(echo.HeaderAuthorization, adminUserJWT)
	resp := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, resp)
	c.SetPath("/book/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	getErr := httpServer.GetBook(c)
	if getErr != nil {
		t.Fatal(getErr)
	}

	var jsonResponse map[string]interface{}
	data, _ := io.ReadAll(resp.Body)
	jsonErr := json.Unmarshal(data, &jsonResponse)
	if jsonErr != nil {
		return
	}
	if resp.Code != netHttp.StatusOK {
		t.Fatal("logged in,", "got", resp.Code, "code")
	}
	t.Log(jsonResponse)
	if _, ok := jsonResponse["error"]; ok {
		t.Log(jsonResponse["error"])
		t.Fatal("got error")
	}
	if _, ok := jsonResponse["book"]; !ok {
		t.Fatal("couldn't get book")
	}
	bookData, ok := jsonResponse["book"].(map[string]interface{})
	if !ok {
		t.Fatal("response map struct is invalid")
	}
	if bookData["name"] == "" {
		t.Fatal("book name is empty")
	}
}
