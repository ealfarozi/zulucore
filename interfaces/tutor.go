package interfaces

import (
	"github.com/ealfarozi/zulucore/structs"
)

type TutorRepository interface {
	GetTutors(insID string) (*[]structs.Tutor, *structs.ErrorMessage)
	GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage)
	GetTutor(nmrInd string, name string, insID string) (*[]structs.Tutor, *structs.ErrorMessage)
	UpdateTutorDetails(tutor structs.Tutor) *structs.ErrorMessage
	CreateTutors(tutor structs.Tutor) *structs.ErrorMessage
	CheckNomorInduk(insID int, nmrInduk string, tutorID int) int
	CheckEmail(email string, usrID int) int
}
