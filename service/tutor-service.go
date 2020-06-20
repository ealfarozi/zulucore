package service

import (
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/structs"
)

type TutorService interface {
	Validate(*structs.Tutor) (*structs.Tutor, *structs.ErrorMessage)
	Validates(*[]structs.Tutor) (*[]structs.Tutor, *structs.ErrorMessage)
	GetTutors(insID string) (*[]structs.Tutor, *structs.ErrorMessage)
	GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage)
}

type service struct{}

var (
	repo interfaces.TutorRepository
)

func NewTutorService(repository interfaces.TutorRepository) TutorService {
	repo = repository
	return &service{}
}

func (*service) GetTutors(insID string) (*[]structs.Tutor, *structs.ErrorMessage) {
	return repo.GetTutors(insID)
}

func (*service) GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage) {
	return repo.GetTutorDetails(tutorID)
}

func (*service) Validate(*structs.Tutor) (*structs.Tutor, *structs.ErrorMessage) {
	return nil, nil
}

func (*service) Validates(*[]structs.Tutor) (*[]structs.Tutor, *structs.ErrorMessage) {
	return nil, nil
}
