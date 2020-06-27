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
	tutorRepository interfaces.TutorRepository = repositories.NewTutorRepository()
	tutorService    service.TutorService       = service.NewTutorService(tutorRepository)
	tutorLogic      api.TutorLogic             = api.NewTutorLogic(tutorService)

	insRepository interfaces.InstitutionRepository = repositories.NewInstitutionRepository()
	insService    service.InstitutionService       = service.NewInstitutionService(insRepository)
	insLogic      api.InstitutionLogic             = api.NewInstitutionLogic(insService)
)

func main() {

	httpRouter.POSTLogin("/api/v1/login", api.Login)
	httpRouter.GET("/api/v1/references", api.GetReferences)
	httpRouter.GET("/api/v1/roles", api.GetRoles)
	httpRouter.GET("/api/v1/address", api.GetAddress)
	httpRouter.POST("/api/v1/userlogin", api.CreateUserLogin)

	httpRouter.GET("/api/v1/institution", insLogic.GetInstitution)
	httpRouter.GET("/api/v1/institutions", insLogic.GetInstitutions)
	httpRouter.POST("/api/v1/institutions", insLogic.CreateInstitutions)

	httpRouter.POST("/api/v1/tutors", tutorLogic.CreateTutors)
	httpRouter.GET("/api/v1/tutors", tutorLogic.GetTutors)
	httpRouter.GET("/api/v1/tutor", tutorLogic.GetTutor)
	httpRouter.POST("/api/v1/tutor/details", tutorLogic.UpdateTutorDetails)
	httpRouter.GET("/api/v1/tutor/details", tutorLogic.GetTutorDetails)
	httpRouter.POST("/api/v1/tutor/educations", tutorLogic.UpdateEducations)
	httpRouter.POST("/api/v1/tutor/experiences", tutorLogic.UpdateExperiences)
	httpRouter.POST("/api/v1/tutor/certificates", tutorLogic.UpdateCertificates)
	httpRouter.POST("/api/v1/tutor/journals", tutorLogic.UpdateJournals)
	httpRouter.POST("/api/v1/tutor/researches", tutorLogic.UpdateResearches)

	httpRouter.POST("/api/v1/students", api.CreateStudents)
	httpRouter.GET("/api/v1/students", api.GetStudents)
	httpRouter.GET("/api/v1/student", api.GetStudent)
	httpRouter.GET("/api/v1/student/details", api.GetStudentDetails)
	httpRouter.POST("/api/v1/student/details", api.UpdateStudentDetais)

	httpRouter.SERVE(":8000")
}
