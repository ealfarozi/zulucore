package api

import (
	"encoding/json"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/repositories"
	"github.com/ealfarozi/zulucore/service"
	"github.com/ealfarozi/zulucore/structs"
)

type stdLogic struct{}

var (
	stdService service.StudentService
	stdRepo    interfaces.StudentRepository = repositories.NewStudentRepository()
)

//TutorLogic is the interface of httpRequest for Tutor
type StudentLogic interface {
	CreateStudents(w http.ResponseWriter, r *http.Request)
	UpdateStudentDetails(w http.ResponseWriter, r *http.Request)
	GetStudentDetails(w http.ResponseWriter, r *http.Request)
	GetStudents(w http.ResponseWriter, r *http.Request)
	GetStudent(w http.ResponseWriter, r *http.Request)
}

//NewTutorLogic is the func to calling the constructor of tutor interface
func NewStudentLogic(service service.StudentService) StudentLogic {
	stdService = service
	return &stdLogic{}
}

//CreateStudents is the func to insert the student data
func (*stdLogic) CreateStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stds []structs.Student
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&stds)

	for j := range stds {
		std, errStr := stdService.Validate(&stds[j])
		if errStr != nil {
			errs = append(errs, *errStr)
			continue
		}

		checkNomorInduk := stdService.CheckNomorIndukStd(stds[j].InsID, stds[j].NomorInduk, 0)
		if checkNomorInduk != 0 {
			errStr := structs.ErrorMessage{Data: stds[j].NomorInduk, Message: structs.NomorInd, SysMessage: "", Code: http.StatusInternalServerError}
			errs = append(errs, errStr)
			continue
		}

		if stds[j].Details != nil {
			checkEmail := stdService.CheckEmail(stds[j].Details.Email, stds[j].UserID)
			if checkEmail != 0 {
				errStr := structs.ErrorMessage{Data: stds[j].NomorInduk, Message: structs.Email, SysMessage: "", Code: http.StatusInternalServerError}
				errs = append(errs, errStr)
				continue
			}
		}

		errStr = stdService.CreateStudents(*std)
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
			continue
		} else {
			errs = append(errs, structs.ErrorMessage{Data: stds[j].NomorInduk, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
	}
	common.JSONErrs(w, &errs)
	return

}

//GetStudent is the func to get the student list based on nomor induk and name
func (*stdLogic) GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	std, errStr := stdService.GetStudent(r.FormValue("nomor_induk"), r.FormValue("name"), r.FormValue("institution_id"), r.FormValue("_page"), r.FormValue("_limit"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(std)
}

//GetStudents in the db (all)
func (*stdLogic) GetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	std, errStr := stdService.GetStudents(r.FormValue("institution_id"), r.FormValue("_page"), r.FormValue("_limit"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(std)
}

//GetStudentDetails is the func to get the student details based on student_id
func (*stdLogic) GetStudentDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	std, errStr := stdService.GetStudentDetails(r.FormValue("student_id"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(std)
}

//UpdateDetails is the func to create/update the student detail (ONLY) on Frontend side for student entity. the update will includes nomor_induk and student_name as well
//please note that status = 0 = soft delete.
//In order to create a new student please refer to CreateStudents func
//email field should be coming from Login func
func (*stdLogic) UpdateStudentDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stds []structs.Student
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&stds)

	for j := range stds {
		std, errStr := stdService.Validate(&stds[j])
		if errStr != nil {
			errs = append(errs, *errStr)
			continue
		}

		checkNomorInduk := stdService.CheckNomorIndukStd(stds[j].InsID, stds[j].NomorInduk, stds[j].ID)
		if checkNomorInduk != 0 {
			errStr := structs.ErrorMessage{Data: stds[j].NomorInduk, Message: structs.NomorInd, SysMessage: "", Code: http.StatusInternalServerError}
			errs = append(errs, errStr)
			continue
		}

		if stds[j].Details != nil {
			checkEmail := stdService.CheckEmail(stds[j].Details.Email, stds[j].UserID)
			if checkEmail != 0 {
				errStr := structs.ErrorMessage{Data: stds[j].NomorInduk, Message: structs.Email, SysMessage: "", Code: http.StatusInternalServerError}
				errs = append(errs, errStr)
				continue
			}
		}

		errStr = stdService.UpdateStudentDetails(*std)
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
			continue
		} else {
			errs = append(errs, structs.ErrorMessage{Data: stds[j].NomorInduk, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
	}
	common.JSONErrs(w, &errs)
	return
}

//UpdateParents is the func to update the student's parents information
func UpdateParents(w http.ResponseWriter, r *http.Request) {

}

//GetScores is the func to get all of the scores of a student
func GetScores(w http.ResponseWriter, r *http.Request) {

}

//UpdateScores is the func to update the scores
func UpdateScores(w http.ResponseWriter, r *http.Request) {

}

//GetSemesters is the func to get all of the semesters status based on institution id
func GetSemesters(w http.ResponseWriter, r *http.Request) {

}
