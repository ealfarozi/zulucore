package interfaces

import "github.com/ealfarozi/zulucore/structs"

type StudentRepository interface {
	CreateStudents(std structs.Student) *structs.ErrorMessage
	CreateParent(prt structs.Parents) *structs.ErrorMessage
	CheckEmail(email string, usrID int) int
	CheckNomorIndukStd(insID int, nmrInduk string, stdID int) int
	CheckFamily(famID int, stdID int, parID int) int
	UpdateStudentDetails(std structs.Student) *structs.ErrorMessage
	GetStudentDetails(stdID string) (*structs.Student, *structs.ErrorMessage)
	GetStudents(insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage)
	GetParents(stdID string, page string, limit string) (*[]structs.Parents, *structs.ErrorMessage)
	GetStudent(nmrInduk string, name string, insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage)
}
