package interfaces

import (
	"github.com/ealfarozi/zulucore/structs"
)

type TutorRepository interface {
	GetTutors(insID string, page string, limit string) (*[]structs.Tutor, *structs.ErrorMessage)
	GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage)
	GetTutor(nmrInd string, name string, insID string, page string, limit string) (*[]structs.Tutor, *structs.ErrorMessage)
	UpdateTutorDetails(tutor structs.Tutor) *structs.ErrorMessage
	CreateTutors(tutor structs.Tutor) *structs.ErrorMessage
	UpdateEducations(edu structs.TutorEducation) *structs.ErrorMessage
	UpdateExperiences(exp structs.TutorExperience) *structs.ErrorMessage
	UpdateCertificates(cert structs.TutorCertificate) *structs.ErrorMessage
	UpdateJournals(jour structs.TutorJournal) *structs.ErrorMessage
	UpdateResearches(res structs.TutorResearch) *structs.ErrorMessage
	CheckNomorInduk(insID int, nmrInduk string, tutorID int) int
	CheckEmail(email string, usrID int) int
}
