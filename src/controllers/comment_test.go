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

func TestCreateComment(t *testing.T){
	post := `{"content": "test", "postId": `+strconv.Itoa(testPost.ID)+`}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/create-comment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.CreateComment(context)){
		if assert.Equal(t, http.StatusCreated, rec.Code){
			newPost := &models.Comment{}
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &newPost))
		}
	}
}

func TestCreateInvalidCommentContent(t *testing.T){
	post := `{"content": "", "postId": `+strconv.Itoa(testPost.ID)+`}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/create-comment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.CreateComment(context))
}

func TestCreateInvalidCommentPostID(t *testing.T){
	post := `{"content": "qweqweqwe", "postId": "`+strconv.Itoa(testPost.ID)+`"}`
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/create-comment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.CreateComment(context))
}

func TestLikeAndUnlikeComment(t *testing.T){
	comment := `{"id": `+strconv.Itoa(testComment.ID)+`}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likeComment", strings.NewReader(comment))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.LikeComment(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, c.HasLikedComment(test.ID, testComment.ID))
	}
	req = httptest.NewRequest(echo.POST, "/unlikeComment", strings.NewReader(comment))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	if assert.NoError(t, c.UnlikeComment(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.False(t, c.HasLikedComment(test.ID, testComment.ID))
	}
}


func TestCantLikeInvalidComment(t *testing.T){
	comment := `{"id": "-1"}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likeComment", strings.NewReader(comment))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.LikeComment(context))
}

func TestCantLikeUnknownComment(t *testing.T){
	comment := `{"id": -1}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likeComment", strings.NewReader(comment))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.LikeComment(context))
}

func TestCantUnlikeUnknownComment(t *testing.T){
	comment := `{"id": -1}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/unlikeComment", strings.NewReader(comment))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UnlikeComment(context))
}

func TestCantUnlikeInvalidComment(t *testing.T){
	comment := `{"id": "-1"}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/unlikeComment", strings.NewReader(comment))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	assert.Error(t, c.UnlikeComment(context))
}

func TestCantLikeTwiceComment(t *testing.T){
	post := `{"id": `+strconv.Itoa(testComment.ID)+`}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likeComment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.LikeComment(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	req = httptest.NewRequest(echo.POST, "/likeComment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	assert.Error(t, c.LikeComment(context))

	req = httptest.NewRequest(echo.POST, "/unlikeComment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	if assert.NoError(t, c.UnlikeComment(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCantUnlikeTwiceComment(t *testing.T){
	post := `{"id": `+strconv.Itoa(testComment.ID)+`}`

	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/likeComment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.LikeComment(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(echo.POST, "/unlikeComment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)
	
	if assert.NoError(t, c.UnlikeComment(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(echo.POST, "/unlikeComment", strings.NewReader(post))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec = httptest.NewRecorder()
	context = e.NewContext(req, rec)

	assert.Error(t, c.UnlikeComment(context))

}

func TestDeleteComment(t *testing.T){
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/deleteComment/"+strconv.Itoa(testComment.ID), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+test.Token)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	c := &controllers.Controller{DB: conn}

	if assert.NoError(t, c.DeleteComment(context)){
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}