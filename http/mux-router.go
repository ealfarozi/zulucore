package router

import (
	"log"
	"net/http"

	"github.com/ealfarozi/zulucore/logic/api"
	"github.com/gorilla/mux"
)

type muxRouter struct{}

var (
	muxDispatcher = mux.NewRouter()
)

func NewMuxRouter() Router {
	return &muxRouter{}
}

func (*muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.Handle(uri, api.IsAuthorized(f)).Methods("GET")
}
func (*muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.Handle(uri, api.IsAuthorized(f)).Methods("POST")
}
func (*muxRouter) SERVE(port string) {
	log.Fatal(http.ListenAndServe(port, muxDispatcher))
}

func (*muxRouter) POSTLogin(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}
