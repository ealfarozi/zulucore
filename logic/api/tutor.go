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
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&tutors)

	j := 0
	for range tutors {
		tutor, errStr := tutorService.Validate(&tutors[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		if tutor.ID == 0 {
			errStr := structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.EmptyID, SysMessage: "", Code: http.StatusInternalServerError}
			errs = append(errs, errStr)
		}

		checkNomorInduk := tutorService.CheckNomorInduk(tutors[j].InsID, tutors[j].NomorInduk, tutors[j].ID)
		if checkNomorInduk != 0 {
			errStr := structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.NomorInd, SysMessage: "", Code: http.StatusInternalServerError}
			errs = append(errs, errStr)
		}

		checkEmail := tutorService.CheckEmail(tutors[j].Details.Email, tutors[j].UserID)
		if checkEmail != 0 {
			errStr := structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.Email, SysMessage: "", Code: http.StatusInternalServerError}
			errs = append(errs, errStr)
		}

		//insert or update
		errStr = tutorService.UpdateTutorDetails(*tutor)
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		}
		errs = append(errs, structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		j++
	}

	common.JSONErrs(w, &errs)
	return
}

//CreateTutors is the func that will insert multiple tutors at once (complete).
//please note that the email in request parameter is the username coming from Login func
func (*logic) CreateTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tutors []structs.Tutor
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&tutors)

	j := 0
	for range tutors {
		tutor, errStr := tutorService.Validate(&tutors[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		checkNomorInduk := common.CheckNomorInduk(tutors[j].InsID, tutors[j].NomorInduk, 0)
		if checkNomorInduk != 0 {
			errStr := structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.NomorInd, SysMessage: "", Code: http.StatusInternalServerError}
			errs = append(errs, errStr)
		}

		if tutors[j].Details != nil {
			checkEmail := tutorService.CheckEmail(tutors[j].Details.Email, tutors[j].UserID)
			if checkEmail != 0 {
				errStr := structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.Email, SysMessage: "", Code: http.StatusInternalServerError}
				errs = append(errs, errStr)
			}
		}

		errStr = tutorService.CreateTutors(*tutor)
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		}

		errs = append(errs, structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		j++

	}
	common.JSONErrs(w, &errs)
	return

}
