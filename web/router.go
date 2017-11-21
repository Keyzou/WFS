package web

import(
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/gorilla/mux"
)

func addRoutes(router *mux.Router){
	router.HandleFunc("/", rootHandler)
}


func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	data, err := ioutil.ReadFile("html/index.html")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, string(data))
}
