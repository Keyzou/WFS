package main

import (
	"os"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"controllers"
	"github.com/jackc/pgx"
)

var dbConfig = pgx.ConnPoolConfig{MaxConnections: 15, ConnConfig: pgx.ConnConfig{Host: "localhost", User: "postgres", Password: "postgres", Database: "spt"}}

func main(){
	e := echo.New()
	e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(controllers.JWTSecretKey),
		TokenLookup: "header:authorization",
		Skipper: func(c echo.Context) bool {
			return true
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.Use(middleware.Secure())

	conn, err := pgx.NewConnPool(dbConfig)
	if err != nil {
		e.Logger.Fatal(err)
		os.Exit(1)
	}

	c := &controllers.Controller{DB: conn}

	e.POST("/signup", c.Signup)
	e.POST("/login", c.Login)

	e.POST("/follow/:id", c.Follow)
	e.POST("/unfollow/:id", c.Unfollow)
	e.GET("/queryUsers/:query", c.QueryUsers)
	e.GET("/user/:id", c.GetUser)
	e.GET("/user/:id/following", c.GetFollowing)
	e.GET("/user/:id/followers", c.GetFollowers)
	e.GET("/user/followers", c.GetFollowersUsers)
	e.GET("/user/following", c.GetFollowingUsers)
	e.POST("/user/update", c.UpdateUser)

	e.POST("/create-post", c.CreatePost)
	e.GET("/feed", c.GetPosts)
	e.POST("/likePost", c.LikePost)
	e.POST("/unlikePost", c.UnlikePost)

	e.DELETE("/deletePost/:id", c.DeletePost)
	e.DELETE("/deleteComment/:id", c.DeleteComment)

	e.POST("/create-comment", c.CreateComment)
	e.POST("/likeComment", c.LikeComment)
	e.POST("/unlikeComment", c.UnlikeComment)

	e.Logger.Fatal(e.Start(":1323"))
}