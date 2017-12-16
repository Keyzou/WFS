package controllers_test

import(
	"net/http/httptest"
	"net/http"
	"testing"
	"strings"
	"controllers"
	"encoding/json"
	"models"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"strconv"
)

func TestUser_CreateUser(t *testing.T){
	user := `{"username": "test2", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		userM := &models.User{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}

	conn.Query(`DELETE FROM users WHERE username='test2'`)	
}

func TestUser_CantCreateEmptyUsername(t *testing.T){
	user := `{"username": "", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateUserInvalidRequest(t *testing.T){
	user := `{"username": "test2", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, "snwd")
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}


func TestUser_CantCreateEmptyEmail(t *testing.T){
	user := `{"username": "test2", "email": "", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateEmptyPassword(t *testing.T){
	user := `{"username": "test2", "email": "test2@test.fr", "password": ""}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateInvalidEmail(t *testing.T){
	user := `{"username": "test2", "email": "test2test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateInvalidUsername(t *testing.T){
	user := `{"username": ".,>>,wtr", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateShortUsername(t *testing.T){
	user := `{"username": "as", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateShortPassword(t *testing.T){
	user := `{"username": "test2", "email": "test2@test.fr", "password": "test"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateDuplicateUsername(t *testing.T){
	user := `{"username": "test", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Signup(context))
}

func TestUser_CantCreateDuplicateEmail(t *testing.T){
	user := `{"username": "tewst2", "email": "test@test.fr", "password": "tewst123"}`
	e := echo.New()
	reader := strings.NewReader(user)
	req := httptest.NewRequest(echo.POST, "/signup", reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	resp := c.Signup(context) 
	assert.Error(t, resp)
}

func TestUser_LoginUser(t *testing.T){
	user := `{"email": "test@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.Login(context)){
		if assert.Equal(t, http.StatusOK, rec.Code){
			var dat map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &dat)
			assert.NotNil(t, dat["jwt"])
		}
	}
}

func TestUser_CantLoginInvalidRequest(t *testing.T){
	user := `{"email": 52, "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Login(context))
}

func TestUser_CantLoginInvalidEmail(t *testing.T){
	user := `{"email": "testtest.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Login(context))
}


func TestUser_CantLoginEmptyEmail(t *testing.T){
	user := `{"email": "", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Login(context))
}

func TestUser_CantLoginEmptyPassword(t *testing.T){
	user := `{"email": "test@test.fr", "password": ""}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Login(context))
}

func TestUser_CantLoginUnknownEmail(t *testing.T){
	user := `{"email": "tes444t@test.fr", "password": "test1324"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Login(context))
}

func TestUser_CantLoginWrongPassword(t *testing.T){
	user := `{"email": "test@test.fr", "password": "test1324"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/login", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.Login(context))
}



func TestUser_UpdateUsername(t *testing.T){
	user := `{"username": "test3"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/update", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.UpdateUser(context)){
		if assert.Equal(t, http.StatusOK, rec.Code){
			newName := &models.User{}
			json.Unmarshal(rec.Body.Bytes(), &newName)
			assert.Equal(t, "test3", newName.Username)
		}
	}
	conn.Query(`update users set username='test' where username='test3'`)	
}

func TestUser_CantUpdateInvalidRequest(t *testing.T){
	user := `{"username": 456}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/update", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	assert.Error(t, c.UpdateUser(context))
}

func TestUser_CantUpdateNoChanges(t *testing.T){
	user := `{}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/update", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UpdateUser(context))
}

func TestUser_CantUpdateInvalidChanges(t *testing.T){
	user := `{"username": "test.,.3"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/update", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UpdateUser(context))
}

func TestUser_CantUpdateShortUsername(t *testing.T){
	user := `{"username": "te"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/update", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UpdateUser(context))
}

func TestUser_CantUpdateShortPassword(t *testing.T){
	user := `{"password": "tsae"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/user/update", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UpdateUser(context))
}

func TestUser_CantUpdateDuplicateUsername(t *testing.T){
	user := `{"username": "duplicate", "email": "duplicate@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		userM := &models.User{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}


	user = `{"username": "duplicate"}`
	req = httptest.NewRequest(echo.POST, "/user/update", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	c = &controllers.Controller{DB: conn}

	assert.Error(t, c.UpdateUser(context))

	conn.Query(`DELETE FROM users WHERE username='duplicate'`)	
}

func TestUser_QueryUsers(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/queryUsers/:query")
	context.SetParamNames("query")
	context.SetParamValues("te")
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.QueryUsers(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
		var users []*models.User
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))
		assert.Equal(t, users[0].ID, test.ID)
	}
}

func TestUser_CantQueryUsersEmptyQuery(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/queryUsers/:query")
	context.SetParamNames("query")
	context.SetParamValues("")
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.QueryUsers(context))
}

func TestUser_FollowUser(t *testing.T){
	user := `{"username": "test2", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	userM := &models.User{}

	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}

	req = httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))
	if assert.NoError(t, c.Follow(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	conn.Query(`DELETE FROM users WHERE username='test2'`)	
}

func TestUser_CantFollowInvalidID(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	c := &controllers.Controller{DB: conn}
	context := e.NewContext(req, rec)
	context.SetPath("/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues("string")
	assert.Error(t, c.Follow(context))

}

func TestUser_CantFollowSelf(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	context.SetPath("/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(test.ID))
	assert.Error(t, c.Follow(context))
}

func TestUser_CantFollowUnknownUser(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	context.SetPath("/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues("-1")
	assert.Error(t, c.Follow(context))
}

func TestUser_CantUnfollowInvalidID(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	c := &controllers.Controller{DB: conn}
	context := e.NewContext(req, rec)
	context.SetPath("/unfollow/:id")
	context.SetParamNames("id")
	context.SetParamValues("string")
	assert.Error(t, c.Unfollow(context))

}

func TestUser_CantUnfollowSelf(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	context.SetPath("/unfollow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(test.ID))
	assert.Error(t, c.Unfollow(context))
}

func TestUser_CantUnfollowUnknownUser(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	context.SetPath("/unfollow/:id")
	context.SetParamNames("id")
	context.SetParamValues("-1")
	assert.Error(t, c.Unfollow(context))
}

func TestUser_CantFollowUserTwice(t *testing.T){
	user := `{"username": "tofollow", "email": "tofollow@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	userM := &models.User{}
	
	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}

	req = httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))

	if assert.NoError(t, c.Follow(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))

	assert.Error(t, c.Follow(context))
	

	conn.Query(`DELETE FROM users WHERE username='tofollow'`)	
}

func TestUser_CantUnfollowUserTwice(t *testing.T){
	user := `{"username": "tounfollow", "email": "tounfollow@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	userM := &models.User{}
	
	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}

	req2 := httptest.NewRequest(echo.POST, "/", nil)
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req2.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec2 := httptest.NewRecorder()
	context2 := e.NewContext(req2, rec2)
	context2.SetPath("/follow/:id")
	context2.SetParamNames("id")
	context2.SetParamValues(strconv.Itoa(userM.ID))

	if assert.NoError(t, c.Follow(context2)){
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	req = httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req2, rec2)
	context.SetPath("/unfollow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))

	if assert.NoError(t, c.Unfollow(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req2, rec2)
	context.SetPath("/unfollow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))

	assert.Error(t, c.Unfollow(context))
	

	conn.Query(`DELETE FROM users WHERE username='tounfollow'`)	
}

func TestUser_UnfollowUser(t *testing.T){
	user := `{"username": "test2", "email": "test2@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	assert.NoError(t, c.Signup(context))
	userM := &models.User{}
	json.Unmarshal(rec.Body.Bytes(), userM)

	req2 := httptest.NewRequest(echo.POST, "/", nil)
	req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req2.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec2 := httptest.NewRecorder()
	context2 := e.NewContext(req2, rec2)
	context2.SetPath("/follow/:id")
	context2.SetParamNames("id")
	context2.SetParamValues(strconv.Itoa(userM.ID))
	assert.NoError(t, c.Follow(context2))

	req = httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req2, rec2)
	context.SetPath("/unfollow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))

	if assert.NoError(t, c.Unfollow(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	conn.Query(`DELETE FROM users WHERE username='test2'`)	
}

func TestUser_GetUser(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/user/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(test.ID))
	c := &controllers.Controller{DB: conn}
	user := &models.User{}
	if assert.NoError(t, c.GetUser(context)){
		json.Unmarshal(rec.Body.Bytes(), user)
		assert.Equal(t, test.ID, user.ID)
	}
}

func TestUser_CantGetInvalidUser(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/user/:id")
	context.SetParamNames("id")
	context.SetParamValues("string")
	c := &controllers.Controller{DB: conn}
	assert.Error(t, c.GetUser(context))
}

func TestUser_CantGetUnknownUser(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/user/:id")
	context.SetParamNames("id")
	context.SetParamValues("-1")
	c := &controllers.Controller{DB: conn}
	assert.Error(t, c.GetUser(context))
}

func TestUser_GetFollowerCount(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/user/:id/followers")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(test.ID))
	c := &controllers.Controller{DB: conn}
	if assert.NoError(t, c.GetFollowers(context)){
		assert.Equal(t, "0", rec.Body.String())
	}
}

func TestUser_CantGetFollowerCountInvalidQuery(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/user/:id/followers")
	context.SetParamNames("id")
	context.SetParamValues("string")
	c := &controllers.Controller{DB: conn}
	assert.Error(t, c.GetFollowers(context))
}

func TestUser_InvalidGetFollowing(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath("/user/:id/following")
	context.SetParamNames("id")
	context.SetParamValues("string")
	c := &controllers.Controller{DB: conn}
	assert.Error(t, c.GetFollowing(context))
}

func TestUser_GetFollowing(t *testing.T){
	user := `{"username": "following", "email": "following@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	userM := &models.User{}

	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}

	
	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/user/:id/following")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))
	if assert.NoError(t, c.GetFollowing(context)){
		assert.Equal(t, "false", rec.Body.String())
	}

}

func TestUser_GetFollowingUsers(t *testing.T){
	user := `{"username": "followinguser", "email": "followinguser@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	userM := &models.User{}

	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}

	
	req = httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/user/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))
	assert.NoError(t, c.Follow(context))

	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/user/following/")
	if assert.NoError(t, c.GetFollowingUsers(context)){
		var users []*models.User
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))
		assert.Equal(t, 1, len(users))
		assert.Equal(t, userM.ID, users[0].ID)
	}
}


func TestUser_GetFollowersUsers(t *testing.T){
	user := `{"username": "followersuser", "email": "followersuser@test.fr", "password": "test123"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/signup", strings.NewReader(user))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}
	userM := &models.User{}
	
	if assert.NoError(t, c.Signup(context)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), userM))
	}
	
	log := `{"email": "followersuser@test.fr", "password": "test123"}`
	req = httptest.NewRequest(echo.POST, "/login", strings.NewReader(log))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	var token string
	if assert.NoError(t, c.Login(context)){
		if assert.Equal(t, http.StatusOK, rec.Code){
			var dat map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &dat)
			assert.NotNil(t, dat["jwt"])
			token = dat["jwt"].(string)
		}
	}
	
	req = httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/user/follow/:id")
	context.SetParamNames("id")
	context.SetParamValues(strconv.Itoa(userM.ID))
	assert.NoError(t, c.Follow(context))

	req = httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	context.SetPath("/user/followers")
	if assert.NoError(t, c.GetFollowersUsers(context)){
		var users []*models.User
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))
		assert.Equal(t, 1, len(users))
		assert.Equal(t, test.ID, users[0].ID)
	}
}