package interfaces

import "github.com/ealfarozi/zulucore/structs"

type StudentRepository interface {
	CreateStudents(std structs.Student) *structs.ErrorMessage
	CheckEmail(email string, usrID int) int
	CheckNomorIndukStd(insID int, nmrInduk string, stdID int) int
	UpdateStudentDetails(std structs.Student) *structs.ErrorMessage
	GetStudentDetails(stdID string) (*structs.Student, *structs.ErrorMessage)
	GetStudents(insID string) (*[]structs.Student, *structs.ErrorMessage)
	GetStudent(nmrInduk string, name string, insID string) (*[]structs.Student, *structs.ErrorMessage)
}
