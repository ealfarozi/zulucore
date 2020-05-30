package main

import (
	"log"
	"net/http"

	"github.com/ealfarozi/zulucore/logic/api"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	//2 - Route Handlers/Endpoints/URL endpoints/Service URL/Path
	r.Handle("/api/v1/tutors", api.IsAuthorized(api.GetTutors)).Methods("GET")
	r.Handle("/api/v1/tutors/{tutorid}", api.IsAuthorized(api.GetTutor)).Methods("GET")
	r.Handle("/api/v1/userLogin", api.IsAuthorized(api.CreateUserLogin)).Methods("POST")
	r.HandleFunc("/api/v1/login", api.Login).Methods("POST")

	//3 - Run the server
	log.Fatal(http.ListenAndServe(":80", r))
}
