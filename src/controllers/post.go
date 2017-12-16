package controllers

import(
	"models"
	"github.com/labstack/echo"
	"net/http"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"fmt"
	"html"
	"strings"
	"strconv"
)

func (c *Controller) CreatePost(co echo.Context) (err error) {
	
	u := &models.User{ID: UserIDFromToken(co)}	
	p := &models.Post{Author: *u}
	err = co.Bind(p)
	fmt.Println(p)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid request"}
	}

	p.Content = strings.TrimSpace(p.Content)

	if p.Content == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid message"}
	}

	p.Content = html.EscapeString(p.Content)

	fmt.Println(u.ID)

	var _id int
	var username, email string
	err = c.DB.QueryRow(`select id, username, email from users where id=$1`, u.ID).Scan(&_id, &username, &email)

	u.Username = username
	u.Email = email

	p.Author = *u

	var id int
	var date pgtype.Timestamptz
	err = c.DB.QueryRow(`insert into posts(author, content, date, likes) values($1, $2, now(), 0) returning id, date`, p.Author.ID, p.Content).Scan(&id, &date)
	p.ID = id
	p.Date = date.Time.Format("_2 January 2006 at 15:04:05")

	return co.JSON(http.StatusCreated, p)
}

func (c *Controller) LikePost(context echo.Context) (err error){
	u := UserIDFromToken(context)
	toLike := &models.Post{}

	err = context.Bind(toLike)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid request"}
	}

	err = c.DB.QueryRow(`select id from posts where id=$1`, toLike.ID).Scan(&toLike.ID)
	if err != nil {
		if err == pgx.ErrNoRows{
			return echo.ErrNotFound
		}
	}

	liked := c.HasLikedPost(u, toLike.ID)
	count := c.GetPostLikes(toLike.ID)
	if liked {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "already liked"}
	}

	c.DB.Exec(`insert into posts_likes(user_id, post_id) values($1, $2)`, u, toLike.ID)

	toLike.Likes = count + 1
	toLike.Liked = true

	c.DB.Exec(`update posts set likes=$1 where id=$2`, toLike.Likes, toLike.ID)

	return context.JSON(http.StatusOK, toLike)
}

func (c *Controller) UnlikePost(context echo.Context) (err error){
	u := UserIDFromToken(context)
	toLike := &models.Post{}

	err = context.Bind(toLike)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid request"}
	}	

	err = c.DB.QueryRow(`select id from posts where id=$1`, toLike.ID).Scan(&toLike.ID)
	if err != nil {
		if err == pgx.ErrNoRows{
			return echo.ErrNotFound
		}
	}

	liked := c.HasLikedPost(u, toLike.ID)
	count := c.GetPostLikes(toLike.ID)
	if !liked {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "already disliked"}
	}

	c.DB.Exec(`delete from posts_likes where user_id=$1 and post_id=$2`, u, toLike.ID)


	toLike.Likes = count - 1
	toLike.Liked = false

	c.DB.Exec(`update posts set likes=$1 where id=$2`, toLike.Likes, toLike.ID)

	return context.JSON(http.StatusOK, toLike)
}

func (c *Controller) DeletePost(context echo.Context) (err error) {
	u := UserIDFromToken(context)
	id, _ := strconv.Atoi(context.Param("id"))
	post := &models.Post{ID: id}

	fmt.Println(post.ID)
	c.DB.Exec(`delete from posts where author=$1 and id=$2`, u, post.ID)

	return context.NoContent(http.StatusOK)
}

func (c *Controller) HasLikedPost(userid int, postid int) (liked bool) {

	var likes int
	c.DB.QueryRow(`select count(*) from posts_likes where user_id=$1 and post_id=$2`, userid, postid).Scan(&likes)
	if likes > 0 {
		liked = true
		return liked
	}
	liked = false
	return liked
}

func (c *Controller) GetPostLikes(postid int) (likes int) {
	c.DB.QueryRow(`select count(*) from posts_likes where post_id=$1`, postid).Scan(&likes)
	return likes
}
	

func (c *Controller) GetPosts(context echo.Context) (err error){
	u := &models.User{ID: UserIDFromToken(context)}	

	posts := []*models.Post{}
	rows, _ := c.DB.Query(`SELECT posts.id, posts.author, posts.content, posts.date, users.email, users.username FROM posts LEFT JOIN users on posts.author=users.id AND (posts.author=$1 OR posts.author IN (SELECT user_followed FROM followers WHERE user_following=$1)) ORDER BY date DESC`, u.ID)
	defer rows.Close()
	for rows.Next() {
		var id, author int
		var content string
		var date pgtype.Timestamptz
		var userEmail string
		var userName string
		rows.Scan(&id, &author, &content, &date, &userEmail, &userName)
		u := models.User{ID: author, Email: userEmail, Username: userName}

		time := date.Time.Format("_2 January 2006 at 15:04:05")
		post := &models.Post{ID: id, Author: u, Content: content, Date: time}
		fmt.Println(post)
		posts = append(posts, post)
		fmt.Println(posts)
	}


	for _, el := range posts {
		el.Comments = c.GetComments(u.ID, el.ID)
		el.Author.Followers = c.GetFollowersCount(el.Author.ID)
		el.Liked = c.HasLikedPost(u.ID, el.ID)
		el.Likes = c.GetPostLikes(el.ID)
		el.IsFollowingAuthor = c.isFollowing(u.ID, el.Author.ID)
	}	


	return context.JSON(http.StatusOK, posts)
}