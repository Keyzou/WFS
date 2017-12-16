package controllers

import(
	"github.com/labstack/echo"
	"net/http"
	"models"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"fmt"
	"strconv"
)

func (c *Controller) CreateComment(context echo.Context) (err error){
	user := UserIDFromToken(context)
	u := &models.User{ID: user}	
	comment := &models.Comment{}

	err = context.Bind(comment)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid message"}
	}

	if comment.Content == ""{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Content cannot be empty !"}
	}

	var _id int
	var username, email string
	err = c.DB.QueryRow(`select id, username, email from users where id=$1`, u.ID).Scan(&_id, &username, &email)

	u.Username = username
	u.Email = email

	comment.Author = *u

	var id int
	var date pgtype.Timestamptz
	err = c.DB.QueryRow(`insert into comments(author, content, date, likes, post_id) values($1, $2, now(), 0, $3) returning id, date`, comment.Author.ID, comment.Content, comment.PostID).Scan(&id, &date)
	comment.ID = id
	comment.Date = date.Time.Format("_2 January 2006 at 15:04:05")

	return context.JSON(http.StatusCreated, comment)
}


func (c *Controller) DeleteComment(context echo.Context) (err error) {
	u := UserIDFromToken(context)
	id, _ := strconv.Atoi(context.Param("id"))
	comment := &models.Comment{ID: id}

	c.DB.Exec(`delete from comments where author=$1 and id=$2`, u, comment.ID)
	
	return context.NoContent(http.StatusOK)
}

func (c *Controller) LikeComment(context echo.Context) (err error){
	u := UserIDFromToken(context)
	toLike := &models.Comment{}

	err = context.Bind(toLike)
	if err != nil {
		fmt.Println(err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "error"}
	}

	err = c.DB.QueryRow(`select id from comments where id=$1`, toLike.ID).Scan(&toLike.ID)
	if err != nil {
		if err == pgx.ErrNoRows{
			return echo.ErrNotFound
		}
	}

	liked := c.HasLikedComment(u, toLike.ID)
	count := c.GetCommentLikes(toLike.ID)
	if liked {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "already liked"}
	}

	c.DB.Exec(`insert into comments_likes(user_id, comment_id) values($1, $2)`, u, toLike.ID)

	toLike.Likes = count + 1
	toLike.Liked = true

	c.DB.Exec(`update comments set likes=$1 where id=$2`, toLike.Likes, toLike.ID)

	return context.JSON(http.StatusOK, toLike)
}

func (c *Controller) UnlikeComment(context echo.Context) (err error){
	u := UserIDFromToken(context)
	toLike := &models.Comment{}

	err = context.Bind(toLike)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "error"}
	}

	err = c.DB.QueryRow(`select id from comments where id=$1`, toLike.ID).Scan(&toLike.ID)
	if err != nil {
		if err == pgx.ErrNoRows{
			return echo.ErrNotFound
		}
	}

	liked := c.HasLikedComment(u, toLike.ID)
	count := c.GetCommentLikes(toLike.ID)
	if !liked {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "already disliked"}
	}

	c.DB.Exec(`delete from comments_likes where user_id=$1 and comment_id=$2`, u, toLike.ID)

	toLike.Likes = count - 1
	toLike.Liked = false

	_, err = c.DB.Exec(`update comments set likes=$1 where id=$2`, toLike.Likes, toLike.ID)

	return context.JSON(http.StatusOK, toLike)
}

func (c *Controller) HasLikedComment(userid int, commentid int) (liked bool) {

	var likes int
	c.DB.QueryRow(`select count(*) from comments_likes where user_id=$1 and comment_id=$2`, userid, commentid).Scan(&likes)
	if likes > 0 {
		liked = true
		return liked
	}
	liked = false
	return liked
}

func (c *Controller) GetCommentLikes(commentid int) (likes int) {
	c.DB.QueryRow(`select count(*) from comments_likes where comment_id=$1`, commentid).Scan(&likes)

	return likes
}

func (c *Controller) GetComments(userid int, idPost int) (comments []*models.Comment){
	comments = []*models.Comment{}
	rows, _ := c.DB.Query(`SELECT comments.id, comments.author, comments.content, comments.date, comments.post_id, users.email, users.username FROM comments JOIN users on comments.author=users.id AND comments.post_id=$1 ORDER BY date DESC`, idPost)
	defer rows.Close()
	for rows.Next() {
		var id, author int
		var post int
		var content string
		var date pgtype.Timestamptz
		var userEmail string
		var userName string
		rows.Scan(&id, &author, &content, &date, &post, &userEmail, &userName)

		u := models.User{ID: author, Email: userEmail, Username: userName}

		time := date.Time.Format("_2 January 2006 at 15:04:05")
		comment := &models.Comment{ID: id, Author: u, Content: content, Date: time, PostID: post}
		comments = append(comments, comment)
	}
	
	for _, el := range comments {
		el.Liked = c.HasLikedComment(userid, el.ID)
		el.Likes = c.GetCommentLikes(el.ID)
		el.Author.Followers = c.GetFollowersCount(el.Author.ID)
		el.IsFollowingAuthor = c.isFollowing(userid, el.Author.ID)
	}
	return comments
}