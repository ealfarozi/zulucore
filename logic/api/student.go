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
	w.Header().Set("Content-Type", "application/json")

	var student structs.Student
	var det structs.StudentDetails

	db := mysql.InitializeMySQL()
	sqlQueryStudent := "select id, nomor_induk, name, degree_id, student_type_id, curr_id, user_id, status from students where id = ?"
	err := db.QueryRow(sqlQueryStudent, r.FormValue("student_id")).Scan(&student.ID, &student.NomorInduk, &student.Name, &student.DegreeID, &student.StudentType, &student.CurrID, &student.UserID, &student.Status)
	if err != nil {
		common.JSONError(w, structs.ErrNotFound, err.Error(), http.StatusInternalServerError)
		return
	}

	//Details
	sqlQueryDetail := "select id, kk_no, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, tutor_id from student_details where student_id = ?"
	resDet, err := db.Query(sqlQueryDetail, r.FormValue("student_id"))
	defer mysql.CloseRows(resDet)
	for resDet.Next() {
		resDet.Scan(&det.ID, &det.KkNO, &det.Ktp, &det.Sim, &det.Npwp, &det.GenderID, &det.PobID, &det.Dob, &det.Phone, &det.Email, &det.StreetAddress, &det.AddressID, &det.InsSource, &det.JoinDate, &det.TutorID)
		student.Details = &det
		student.Details.AddressDetail = common.GetAddressOnly(det.AddressID)
	}

	json.NewEncoder(w).Encode(student)
}

//UpdateDetails is the func to create/update the student detail (ONLY) on Frontend side for student entity. the update will includes nomor_induk and student_name as well
//please note that status = 0 = soft delete.
//In order to create a new student please refer to CreateStudents func
//email field should be coming from Login func
func UpdateStudentDetais(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var students []structs.Student
	var errstr structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&students)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	//updating email will update the username in users table
	insertDet := "insert into student_details (kk_no, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, student_id, tutor_id, user_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	updateDet := "update student_details set kk_no = ?, ktp = ?, sim = ?, npwp = ?, gender_id = ?, pob_id = ?, dob = ?, phone = ?, email = ?, street_address = ?, address_id = ?, institution_source_name = ?, join_date = ?, tutor_id = ?, updated_at = now(), updated_by = 'API' where id = ?"
	updateStd := "update students set nomor_induk = ?, name = ?, degree_id = ?, student_type_id = ?, curr_id = ?, status = ? where id = ?"
	updateUsr := "update users set username = ? where id = ?"

	j := 0
	for range students {
		students[j].Details.StudentID = students[j].ID
		v := validator.New()
		err = v.Struct(students[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		checkNmrInduk := common.CheckNomorIndukStd(students[j].InsID, students[j].NomorInduk, students[j].ID)

		if checkNmrInduk != 0 {
			tx.Rollback()
			common.JSONError(w, structs.NomorInd, "", http.StatusInternalServerError)
			return
		}

		checkEmail := common.CheckEmail(students[j].Details.Email, students[j].UserID)

		if checkEmail != 0 {
			tx.Rollback()
			common.JSONError(w, structs.Email, "", http.StatusInternalServerError)
			return
		}

		if students[j].ID != 0 {
			//update
			_, err := tx.Exec(updateDet, &students[j].Details.KkNO, &students[j].Details.Ktp, &students[j].Details.Sim, &students[j].Details.Npwp, &students[j].Details.GenderID, &students[j].Details.PobID, &students[j].Details.Dob, &students[j].Details.Phone, &students[j].Details.Email, &students[j].Details.StreetAddress, &students[j].Details.AddressID, &students[j].Details.InsSource, &students[j].Details.JoinDate, &students[j].Details.TutorID, &students[j].Details.ID)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err2 := tx.Exec(updateStd, &students[j].NomorInduk, &students[j].Name, &students[j].DegreeID, &students[j].StudentType, &students[j].CurrID, &students[j].Status, &students[j].ID)
			if err2 != nil {
				tx.Rollback()
				common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err3 := tx.Exec(updateUsr, &students[j].Details.Email, &students[j].UserID)
			if err3 != nil {
				tx.Rollback()
				common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			//insert
			_, err := tx.Exec(insertDet, &students[j].Details.KkNO, &students[j].Details.Ktp, &students[j].Details.Sim, &students[j].Details.Npwp, &students[j].Details.GenderID, &students[j].Details.PobID, &students[j].Details.Dob, &students[j].Details.Phone, &students[j].Details.Email, &students[j].Details.StreetAddress, &students[j].Details.AddressID, &students[j].Details.InsSource, &students[j].Details.JoinDate, &students[j].ID, &students[j].Details.TutorID, &students[j].UserID)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		j++
	}

	errstr.Message = structs.Success
	errstr.Code = http.StatusOK
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
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
