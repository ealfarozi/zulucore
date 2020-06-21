package api

import (
	"encoding/json"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/repositories"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/service"
	"github.com/ealfarozi/zulucore/structs"
	"gopkg.in/go-playground/validator.v9"

	//blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type logic struct{}

var (
	tutorService service.TutorService
	repo         interfaces.TutorRepository = repositories.NewMysqlRepository()
)

type TutorLogic interface {
	GetTutors(w http.ResponseWriter, r *http.Request)
	GetTutorDetails(w http.ResponseWriter, r *http.Request)
	GetTutor(w http.ResponseWriter, r *http.Request)
	UpdateTutorDetails(w http.ResponseWriter, r *http.Request)
	CreateTutors(w http.ResponseWriter, r *http.Request)
}

func NewTutorLogic(service service.TutorService) TutorLogic {
	tutorService = service
	return &logic{}
}

//GetTutors in the db (all)
func (*logic) GetTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tutors, errStr := tutorService.GetTutors(r.FormValue("institution_id"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(tutors)

}

//GetTutorDetails in the db (by tutor ID)
func (*logic) GetTutorDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tutors, errStr := tutorService.GetTutorDetails(r.FormValue("tutor_id"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(tutors)

}

//GetTutor in the db based on nomor_induk parameter (search tutor)
func (*logic) GetTutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tutors, errStr := tutorService.GetTutor(r.FormValue("nomor_induk"), r.FormValue("name"), r.FormValue("institution_id"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(tutors)

}

//UpdateEducations is the func to create/update the education in tutor entity. please note that status = 0 (soft delete)
func UpdateEducations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var edus []structs.TutorEducation
	var errstr structs.ErrorMessage
	var queryStr string
	var refID int

	_ = json.NewDecoder(r.Body).Decode(&edus)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	insertEdu := "insert into tutor_educations (univ_degree_id, univ_name, years, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateEdu := "update tutor_educations set univ_degree_id = ?, univ_name = ?, years = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	j := 0
	for range edus {
		v := validator.New()
		err = v.Struct(edus[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		if edus[j].ID != 0 {
			queryStr = updateEdu
			refID = edus[j].ID
		} else {
			queryStr = insertEdu
			refID = edus[j].TutorID
			edus[j].Status = 1
		}

		_, err2 := tx.Exec(queryStr, &edus[j].UnivDegreeID, &edus[j].UnivName, &edus[j].Years, &edus[j].Status, &refID)
		if err2 != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
			return
		}
		j++
	}

	errstr.Message = structs.Success
	errstr.Code = http.StatusOK
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}

//UpdateCertificates is the func to create/update the certificates in tutor entity. please note that status = 0 (soft delete)
func UpdateCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var certs []structs.TutorCertificate
	var errstr structs.ErrorMessage
	var queryStr string
	var refID int

	_ = json.NewDecoder(r.Body).Decode(&certs)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	insertCert := "insert into tutor_certificates (cert_name, cert_date, status, tutor_id) values (?, ?, ?, ?)"
	updateCert := "update tutor_certificates set cert_name = ?, cert_date = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	j := 0
	for range certs {
		v := validator.New()
		err = v.Struct(certs[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		if certs[j].ID != 0 {
			queryStr = updateCert
			refID = certs[j].ID
		} else {
			queryStr = insertCert
			refID = certs[j].TutorID
			certs[j].Status = 1
		}

		_, err2 := tx.Exec(queryStr, &certs[j].CertName, &certs[j].CertDate, &certs[j].Status, &refID)
		if err2 != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
			return
		}
		j++
	}

	errstr.Message = structs.Success
	errstr.Code = http.StatusOK
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}

//UpdateExperiences is the func to create/update the experiences in tutor entity. please note that status = 0 (soft delete)
func UpdateExperiences(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var exps []structs.TutorExperience
	var errstr structs.ErrorMessage
	var queryStr string
	var refID int

	_ = json.NewDecoder(r.Body).Decode(&exps)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	insertExp := "insert into tutor_experiences (exp_name, description, years, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateExp := "update tutor_experiences set exp_name = ?, description = ?, years = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	j := 0
	for range exps {
		v := validator.New()
		err = v.Struct(exps[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		if exps[j].ID != 0 {
			queryStr = updateExp
			refID = exps[j].ID
		} else {
			queryStr = insertExp
			refID = exps[j].TutorID
			exps[j].Status = 1
		}

		_, err2 := tx.Exec(queryStr, &exps[j].ExpName, &exps[j].Description, &exps[j].Years, &exps[j].Status, &refID)
		if err2 != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
			return
		}
		j++
	}

	errstr.Message = structs.Success
	errstr.Code = http.StatusOK
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}

//UpdateJournals is the func to create/update the journal in tutor entity. please note that status = 0 (soft delete)
func UpdateJournals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jours []structs.TutorJournal
	var errstr structs.ErrorMessage
	var queryStr string
	var refID int

	_ = json.NewDecoder(r.Body).Decode(&jours)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	insertJour := "insert into tutor_journals (journal_name, publish_at, publish_date, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateJour := "update tutor_journals set journal_name = ?, publish_at = ?, publish_date = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	j := 0
	for range jours {
		v := validator.New()
		err = v.Struct(jours[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		if jours[j].ID != 0 {
			queryStr = updateJour
			refID = jours[j].ID
		} else {
			queryStr = insertJour
			refID = jours[j].TutorID
			jours[j].Status = 1
		}

		_, err2 := tx.Exec(queryStr, &jours[j].JourName, &jours[j].PublishAt, &jours[j].PublishDate, &jours[j].Status, &refID)
		if err2 != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
			return
		}
		j++
	}

	errstr.Message = structs.Success
	errstr.Code = http.StatusOK
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}

//UpdateResearches is the func to create/update the journal in tutor entity. please note that status = 0 (soft delete)
func UpdateResearches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rschs []structs.TutorResearch
	var errstr structs.ErrorMessage
	var queryStr string
	var refID int

	_ = json.NewDecoder(r.Body).Decode(&rschs)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	insertRsch := "insert into tutor_researches (res_name, description, years, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateRsch := "update tutor_researches set res_name = ?, description = ?, years = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	j := 0
	for range rschs {
		v := validator.New()
		err = v.Struct(rschs[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		if rschs[j].ID != 0 {
			queryStr = updateRsch
			refID = rschs[j].ID
		} else {
			queryStr = insertRsch
			refID = rschs[j].TutorID
			rschs[j].Status = 1
		}

		_, err2 := tx.Exec(queryStr, &rschs[j].ResName, &rschs[j].Description, &rschs[j].Years, &rschs[j].Status, &refID)
		if err2 != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
			return
		}
		j++
	}

	errstr.Message = structs.Success
	errstr.Code = http.StatusOK
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}

//UpdateTutorDetails is the func to create/update the tutor detail (ONLY) on Frontend side for tutor entity. The update will includes nomor_induk and tutor_name as well.
//Please note that Tutor.status = 0 (soft delete). In order to create a new tutor please refer to CreateTutors func.
//Email field should be coming from Login func.
//Both of ID (tutor and tutor_details) are needed in this API
func (*logic) UpdateTutorDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tutors []structs.Tutor

	_ = json.NewDecoder(r.Body).Decode(&tutors)

	j := 0
	for range tutors {
		tutor, errStr := tutorService.Validate(&tutors[j])
		if errStr != nil {
			common.JSONErr(w, errStr)
			return
		}

		if tutor.ID == 0 {
			errStr := structs.ErrorMessage{Message: structs.EmptyID, SysMessage: "", Code: http.StatusInternalServerError}
			common.JSONErr(w, &errStr)
			return
		}

		checkNomorInduk := tutorService.CheckNomorInduk(tutors[j].InsID, tutors[j].NomorInduk, tutors[j].ID)
		if checkNomorInduk != 0 {
			errStr := structs.ErrorMessage{Message: structs.NomorInd, SysMessage: "", Code: http.StatusInternalServerError}
			common.JSONErr(w, &errStr)
			return
		}

		checkEmail := tutorService.CheckEmail(tutors[j].Details.Email, tutors[j].UserID)
		if checkEmail != 0 {
			errStr := structs.ErrorMessage{Message: structs.Email, SysMessage: "", Code: http.StatusInternalServerError}
			common.JSONErr(w, &errStr)
			return
		}

		//insert or update
		errStr = tutorService.UpdateTutorDetails(*tutor)
		if errStr.Code != http.StatusOK {
			common.JSONErr(w, errStr)
			return
		}
		j++
	}
	errStr := structs.ErrorMessage{Message: structs.Success, SysMessage: "", Code: http.StatusOK}
	common.JSONErr(w, &errStr)
	return
}

//CreateTutors is the func that will insert multiple tutors at once (complete).
//please note that the email in request parameter is the username coming from Login func
func (*logic) CreateTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tutors []structs.Tutor

	_ = json.NewDecoder(r.Body).Decode(&tutors)

	j := 0
	for range tutors {
		tutor, errStr := tutorService.Validate(&tutors[j])
		if errStr != nil {
			common.JSONErr(w, errStr)
			return
		}

		checkNomorInduk := common.CheckNomorInduk(tutors[j].InsID, tutors[j].NomorInduk, 0)
		if checkNomorInduk != 0 {
			errStr := structs.ErrorMessage{Message: structs.NomorInd, SysMessage: "", Code: http.StatusInternalServerError}
			common.JSONErr(w, &errStr)
			return
		}

		checkEmail := tutorService.CheckEmail(tutors[j].Details.Email, tutors[j].UserID)
		if checkEmail != 0 {
			errStr := structs.ErrorMessage{Message: structs.Email, SysMessage: "", Code: http.StatusInternalServerError}
			common.JSONErr(w, &errStr)
			return
		}

		errStr = tutorService.CreateTutors(*tutor)
		if errStr.Code != http.StatusOK {
			common.JSONErr(w, errStr)
			return
		}
		j++
	}

	errStr := structs.ErrorMessage{Message: structs.Success, SysMessage: "", Code: http.StatusOK}
	common.JSONErr(w, &errStr)
	return

}

//CreateTutors is the func that will insert multiple tutors at once (complete).
//please note that the email in request parameter is the username coming from Login func
/*
func CreateTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutors []structs.Tutor
	var errstr structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&tutors)

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	//start checking insert
	if err != nil {
		tx.Rollback()
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	j := 0
	for range tutors {
		check := common.CheckNomorInduk(tutors[j].InsID, tutors[j].NomorInduk, 0)

		if check != 0 {
			tx.Rollback()
			fmt.Printf(tutors[j].NomorInduk)
			common.JSONError(w, structs.NomorInd, "", http.StatusInternalServerError)
			return
		}

		sqlQuery := "insert into tutors (nomor_induk, name, tutor_type_id, user_id) values (?, ?, ?, ?)"
		res, err := tx.Exec(sqlQuery, &tutors[j].NomorInduk, &tutors[j].Name, &tutors[j].TutorTypeID, &tutors[j].UserID)
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
			return
		}

		lastID, err := res.LastInsertId()
		lastTutorID := int(lastID)
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.LastIDErr, err.Error(), http.StatusInternalServerError)
			return
		}

		//insert details
		if tutors[j].Details != nil {
			tutors[j].Details.TutorID = lastTutorID
			tutors[j].Details.UserID = tutors[j].UserID

			sqlQueryDetail := "insert into tutor_details (education_degree_front, education_degree_back, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, tutor_id, user_id ) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
			res, err := tx.Exec(sqlQueryDetail, &tutors[j].Details.EducationFront, &tutors[j].Details.EducationBack, &tutors[j].Details.Ktp, &tutors[j].Details.Sim, &tutors[j].Details.Npwp, &tutors[j].Details.GenderID, &tutors[j].Details.PobID, &tutors[j].Details.Dob, &tutors[j].Details.Phone, &tutors[j].Details.Email, &tutors[j].Details.StreetAddress, &tutors[j].Details.AddressID, &tutors[j].Details.InsSource, &tutors[j].Details.JoinDate, &lastTutorID, &tutors[j].UserID)
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
			sqlQueryEdu := "insert into tutor_educations (univ_degree_id, univ_name, years, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Education {
				tutors[j].Education[k].TutorID = lastTutorID

				_, err2 := tx.Exec(sqlQueryEdu, &tutors[j].Education[k].UnivDegreeID, &tutors[j].Education[k].UnivName, &tutors[j].Education[k].Years, &lastTutorID)
				if err2 != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
					return
				}
				k++
			}
		}

		//insert certificates
		m := 0
		if len(tutors[j].Certificate) != 0 {
			sqlQueryCert := "insert into tutor_certificates (cert_name, cert_date, tutor_id) values (?, ?, ?)"
			for range tutors[j].Certificate {
				tutors[j].Certificate[m].TutorID = lastTutorID

				_, err2 := tx.Exec(sqlQueryCert, &tutors[j].Certificate[m].CertName, &tutors[j].Certificate[m].CertDate, &lastTutorID)
				if err2 != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
					return
				}
				m++
			}
		}

		//insert Experiences
		n := 0
		if len(tutors[j].Experience) != 0 {
			sqlQueryExp := "insert into tutor_experiences (exp_name, description, years, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Experience {
				tutors[j].Experience[n].TutorID = lastTutorID

				_, err2 := tx.Exec(sqlQueryExp, &tutors[j].Experience[n].ExpName, &tutors[j].Experience[n].Description, &tutors[j].Experience[n].Years, &lastTutorID)
				if err2 != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
					return
				}
				n++
			}
		}

		//insert Journal
		p := 0
		if len(tutors[j].Journal) != 0 {
			sqlQueryJour := "insert into tutor_journals (journal_name, publish_at, publish_date, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Journal {
				tutors[j].Journal[p].TutorID = lastTutorID

				_, err2 := tx.Exec(sqlQueryJour, &tutors[j].Journal[p].JourName, &tutors[j].Journal[p].PublishAt, &tutors[j].Journal[p].PublishDate, &lastTutorID)
				if err2 != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
					return
				}
				p++
			}
		}

		//insert research
		a := 0
		if len(tutors[j].Research) != 0 {
			sqlQueryRes := "insert into tutor_researches (res_name, description, years, tutor_id) values (?, ?, ?, ?)"
			for range tutors[j].Research {
				tutors[j].Research[a].TutorID = lastTutorID

				_, err2 := tx.Exec(sqlQueryRes, &tutors[j].Research[a].ResName, &tutors[j].Research[a].Description, &tutors[j].Research[a].Years, &lastTutorID)
				if err2 != nil {
					tx.Rollback()
					common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
					return
				}
				a++
			}
		}

		v := validator.New()
		err = v.Struct(tutors[j])
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
*/
