package controllers

import(
	"models"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx"
	"time"
	"fmt"
	"strconv"
	"strings"
	"github.com/asaskevich/govalidator"
)

func (c *Controller) Signup(context echo.Context) (err error) {
	 u := &models.User{}
	 err = context.Bind(u); 
	 if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid request"}
	 }

	 u.Email = strings.TrimSpace(u.Email)
	 u.Username = strings.TrimSpace(u.Username)
	 u.Password = strings.TrimSpace(u.Password)

	 if u.Email == "" || u.Username == "" || u.Password == "" {
		 return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email, password, or username !"}
	 }

	 if !govalidator.IsEmail(u.Email) || !govalidator.IsAlphanumeric(u.Username){
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email, password, or username !"}
	 }

	 u.Email = strings.ToLower(u.Email)
	 u.Username = strings.ToLower(u.Username)

	 if len(u.Username) < 3 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username must be at least 3 characters long !"}
	 }
	 if len(u.Password) < 6 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Password must be at least 6 characters long !"}
	 }
		 
	 if !c.isUsernameUnique(u.Username){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username already taken !"}
	}

	if !c.isEmailUnique(u.Email){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Email already taken !"}
	}

	u.HashPassword()

	var id int
	err = c.DB.QueryRow(`insert into users (username, email, password) VALUES($1, $2, $3) returning id`, u.Username, u.Email, u.Password).Scan(&id)
	u.ID = id
	u.Password = ""
	return context.JSON(http.StatusCreated, u)
}


func (c *Controller) Login(context echo.Context) (err error){
	u := &models.User{}
	err = context.Bind(u)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid request"}
	}

	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	if u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password "}
	}

	if !govalidator.IsEmail(u.Email){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password "}
	}

	u.Email = strings.ToLower(u.Email)
		

	var id int
	var username, password string
	err = c.DB.QueryRow(`select id, username, password from users where email=$1`, u.Email).Scan(&id, &username, &password)
	if err != nil {
		if err == pgx.ErrNoRows{
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password)) != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["expiring"] = time.Now().Add(time.Hour * 72).Unix()
	u.Email = ""
	u.Password = ""
	u.Token, err = token.SignedString([]byte(JWTSecretKey))
	return context.JSON(http.StatusOK, u)
}

func (c *Controller) UpdateUser(context echo.Context) (err error){
	id := UserIDFromToken(context)
	u := &models.User{}
	err = context.Bind(u)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid request"}
	}

	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)

	if !govalidator.IsAlphanumeric(u.Username){
		 return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid username !"}
	}

	u.Username = strings.ToLower(u.Username)

	if u.Password == "" && u.Username == ""{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "no changes"}
	}

	if len(u.Username) < 3 && u.Username != "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username must be at least 3 characters long !"}
	 }
	 if len(u.Password) < 6 && u.Password != "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Password must be at least 6 characters long !"}
	 }
	

	if !c.isUsernameUnique(u.Username){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username already taken !"}
	}

	u.HashPassword()

	if u.Username != ""{
		c.DB.Exec(`update users set username=$1 where id=$2`, u.Username, id)
	}
	if u.Password != ""{
		c.DB.Exec(`update users set password=$1 where id=$2`, u.Password, id)
	}

	return context.JSON(http.StatusOK, &models.User{Username: u.Username})
}

func (c *Controller) QueryUsers(context echo.Context) (err error){
	query := context.Param("query")
	users := []*models.User{}

	query = strings.TrimSpace(query)

	if query == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "no query"}
	}

	rows, err := c.DB.Query(`select id, username from users where username like '%'||$1||'%' or email like '%'||$1||'%'`, query)

	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		rows.Scan(&id, &username)
		u := &models.User{ID: id, Username: username}

		users = append(users, u)
	}

	return context.JSON(http.StatusOK, users)
}

func (c *Controller) Follow(context echo.Context) (err error){
	u := UserIDFromToken(context)
	if !govalidator.IsInt(context.Param("id")) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "id must be int"}
	}
	toFollow, _ := strconv.Atoi(context.Param("id"))
	if (toFollow == u){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "cant follow yourself"}
	}
	err = c.DB.QueryRow(`select id from users where id=$1`, toFollow).Scan(&toFollow)
	if err != nil {
		if err == pgx.ErrNoRows{
			return echo.ErrNotFound
		}
	}

	var count int
	err = c.DB.QueryRow(`select count(*) from followers where user_followed=$1 and user_following=$2`, toFollow, u).Scan(&count)

	if count > 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "already following"}
	}

	c.DB.Exec(`insert into followers(user_followed, user_following) values($1, $2)`, toFollow, u)


	return context.JSON(http.StatusOK, &models.User{Followers: count+1})
}

func (c *Controller) isUsernameUnique(username string) (usernameUnique bool){
		var count int
		c.DB.QueryRow(`select count(*) from users where username=$1`, username).Scan(&count)
		return count == 0
	}

func (c *Controller) isEmailUnique(email string) (emailUnique bool){
		var count int
		c.DB.QueryRow(`select count(*) from users where email LIKE '%'||$1||'%'`, email).Scan(&count)
		return count == 0
	}

func (c *Controller) GetFollowersCount(userid int) (followers int){

	c.DB.QueryRow(`select count(*) from followers where user_followed=$1`, userid).Scan(&followers)
	return followers
}

func (c *Controller) GetFollowers(context echo.Context) (err error){
	
		if !govalidator.IsInt(context.Param("id")) {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: "id must be int"}
		}
		u, err := strconv.Atoi(context.Param("id"))
	
		return context.JSON(http.StatusOK, c.GetFollowersCount(u))
		
	}

func (c *Controller) isFollowing(u1 int, u2 int) (following bool){
		var count int
		c.DB.QueryRow(`select count(*) from followers where user_followed=$2 and user_following=$1`, u1, u2).Scan(&count)
		fmt.Println(count)
		return count > 0
	}
	
	

func (c *Controller) Unfollow(context echo.Context) (err error){
	u := UserIDFromToken(context)
	if !govalidator.IsInt(context.Param("id")) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "id must be int"}
	}
	toFollow, _ := strconv.Atoi(context.Param("id"))

	if (toFollow == u){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "cant unfollow yourself"}
	}

	err = c.DB.QueryRow(`select id from users where id=$1`, toFollow).Scan(&toFollow)
	if err != nil {
		if err == pgx.ErrNoRows{
			return echo.ErrNotFound
		}
	}

	var count int
	err = c.DB.QueryRow(`select count(*) from followers where user_followed=$1 and user_following=$2`, toFollow, u).Scan(&count)

	if count <= 0 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "not following"}
	}

	c.DB.Exec(`delete from followers where user_followed=$1 and user_following=$2`, toFollow, u)

	return context.JSON(http.StatusOK, &models.User{Followers: count-1})
}


func (c *Controller) GetFollowingUsers(context echo.Context) (err error){
	
		user := UserIDFromToken(context)

		users := []*models.User{}

		rows, _ := c.DB.Query(`SELECT users.id, users.email, users.username FROM users JOIN followers ON followers.user_following=$1 AND followers.user_followed=users.id`, user)
		defer rows.Close()
		
		for rows.Next(){
			var id int
			var email, username string
			rows.Scan(&id, &email, &username)
			u := &models.User{ID: id, Email: email, Username: username}

			users = append(users, u)
		}
		
	
		return context.JSON(http.StatusOK, users)
}

func (c *Controller) GetFollowersUsers(context echo.Context) (err error){
	
		user := UserIDFromToken(context)

		users := []*models.User{}

		rows, _ := c.DB.Query(`SELECT users.id, users.email, users.username FROM users JOIN followers ON followers.user_followed=$1 AND followers.user_following=users.id`, user)
		defer rows.Close()
		
		for rows.Next(){
			var id int
			var email, username string
			rows.Scan(&id, &email, &username)
			u := &models.User{ID: id, Email: email, Username: username}

			users = append(users, u)
		}
		
	
		return context.JSON(http.StatusOK, users)
}
	

func (c *Controller) GetFollowing(context echo.Context) (err error){

	u1 := UserIDFromToken(context)
	if !govalidator.IsInt(context.Param("id")) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "id must be int"}
	}
	u2, err := strconv.Atoi(context.Param("id"))

	return context.JSON(http.StatusOK, c.isFollowing(u1, u2))
	
}

func (c *Controller) GetUser(context echo.Context) (err error){
	if !govalidator.IsInt(context.Param("id")) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "id must be int"}
	}
	u, err := strconv.Atoi(context.Param("id"))

	user := &models.User{ID: u}
	var email, username string
	err = c.DB.QueryRow(`SELECT email, username FROM users WHERE users.id=$1`, u).Scan(&email, &username)
	if err != nil {
		fmt.Println(err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}
	user.Email = email
	user.Username = username
	user.Followers = c.GetFollowersCount(user.ID)
	return context.JSON(http.StatusOK, user)
	
}

func UserIDFromToken(c echo.Context) int{
	u := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", -1)
	type Claims struct{
		ID float64 `json:"id"`
		Username string `json:"username"`
		Expiring int `json:"expiring"`
		jwt.StandardClaims
	}
	token, err := jwt.ParseWithClaims(u, &Claims{},  func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})
	if err != nil{
		fmt.Println(err)
		return -1
	}
	claims := token.Claims.(*Claims)
	return int(claims.ID)
}