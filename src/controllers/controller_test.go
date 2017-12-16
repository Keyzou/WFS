package controllers_test

import(
	"net/http/httptest"
	"testing"
	"strings"
	"controllers"
	"encoding/json"
	"models"
	"github.com/labstack/echo"
	"github.com/jackc/pgx"
	"os"
	"strconv"
)

var(
	dbConfig = pgx.ConnPoolConfig{MaxConnections: 15, ConnConfig: pgx.ConnConfig{Host: "localhost", User: "postgres", Password: "postgres", Database: "spt_test"}}
	conn, err = pgx.NewConnPool(dbConfig)
	test = &models.User{}
	testPost = &models.Post{}
	testComment = &models.Comment{}
)

func TestMain(m *testing.M) {
	conn.Query(`DELETE FROM users`)
	conn.Query(`DELETE FROM posts`)	
	conn.Query(`DELETE FROM comments`)
	conn.Query(`DELETE FROM followers`)

	user := `{"username": "Test", "email": "test@test.fr", "password": "test123"}`
	post := `{"content": "test"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	c.Signup(context)
	json.Unmarshal(rec.Body.Bytes(), test)

	req = httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	c.Login(context)
	var dat map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &dat)
	test.Token = dat["jwt"].(string)

	req = httptest.NewRequest(echo.POST, "/create-post", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	c.CreatePost(context)
	json.Unmarshal(rec.Body.Bytes(), testPost)


	comm := `{"content": "test", "postId": `+strconv.Itoa(testPost.ID)+`}`
	req = httptest.NewRequest(echo.POST, "/create-comment", strings.NewReader(comm))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	c.CreateComment(context)
	json.Unmarshal(rec.Body.Bytes(), testComment)

	code := m.Run()

	conn.Query(`DELETE FROM users`)	
	conn.Query(`DELETE FROM posts`)	
	conn.Query(`DELETE FROM comments`)
	conn.Query(`DELETE FROM followers`)
	
	
	os.Exit(code)
}