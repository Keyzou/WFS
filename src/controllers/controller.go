package controllers

import(
	"github.com/jackc/pgx"
)

type(
	Controller struct {
		DB *pgx.ConnPool
	}
)


const(
	JWTSecretKey = "secretkey"
)
