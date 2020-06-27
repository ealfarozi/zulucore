package service

import (
	"net/http"

	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/structs"
	"gopkg.in/go-playground/validator.v9"
)

type TutorService interface {
	Validate(tutor *structs.Tutor) (*structs.Tutor, *structs.ErrorMessage)
	ValidateEdu(edu *structs.TutorEducation) *structs.ErrorMessage
	ValidateExp(exp *structs.TutorExperience) *structs.ErrorMessage
	ValidateCert(cert *structs.TutorCertificate) *structs.ErrorMessage
	ValidateJour(jour *structs.TutorJournal) *structs.ErrorMessage
	ValidateRes(res *structs.TutorResearch) *structs.ErrorMessage
	Validates(tutors *[]structs.Tutor) (*structs.Tutor, *structs.ErrorMessage)
	UpdateTutorDetails(tutors structs.Tutor) *structs.ErrorMessage
	CreateTutors(tutor structs.Tutor) *structs.ErrorMessage
	UpdateEducations(edu structs.TutorEducation) *structs.ErrorMessage
	UpdateExperiences(exp structs.TutorExperience) *structs.ErrorMessage
	UpdateJournals(jour structs.TutorJournal) *structs.ErrorMessage
	UpdateCertificates(cert structs.TutorCertificate) *structs.ErrorMessage
	UpdateResearches(res structs.TutorResearch) *structs.ErrorMessage
	GetTutors(insID string) (*[]structs.Tutor, *structs.ErrorMessage)
	GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage)
	GetTutor(nmrInd string, name string, insID string) (*[]structs.Tutor, *structs.ErrorMessage)
	CheckNomorInduk(insID int, nmrInduk string, tutorID int) int
	CheckEmail(email string, usrID int) int
}

type service struct{}

var (
	repo interfaces.TutorRepository
)

func NewTutorService(repository interfaces.TutorRepository) TutorService {
	repo = repository
	return &service{}
}

func (*service) GetTutor(nmrInd string, name string, insID string) (*[]structs.Tutor, *structs.ErrorMessage) {
	return repo.GetTutor(nmrInd, name, insID)
}

func (*service) GetTutors(insID string) (*[]structs.Tutor, *structs.ErrorMessage) {
	return repo.GetTutors(insID)
}

func (*service) GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage) {
	return repo.GetTutorDetails(tutorID)
}

func (*service) UpdateTutorDetails(tutor structs.Tutor) *structs.ErrorMessage {
	return repo.UpdateTutorDetails(tutor)
}

func (*service) UpdateEducations(edu structs.TutorEducation) *structs.ErrorMessage {
	return repo.UpdateEducations(edu)
}

func (*service) UpdateJournals(jour structs.TutorJournal) *structs.ErrorMessage {
	return repo.UpdateJournals(jour)
}

func (*service) UpdateCertificates(cert structs.TutorCertificate) *structs.ErrorMessage {
	return repo.UpdateCertificates(cert)
}

func (*service) UpdateExperiences(exp structs.TutorExperience) *structs.ErrorMessage {
	return repo.UpdateExperiences(exp)
}

func (*service) UpdateResearches(res structs.TutorResearch) *structs.ErrorMessage {
	return repo.UpdateResearches(res)
}

func (*service) CreateTutors(tutor structs.Tutor) *structs.ErrorMessage {
	return repo.CreateTutors(tutor)
}

func (*service) CheckNomorInduk(insID int, nmrInduk string, tutorID int) int {
	return repo.CheckNomorInduk(insID, nmrInduk, tutorID)
}

func (*service) CheckEmail(email string, usrID int) int {
	return repo.CheckEmail(email, usrID)
}

func (*service) Validate(tutor *structs.Tutor) (*structs.Tutor, *structs.ErrorMessage) {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(tutor)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = tutor.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
	return tutor, nil
}

func (*service) ValidateCert(cert *structs.TutorCertificate) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(cert)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = cert.CertName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}
	return nil
}

func (*service) ValidateEdu(edu *structs.TutorEducation) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(edu)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = edu.UnivName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}
	return nil
}

func (*service) ValidateJour(jour *structs.TutorJournal) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(jour)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = jour.JourName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}
	return nil
}

func (*service) ValidateExp(exp *structs.TutorExperience) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(exp)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = exp.ExpName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}
	return nil
}

func (*service) ValidateRes(res *structs.TutorResearch) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(res)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = res.ResName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}
	return nil
}

func (*service) Validates(tutors *[]structs.Tutor) (*structs.Tutor, *structs.ErrorMessage) {
	return nil, nil
}
