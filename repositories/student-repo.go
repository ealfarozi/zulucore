package repositories

import (
	"log"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
)

type stdRepo struct{}

//NewStudentRepository is the constructor for student-repo
func NewStudentRepository() interfaces.StudentRepository {
	return &stdRepo{}
}

func (*stdRepo) CreateStudents(std structs.Student) *structs.ErrorMessage {
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	//start checking insert
	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = std.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	sqlQuery := "insert into students (nomor_induk, name, degree_id, student_type_id, curr_id, user_id) values (?, ?, ?, ?, ?, ?)"
	res, err := tx.Exec(sqlQuery, &std.NomorInduk, &std.Name, &std.DegreeID, &std.StudentType, &std.CurrID, &std.UserID)
	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = std.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	lastID, err := res.LastInsertId()
	lastStudentID := int(lastID)
	if err != nil {
		tx.Rollback()
		errors.Message = structs.LastIDErr
		errors.Data = std.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	//insert details
	if std.Details != nil {
		std.Details.StudentID = lastStudentID
		std.Details.UserID = std.UserID

		sqlQueryDetail := "insert into student_details (kk_no, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, student_id, user_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		res, err := tx.Exec(sqlQueryDetail, &std.Details.KkNO, &std.Details.Ktp, &std.Details.Sim, &std.Details.Npwp, &std.Details.GenderID, &std.Details.PobID, &std.Details.Dob, &std.Details.Phone, &std.Details.Email, &std.Details.StreetAddress, &std.Details.AddressID, &std.Details.InsSource, &std.Details.JoinDate, &lastStudentID, &std.Details.UserID)
		if err != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = std.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			errors.Message = structs.LastIDErr
			errors.Data = std.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			log.Println(lastID)
			return &errors
		}
	}

	errors.Message = structs.Success
	errors.Code = http.StatusOK

	tx.Commit()
	return &errors
}

func (*stdRepo) CheckEmail(email string, usrID int) int {
	db := mysql.InitializeMySQL()
	sqlQueryCheck := "select count(1) from users where username = ? and id != ?"
	check := 0
	err := db.QueryRow(sqlQueryCheck, &email, &usrID).Scan(&check)

	if err != nil {
		check = 99
	}
	return check
}

func (*stdRepo) CheckNomorIndukStd(insID int, nmrInduk string, stdID int) int {
	db := mysql.InitializeMySQL()
	sqlQueryCheck := "SELECT count(1) FROM students std inner join (select user_id from user_roles where institution_id = ?) ur on std.user_id = ur.user_id where std.nomor_induk = ? "
	check := 0
	if stdID != 0 {
		sqlQueryCheck += "and std.id != ?"
		err := db.QueryRow(sqlQueryCheck, &insID, &nmrInduk, &stdID).Scan(&check)
		if err != nil {
			check = 99
		}
	} else {
		err := db.QueryRow(sqlQueryCheck, &insID, &nmrInduk).Scan(&check)
		if err != nil {
			check = 99
		}
	}

	return check
}

//UpdateStudentDetails is the func to create/update the student detail (ONLY) on Frontend side for student entity. the update will includes nomor_induk and student_name as well
//please note that status = 0 = soft delete.
//In order to create a new student please refer to CreateStudents func
//email field should be coming from Login func
func (*stdRepo) UpdateStudentDetails(std structs.Student) *structs.ErrorMessage {
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	//start checking insert
	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = std.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	insertDet := "insert into student_details (kk_no, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, student_id, tutor_id, user_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	updateDet := "update student_details set kk_no = ?, ktp = ?, sim = ?, npwp = ?, gender_id = ?, pob_id = ?, dob = ?, phone = ?, email = ?, street_address = ?, address_id = ?, institution_source_name = ?, join_date = ?, tutor_id = ?, updated_at = now(), updated_by = 'API' where id = ?"
	updateStd := "update students set nomor_induk = ?, name = ?, degree_id = ?, student_type_id = ?, curr_id = ?, status = ? where id = ?"
	updateUsr := "update users set username = ? where id = ?"

	std.Details.StudentID = std.ID

	if std.Details.ID != 0 {
		//update
		_, err := tx.Exec(updateDet, &std.Details.KkNO, &std.Details.Ktp, &std.Details.Sim, &std.Details.Npwp, &std.Details.GenderID, &std.Details.PobID, &std.Details.Dob, &std.Details.Phone, &std.Details.Email, &std.Details.StreetAddress, &std.Details.AddressID, &std.Details.InsSource, &std.Details.JoinDate, &std.Details.TutorID, &std.Details.ID)
		if err != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = std.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}

		_, err2 := tx.Exec(updateStd, &std.NomorInduk, &std.Name, &std.DegreeID, &std.StudentType, &std.CurrID, &std.Status, &std.ID)
		if err2 != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = std.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}

		_, err3 := tx.Exec(updateUsr, &std.Details.Email, &std.UserID)
		if err3 != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = std.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}
	} else {
		//insert
		_, err := tx.Exec(insertDet, &std.Details.KkNO, &std.Details.Ktp, &std.Details.Sim, &std.Details.Npwp, &std.Details.GenderID, &std.Details.PobID, &std.Details.Dob, &std.Details.Phone, &std.Details.Email, &std.Details.StreetAddress, &std.Details.AddressID, &std.Details.InsSource, &std.Details.JoinDate, &std.ID, &std.Details.TutorID, &std.UserID)
		if err != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = std.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}
	}

	errors.Message = structs.Success
	errors.Code = http.StatusOK

	tx.Commit()
	return &errors
}

//GetStudent is the func to get the student list based on nomor induk and name
func (*stdRepo) GetStudent(nmrInduk string, name string, insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage) {
	var students []structs.Student
	var errors structs.ErrorMessage
	var prm string

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT std.id, std.nomor_induk, std.name, std.degree_id, std.student_type_id, std.curr_id, std.user_id, std.status FROM students std inner join (select user_id from user_roles where institution_id = ?) ur on std.user_id = ur.user_id where "

	if nmrInduk != "" {
		sqlQuery += "std.nomor_induk like ?"
		prm = "%" + nmrInduk + "%"
	}
	if name != "" {
		sqlQuery += "std.name like ?"
		prm = "%" + name + "%"
	}

	sqlQuery += common.SetPageLimit(page, limit)
	res, err := db.Query(sqlQuery, insID, prm)
	defer mysql.CloseRows(res)
	if err != nil {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

	student := structs.Student{}
	for res.Next() {
		res.Scan(&student.ID, &student.NomorInduk, &student.Name, &student.DegreeID, &student.StudentType, &student.CurrID, &student.UserID, &student.Status)
		students = append(students, student)
	}

	if len(students) != 0 {
		return &students, nil
	} else {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = ""
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
}

//GetStudents in the db (all)
func (*stdRepo) GetStudents(insID string, page string, limit string) (*[]structs.Student, *structs.ErrorMessage) {
	var students []structs.Student
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT std.id, std.nomor_induk, std.name, std.degree_id, std.student_type_id, std.curr_id, std.user_id, std.status FROM students std inner join (select user_id from user_roles where institution_id = ?) ur on std.user_id = ur.user_id "
	sqlQuery += common.SetPageLimit(page, limit)
	res, err := db.Query(sqlQuery, insID)
	defer mysql.CloseRows(res)
	if err != nil {
		errors.Message = structs.QueryErr
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

	student := structs.Student{}
	for res.Next() {
		res.Scan(&student.ID, &student.NomorInduk, &student.Name, &student.DegreeID, &student.StudentType, &student.CurrID, &student.UserID, &student.Status)
		students = append(students, student)
	}

	if len(students) != 0 {
		return &students, nil
	} else {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = ""
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
}

//GetStudentDetails is the func to get the student details based on student_id
func (*stdRepo) GetStudentDetails(stdID string) (*structs.Student, *structs.ErrorMessage) {
	var student structs.Student
	var det structs.StudentDetails
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()
	sqlQueryStudent := "select id, nomor_induk, name, degree_id, student_type_id, curr_id, user_id, status from students where id = ?"
	err := db.QueryRow(sqlQueryStudent, stdID).Scan(&student.ID, &student.NomorInduk, &student.Name, &student.DegreeID, &student.StudentType, &student.CurrID, &student.UserID, &student.Status)
	if err != nil {
		errors.Message = structs.QueryErr
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

	//Details
	sqlQueryDetail := "select id, kk_no, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, tutor_id from student_details where student_id = ?"
	resDet, err := db.Query(sqlQueryDetail, stdID)
	defer mysql.CloseRows(resDet)
	if err != nil {
		errors.Message = structs.QueryErr
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

	for resDet.Next() {
		resDet.Scan(&det.ID, &det.KkNO, &det.Ktp, &det.Sim, &det.Npwp, &det.GenderID, &det.PobID, &det.Dob, &det.Phone, &det.Email, &det.StreetAddress, &det.AddressID, &det.InsSource, &det.JoinDate, &det.TutorID)
		student.Details = &det
		student.Details.AddressDetail = common.GetAddressOnly(det.AddressID)
	}

	if student.ID != 0 {
		return &student, nil
	} else {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = ""
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
}

/*
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
*/
