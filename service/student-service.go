package service

import (
	"net/http"

	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/structs"
	"gopkg.in/go-playground/validator.v9"
)

type StudentService interface {
	CreateStudents(std structs.Student) *structs.ErrorMessage
	CreateParent(prt structs.Parents) *structs.ErrorMessage
	CheckEmail(email string, usrID int) int
	CheckFamily(famID int, stdID int) int
	CheckNomorIndukStd(insID int, nmrInduk string, stdID int) int
	UpdateStudentDetails(std structs.Student) *structs.ErrorMessage
	Validate(std *structs.Student) (*structs.Student, *structs.ErrorMessage)
	ValidateParent(prt *structs.Parents) (*structs.Parents, *structs.ErrorMessage)
	GetStudentDetails(stdID string) (*structs.Student, *structs.ErrorMessage)
	GetStudents(insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage)
	GetStudent(nmrInduk string, name string, insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage)
}

type stdService struct{}

var (
	stdRepo interfaces.StudentRepository
)

func NewStudentService(repository interfaces.StudentRepository) StudentService {
	stdRepo = repository
	return &stdService{}
}

func (*stdService) CreateStudents(std structs.Student) *structs.ErrorMessage {
	return stdRepo.CreateStudents(std)
}

func (*stdService) CreateParent(prt structs.Parents) *structs.ErrorMessage {
	return stdRepo.CreateParent(prt)
}

func (*stdService) UpdateStudentDetails(std structs.Student) *structs.ErrorMessage {
	return stdRepo.UpdateStudentDetails(std)
}

func (*stdService) GetStudentDetails(stdID string) (*structs.Student, *structs.ErrorMessage) {
	return stdRepo.GetStudentDetails(stdID)
}

func (*stdService) GetStudents(insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage) {
	return stdRepo.GetStudents(insID, page, limit)
}

func (*stdService) GetStudent(nmrInduk string, name string, insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage) {
	return stdRepo.GetStudent(nmrInduk, name, insID, page, limit)
}

func (*stdService) CheckEmail(email string, usrID int) int {
	return stdRepo.CheckEmail(email, usrID)
}
func (*stdService) CheckFamily(famID int, stdID int) int {
	return stdRepo.CheckFamily(famID, stdID)
}

func (*stdService) CheckNomorIndukStd(insID int, nmrInduk string, stdID int) int {
	return stdRepo.CheckNomorIndukStd(insID, nmrInduk, stdID)
}

func (*stdService) Validate(std *structs.Student) (*structs.Student, *structs.ErrorMessage) {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(std)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = std.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
	return std, nil
}

func (*stdService) ValidateParent(prt *structs.Parents) (*structs.Parents, *structs.ErrorMessage) {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(prt)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = prt.Name
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
	return prt, nil
}
