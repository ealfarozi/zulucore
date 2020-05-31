package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
	"gopkg.in/go-playground/validator.v9"

	//blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
)

//GetTutors in the db (all)
func GetTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutors []structs.Tutor

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT ttr.id, ttr.nomor_induk, ttr.name, ttr.tutor_type_id, ttr.user_id, ttr.status FROM tutors ttr inner join (select user_id from user_roles where institution_id = ?) ur on ttr.user_id = ur.user_id"

	//var tutor structs.Tutor

	res, err := db.Query(sqlQuery, r.FormValue("institution_id"))
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Next() != false {
		tutor := structs.Tutor{}
		for res.Next() {
			res.Scan(&tutor.ID, &tutor.NomorInduk, &tutor.Name, &tutor.TutorTypeID, &tutor.UserID, &tutor.Status)
			tutors = append(tutors, tutor)
		}
		json.NewEncoder(w).Encode(tutors)
	} else {
		common.JSONError(w, structs.ErrNotFound, "", http.StatusInternalServerError)
		return
	}

}

//GetTutor in the db based on nomor_induk parameter
func GetTutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prm string

	//params := mux.Vars(r) //get params

	tutors := []structs.Tutor{}
	db := mysql.InitializeMySQL()

	//sqlQuery := "SELECT id, nomor_induk, name, tutor_type_id, user_id, status FROM tutors where "
	sqlQuery := "SELECT ttr.id, ttr.nomor_induk, ttr.name, ttr.tutor_type_id, ttr.user_id, ttr.status FROM tutors ttr inner join (select user_id from user_roles where institution_id = ?) ur on ttr.user_id = ur.user_id where "

	if r.FormValue("nomor_induk") != "" {
		sqlQuery += "ttr.nomor_induk = ?"
		prm = r.FormValue("nomor_induk")
	}
	if r.FormValue("name") != "" {
		sqlQuery += "ttr.name like ?"
		prm = "%" + r.FormValue("name") + "%"
	}

	res, err := db.Query(sqlQuery, r.FormValue("institution_id"), prm)
	if err != nil {
		common.JSONError(w, structs.ErrNotFound, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Next() != false {
		tutor := structs.Tutor{}
		for res.Next() {
			res.Scan(&tutor.ID, &tutor.NomorInduk, &tutor.Name, &tutor.TutorTypeID, &tutor.UserID, &tutor.Status)
			tutors = append(tutors, tutor)
		}
		json.NewEncoder(w).Encode(tutors)
	} else {
		common.JSONError(w, structs.ErrNotFound, "", http.StatusInternalServerError)
		return
	}

}

//CreateTutors is func that will insert multiple tutors at once (complete)
func CreateTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutors []structs.Tutor
	var errstr structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&tutors)

	db := mysql.InitializeMySQL()
	sqlQuery := "insert into tutors (nomor_induk, name, tutor_type_id, user_id) values (?, ?, ?, ?)"
	tx, err := db.Begin()

	j := 0
	for range tutors {
		if err != nil {
			tx.Rollback()
			common.JSONError(w, tutors[j].NomorInduk, err.Error(), http.StatusInternalServerError)
			return
		}

		v := validator.New()
		err := v.Struct(tutors[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := tx.Exec(sqlQuery, &tutors[j].NomorInduk, &tutors[j].Name, &tutors[j].TutorTypeID, &tutors[j].UserID)
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
			return
		}

		lastTutorID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.LastIDErr, err.Error(), http.StatusInternalServerError)
			return
		}

		//insert details
		if tutors[j].Details != nil {
			v := validator.New()
			err := v.Struct(tutors[j].Details)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
				return
			}

			sqlQueryDetail := "insert into tutor_details (education_degree_front, education_degree_back, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, tutor_id, user_id ) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			res, err := tx.Exec(sqlQueryDetail, &tutors[j].Details.EducationFront, &tutors[j].Details.EducationBack, &tutors[j].Details.Ktp, &tutors[j].Details.Sim, &tutors[j].Details.Npwp, &tutors[j].Details.GenderID, &tutors[j].Details.PobID, &tutors[j].Details.Dob, &tutors[j].Details.Phone, &tutors[j].Details.Email, &tutors[j].Details.StreetAddress, &tutors[j].Details.AddressID, &tutors[j].Details.InsSource, &tutors[j].Details.JoinDate, &tutors[j].UserID, &lastTutorID)
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

		//insert educations
		k := 0

		if len(tutors[j].Education) != 0 {
			v := validator.New()
			err := v.Struct(tutors[j].Education)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
				return
			}

			sqlQueryEdu := "insert into tutor_educations (univ_degree_id, univ_name, years, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Education {
				_, err := tx.Exec(sqlQueryEdu, &tutors[j].Education[k].UnivDegreeID, &tutors[j].Education[k].UnivName, &tutors[j].Education[k].Years, &lastTutorID)
				if err != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
					return
				}
				k++
			}
		}

		//insert certificates
		m := 0
		if len(tutors[j].Certificate) != 0 {
			v := validator.New()
			err := v.Struct(tutors[j].Certificate)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
				return
			}

			sqlQueryCert := "insert into tutor_certificates (cert_name, cert_date, tutor_id) values (?, ?, ?)"
			for range tutors[j].Certificate {
				_, err := tx.Exec(sqlQueryCert, &tutors[j].Certificate[m].CertName, &tutors[j].Certificate[m].CertDate, &lastTutorID)
				if err != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
					return
				}
				m++
			}
		}

		//insert Experiences
		n := 0
		if len(tutors[j].Experience) != 0 {
			v := validator.New()
			err := v.Struct(tutors[j].Experience)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
				return
			}

			sqlQueryExp := "insert into tutor_experiences (exp_name, description, years, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Experience {
				_, err := tx.Exec(sqlQueryExp, &tutors[j].Experience[n].ExpName, &tutors[j].Experience[n].Description, &tutors[j].Experience[n].Years, &lastTutorID)
				if err != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
					return
				}
				n++
			}
		}

		//insert Journal
		p := 0
		if len(tutors[j].Journal) != 0 {
			v := validator.New()
			err := v.Struct(tutors[j].Journal)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
				return
			}

			sqlQueryJour := "insert into tutor_journals (journal_name, publish_at, publish_date, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Journal {
				_, err := tx.Exec(sqlQueryJour, &tutors[j].Journal[p].JourName, &tutors[j].Journal[p].PublishAt, &tutors[j].Journal[p].PublishDate, &lastTutorID)
				if err != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
					return
				}
				p++
			}
		}

		//insert research
		a := 0
		if len(tutors[j].Research) != 0 {
			v := validator.New()
			err := v.Struct(tutors[j].Research)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
				return
			}

			sqlQueryRes := "insert into tutor_researches (res_name, description, years, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Research {
				_, err := tx.Exec(sqlQueryRes, &tutors[j].Research[a].ResName, &tutors[j].Research[a].Description, &tutors[j].Research[a].Years, &lastTutorID)
				if err != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
					return
				}
				a++
			}
		}

		errstr.Message = structs.Success
		errstr.Code = http.StatusOK
		j++
	}
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)

}

//UpdateTutors (multiple)
func UpdateTutors(w http.ResponseWriter, r *http.Request) {

}

//DeleteTutors (multiple). this func will soft-delete the data
func DeleteTutors(w http.ResponseWriter, r *http.Request) {

}
