package main

import (
	router "github.com/ealfarozi/zulucore/http"
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/logic/api"
	"github.com/ealfarozi/zulucore/repositories"
	"github.com/ealfarozi/zulucore/service"
)

var (
	httpRouter      router.Router              = router.NewMuxRouter()
	tutorRepository interfaces.TutorRepository = repositories.NewMysqlRepository()
	tutorService    service.TutorService       = service.NewTutorService(tutorRepository)
	tutorLogic      api.TutorLogic             = api.NewTutorLogic(tutorService)
)

func main() {

	httpRouter.POSTLogin("/api/v1/login", api.Login)
	httpRouter.GET("/api/v1/references", api.GetReferences)
	httpRouter.GET("/api/v1/roles", api.GetRoles)
	httpRouter.GET("/api/v1/address", api.GetAddress)
	httpRouter.POST("/api/v1/userlogin", api.CreateUserLogin)

	httpRouter.GET("/api/v1/institution", api.GetInstitution)
	httpRouter.GET("/api/v1/institutions", api.GetInstitutions)
	httpRouter.POST("/api/v1/institutions", api.CreateInstitutions)

	httpRouter.POST("/api/v1/tutors", api.CreateTutors)
	httpRouter.GET("/api/v1/tutors", tutorLogic.GetTutors)
	httpRouter.GET("/api/v1/tutor", api.GetTutor)
	httpRouter.POST("/api/v1/tutor/details", api.UpdateTutorDetails)
	httpRouter.GET("/api/v1/tutor/details", tutorLogic.GetTutorDetails)
	httpRouter.POST("/api/v1/tutor/educations", api.UpdateEducations)
	httpRouter.POST("/api/v1/tutor/experiences", api.UpdateExperiences)
	httpRouter.POST("/api/v1/tutor/certificates", api.UpdateCertificates)
	httpRouter.POST("/api/v1/tutor/journals", api.UpdateJournals)
	httpRouter.POST("/api/v1/tutor/researches", api.UpdateResearches)

	httpRouter.POST("/api/v1/students", api.CreateStudents)
	httpRouter.GET("/api/v1/students", api.GetStudents)
	httpRouter.GET("/api/v1/student", api.GetStudent)
	httpRouter.GET("/api/v1/student/details", api.GetStudentDetails)
	httpRouter.POST("/api/v1/student/details", api.UpdateStudentDetais)

	httpRouter.SERVE(":8000")
}
