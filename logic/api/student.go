package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
	"gopkg.in/go-playground/validator.v9"
)

//CreateStudents is the func to insert the student data
func CreateStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var students []structs.Student
	var errstr structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&students)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	//start checking insert
	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	j := 0
	for range students {
		check := common.CheckNomorIndukStd(students[j].InsID, students[j].NomorInduk, 0)

		if check != 0 {
			tx.Rollback()
			fmt.Printf(students[j].NomorInduk)
			common.JSONError(w, structs.NomorInd, "", http.StatusInternalServerError)
			return
		}

		sqlQuery := "insert into students (nomor_induk, name, degree_id, student_type_id, curr_id) values (?, ?, ?, ?, ?)"
		res, err := tx.Exec(sqlQuery, &students[j].NomorInduk, &students[j].Name, &students[j].DegreeID, &students[j].StudentType, &students[j].CurrID)
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
			return
		}

		lastID, err := res.LastInsertId()
		lastStudentID := int(lastID)
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.LastIDErr, err.Error(), http.StatusInternalServerError)
			return
		}

		//insert details
		if students[j].Details != nil {
			students[j].Details.StudentID = lastStudentID
			students[j].Details.UserID = students[j].UserID

			sqlQueryDetail := "insert into student_details (kk_no, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, student_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			res, err := tx.Exec(sqlQueryDetail, &students[j].Details.KkNO, &students[j].Details.Ktp, &students[j].Details.Sim, &students[j].Details.Npwp, &students[j].Details.GenderID, &students[j].Details.PobID, &students[j].Details.Dob, &students[j].Details.Phone, &students[j].Details.Email, &students[j].Details.StreetAddress, &students[j].Details.AddressID, &students[j].Details.InsSource, &students[j].Details.JoinDate, &lastStudentID)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
				return
			}

			lastID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.LastIDErr, err.Error(), http.StatusInternalServerError)
				log.Println(lastID)
				return
			}
		}

		v := validator.New()
		err = v.Struct(students[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		errstr.Message = structs.Success
		errstr.Code = http.StatusOK
		j++
	}
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}

//GetStudent is the func to get the student list based on nomor induk and name
func GetStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prm string

	students := []structs.Student{}
	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT std.id, std.nomor_induk, std.name, std.degree_id, std.student_type_id, std.curr_id, std.user_id, std.status FROM students std inner join (select user_id from user_roles where institution_id = ?) ur on std.user_id = ur.user_id where "

	if r.FormValue("nomor_induk") != "" {
		sqlQuery += "std.nomor_induk like ?"
		prm = "%" + r.FormValue("nomor_induk") + "%"
	}
	if r.FormValue("name") != "" {
		sqlQuery += "std.name like ?"
		prm = "%" + r.FormValue("name") + "%"
	}
	res, err := db.Query(sqlQuery, r.FormValue("institution_id"), prm)
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.ErrNotFound, err.Error(), http.StatusInternalServerError)
		return
	}

	student := structs.Student{}
	for res.Next() {
		res.Scan(&student.ID, &student.NomorInduk, &student.Name, &student.DegreeID, &student.StudentType, &student.CurrID, &student.UserID, &student.Status)
		students = append(students, student)
	}

	if len(students) != 0 {
		json.NewEncoder(w).Encode(students)
	} else {
		common.JSONError(w, structs.ErrNotFound, "", http.StatusInternalServerError)
		return
	}
}

//GetStudents in the db (all)
func GetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var students []structs.Student

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT std.id, std.nomor_induk, std.name, std.degree_id, std.student_type_id, std.curr_id, std.user_id, std.status FROM students std inner join (select user_id from user_roles where institution_id = ?) ur on std.user_id = ur.user_id "

	res, err := db.Query(sqlQuery, r.FormValue("institution_id"))
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	student := structs.Student{}
	for res.Next() {
		res.Scan(&student.ID, &student.NomorInduk, &student.Name, &student.DegreeID, &student.StudentType, &student.CurrID, &student.UserID, &student.Status)
		students = append(students, student)
	}

	if len(students) != 0 {
		json.NewEncoder(w).Encode(students)
	} else {
		common.JSONError(w, structs.ErrNotFound, "", http.StatusInternalServerError)
		return
	}
}

//GetStudentDetails is the func to get the student details based on student_id
func GetStudentDetails(w http.ResponseWriter, r *http.Request) {

}

//UpdateDetails is the func to update all of the student's information
func UpdateDetais(w http.ResponseWriter, r *http.Request) {

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
