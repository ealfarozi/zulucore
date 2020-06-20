package interfaces

import (
	"github.com/ealfarozi/zulucore/structs"
)

type TutorRepository interface {
	GetTutors(insID string) (*[]structs.Tutor, *structs.ErrorMessage)
	GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage)
}
