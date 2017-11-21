package main

import (
	"log"
	"net/http"
	"github.com/keyzou/spt-go/web"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	web.addRoutes(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
