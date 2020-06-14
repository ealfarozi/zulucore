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
	r.Handle("/api/v1/students", api.IsAuthorized(api.CreateStudents)).Methods("POST")
	r.Handle("/api/v1/students", api.IsAuthorized(api.GetStudents)).Methods("GET")
	r.Handle("/api/v1/student", api.IsAuthorized(api.GetStudent)).Methods("GET")
	r.Handle("/api/v1/student/details", api.IsAuthorized(api.GetStudentDetails)).Methods("GET")
	r.Handle("/api/v1/student/details", api.IsAuthorized(api.UpdateStudentDetais)).Methods("POST")

	r.Handle("/api/v1/tutors", api.IsAuthorized(api.CreateTutors)).Methods("POST")
	r.Handle("/api/v1/tutors", api.IsAuthorized(api.GetTutors)).Methods("GET")
	r.Handle("/api/v1/tutor", api.IsAuthorized(api.GetTutor)).Methods("GET")
	r.Handle("/api/v1/tutor/details", api.IsAuthorized(api.UpdateTutorDetails)).Methods("POST")
	r.Handle("/api/v1/tutor/details", api.IsAuthorized(api.GetTutorDetails)).Methods("GET")
	r.Handle("/api/v1/tutor/educations", api.IsAuthorized(api.UpdateEducations)).Methods("POST")
	r.Handle("/api/v1/tutor/certificates", api.IsAuthorized(api.UpdateCertificates)).Methods("POST")
	r.Handle("/api/v1/tutor/experiences", api.IsAuthorized(api.UpdateExperiences)).Methods("POST")
	r.Handle("/api/v1/tutor/journals", api.IsAuthorized(api.UpdateJournals)).Methods("POST")
	r.Handle("/api/v1/tutor/researches", api.IsAuthorized(api.UpdateResearches)).Methods("POST")

	r.Handle("/api/v1/institution", api.IsAuthorized(api.GetInstitution)).Methods("GET")
	r.Handle("/api/v1/institutions", api.IsAuthorized(api.GetInstitutions)).Methods("GET")
	r.Handle("/api/v1/institutions", api.IsAuthorized(api.CreateInstitutions)).Methods("POST")

	r.Handle("/api/v1/references", api.IsAuthorized(api.GetReferences)).Methods("GET")
	r.Handle("/api/v1/roles", api.IsAuthorized(api.GetRoles)).Methods("GET")
	r.Handle("/api/v1/address", api.IsAuthorized(api.GetAddress)).Methods("GET")
	r.Handle("/api/v1/userlogin", api.IsAuthorized(api.CreateUserLogin)).Methods("POST")
	r.HandleFunc("/api/v1/login", api.Login).Methods("POST")

	//3 - Run the server
	log.Fatal(http.ListenAndServe(":8000", r))
}
