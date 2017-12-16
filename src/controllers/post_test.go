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

func TestPost_CreatePost(t *testing.T){
	post := `{"content": "test"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/create-post", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.CreatePost(context)){
		if assert.Equal(t, http.StatusCreated, rec.Code){
			newPost := &models.Post{}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &newPost))
		}
	}
}

func TestPost_CreateInvalidPost(t *testing.T){
	post := `{"content": ""}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/create-post", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.CreatePost(context))
}

func TestPost_CreateInvalidHeaderPost(t *testing.T){
	post := `{"content": ""}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/create-post", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, "sokd")
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.CreatePost(context))
}

func TestPost_CantCreatePostUnknownUser(t *testing.T){
	post := `{"content": ""}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/create-post", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmluZyI6MTUxMzY1MzY0NCwiaWQiOi0xLCJ1c2VybmFtZSI6InRlc3QifQ.OdQUAmlPFy0ONzV2KguDdJOCwM9S4liG4mgzdTYUPIo")
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.CreatePost(context))
}


func TestPost_LikeAndUnlikePost(t *testing.T){
	post := `{"id": `+strconv.Itoa(testPost.ID)+`}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.LikePost(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	req = httptest.NewRequest(echo.POST, "/unlikePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	if assert.NoError(t, c.UnlikePost(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestPost_CantLikeInvalidPost(t *testing.T){
	post := `{"id": "-1"}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.LikePost(context))
}


func TestPost_CantUnlikeInvalidPost(t *testing.T){
	post := `{"id": "-1"}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/unlikePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UnlikePost(context))
}

func TestPost_CantUnlikeUnknownPost(t *testing.T){
	post := `{"id": -1}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/unlikePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UnlikePost(context))
}

func TestPost_CantLikeUnknownPost(t *testing.T){
	post := `{"id": -1}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.LikePost(context))
}

func TestPost_CantLikeTwicePost(t *testing.T){
	post := `{"id": `+strconv.Itoa(testPost.ID)+`}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.LikePost(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	req = httptest.NewRequest(echo.POST, "/likePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	assert.Error(t, c.LikePost(context))

	req = httptest.NewRequest(echo.POST, "/unlikePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	if assert.NoError(t, c.UnlikePost(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func TestPost_CantUnlikeTwicePost(t *testing.T){
	post := `{"id": `+strconv.Itoa(testPost.ID)+`}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.LikePost(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	
	req = httptest.NewRequest(echo.POST, "/unlikePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	if assert.NoError(t, c.UnlikePost(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(echo.POST, "/unlikePost", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	assert.Error(t, c.UnlikePost(context))


}


func TestPost_DeletePost(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/deletePost/"+strconv.Itoa(testPost.ID), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.DeletePost(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestPost_GetPostsAndComments(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/feed", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.GetPosts(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
		var posts []*models.Post
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &posts))
		assert.Equal(t, posts[1].Author.ID, test.ID)
		assert.Equal(t, posts[1].Comments[1].Author.ID, test.ID)
	}
}
