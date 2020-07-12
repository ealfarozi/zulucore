package repositories

import (
	"log"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
)

type repo struct{}

//NewTutorRepository is the constructor for tutor-repo
func NewTutorRepository() interfaces.TutorRepository {
	return &repo{}
}

func (*repo) GetTutorDetails(tutorID string) (*structs.Tutor, *structs.ErrorMessage) {

	var tutor structs.Tutor
	var errors structs.ErrorMessage
	var det structs.TutorDetails
	var edus []structs.TutorEducation
	var certs []structs.TutorCertificate
	var exps []structs.TutorExperience
	var jours []structs.TutorJournal
	var rschs []structs.TutorResearch

	db := mysql.InitializeMySQL()
	sqlQueryTutor := "select id, nomor_induk, name, tutor_type_id, user_id, status from tutors where id = ?"
	err := db.QueryRow(sqlQueryTutor, tutorID).Scan(&tutor.ID, &tutor.NomorInduk, &tutor.Name, &tutor.TutorTypeID, &tutor.UserID, &tutor.Status)
	if err != nil {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

	//Details
	sqlQueryDetail := "select id, education_degree_front, education_degree_back, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date from tutor_details where tutor_id = ?"
	resDet, err := db.Query(sqlQueryDetail, tutorID)
	defer mysql.CloseRows(resDet)
	for resDet.Next() {
		resDet.Scan(&det.ID, &det.EducationFront, &det.EducationBack, &det.Ktp, &det.Sim, &det.Npwp, &det.GenderID, &det.PobID, &det.Dob, &det.Phone, &det.Email, &det.StreetAddress, &det.AddressID, &det.InsSource, &det.JoinDate)
		tutor.Details = &det
		tutor.Details.AddressDetail = common.GetAddressOnly(det.AddressID)
	}

	//Educations
	sqlQueryEdu := "select id, univ_degree_id, univ_name, years from tutor_educations where status = 1 and tutor_id = ?"
	res, err := db.Query(sqlQueryEdu, tutorID)
	defer mysql.CloseRows(res)

	edu := structs.TutorEducation{}
	for res.Next() {
		res.Scan(&edu.ID, &edu.UnivDegreeID, &edu.UnivName, &edu.Years)
		edus = append(edus, edu)
	}
	tutor.Education = edus

	//Certificates
	sqlQueryCert := "select id, cert_name, cert_date from tutor_certificates where status = 1 and tutor_id = ?"
	resCert, err := db.Query(sqlQueryCert, tutorID)
	defer mysql.CloseRows(resCert)

	cert := structs.TutorCertificate{}
	for resCert.Next() {
		resCert.Scan(&cert.ID, &cert.CertName, &cert.CertDate)
		certs = append(certs, cert)
	}
	tutor.Certificate = certs

	//Experiences
	sqlQueryExp := "select id, exp_name, description, years from tutor_experiences where status = 1 and tutor_id = ?"
	resExp, err := db.Query(sqlQueryExp, tutorID)
	defer mysql.CloseRows(resExp)

	exp := structs.TutorExperience{}
	for resExp.Next() {
		resExp.Scan(&exp.ID, &exp.ExpName, &exp.Description, &exp.Years)
		exps = append(exps, exp)
	}
	tutor.Experience = exps

	//Journals
	sqlQueryJour := "select id, journal_name, publish_at, publish_date from tutor_journals where status = 1 and tutor_id = ?"
	resJour, err := db.Query(sqlQueryJour, tutorID)
	defer mysql.CloseRows(resJour)

	jour := structs.TutorJournal{}

	for resJour.Next() {
		resJour.Scan(&jour.ID, &jour.JourName, &jour.PublishAt, &jour.PublishDate)
		jours = append(jours, jour)
	}
	tutor.Journal = jours

	//Researches
	sqlQueryRes := "select id, res_name, description, years from tutor_researches where status = 1 and tutor_id = ?"
	resRes, err := db.Query(sqlQueryRes, tutorID)
	defer mysql.CloseRows(resRes)

	rsch := structs.TutorResearch{}

	for resRes.Next() {
		resRes.Scan(&rsch.ID, &rsch.ResName, &rsch.Description, &rsch.Years)
		rschs = append(rschs, rsch)
	}
	tutor.Research = rschs

	return &tutor, nil
}

func (*repo) GetTutors(insID string, page string, limit string) (*[]structs.Tutor, *structs.ErrorMessage) {
	var tutors []structs.Tutor
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT ttr.id, ttr.nomor_induk, ttr.name, ttr.tutor_type_id, ttr.user_id, ttr.status FROM tutors ttr inner join (select user_id from user_roles where institution_id = ?) ur on ttr.user_id = ur.user_id"
	sqlQuery += common.SetPageLimit(page, limit)

	res, err := db.Query(sqlQuery, insID)
	defer mysql.CloseRows(res)
	if err != nil {
		errors.Message = structs.QueryErr
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

	tutor := structs.Tutor{}
	for res.Next() {
		res.Scan(&tutor.ID, &tutor.NomorInduk, &tutor.Name, &tutor.TutorTypeID, &tutor.UserID, &tutor.Status)
		tutors = append(tutors, tutor)
	}

	if len(tutors) != 0 {
		return &tutors, nil
	} else {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = ""
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
}

func (*repo) GetTutor(nmrInd string, name string, insID string, page string, limit string) (*[]structs.Tutor, *structs.ErrorMessage) {
	var prm string

	var tutors []structs.Tutor
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT ttr.id, ttr.nomor_induk, ttr.name, ttr.tutor_type_id, ttr.user_id, ttr.status FROM tutors ttr inner join (select user_id from user_roles where institution_id = ?) ur on ttr.user_id = ur.user_id where "

	if nmrInd != "" {
		sqlQuery += "ttr.nomor_induk like ?"
		prm = "%" + nmrInd + "%"
	}
	if name != "" {
		sqlQuery += "ttr.name like ?"
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

	tutor := structs.Tutor{}
	for res.Next() {
		res.Scan(&tutor.ID, &tutor.NomorInduk, &tutor.Name, &tutor.TutorTypeID, &tutor.UserID, &tutor.Status)
		tutors = append(tutors, tutor)
	}

	if len(tutors) != 0 {
		return &tutors, nil
	} else {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = ""
		errors.Code = http.StatusInternalServerError
		return nil, &errors

	}
}

//UpdateTutorDetails is the func to create/update the tutor detail (ONLY) on Frontend side for tutor entity. The update will includes nomor_induk and tutor_name as well.
//Please note that Tutor.status = 0 (soft delete). In order to create a new tutor please refer to CreateTutors func.
//Email field should be coming from Login func.
//Both of ID (tutor and tutor_details) are needed in this API
func (*repo) UpdateTutorDetails(tutor structs.Tutor) *structs.ErrorMessage {
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = tutor.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	//updating email will update the username in users table
	insertDet := "insert into tutor_details (education_degree_front, education_degree_back, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, tutor_id, user_id) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	updateDet := "update tutor_details set education_degree_front = ?, education_degree_back = ?, ktp = ?, sim = ?, npwp = ?, gender_id = ?, pob_id = ?, dob = ?, phone = ?, email = ?, street_address = ?, address_id = ?, institution_source_name = ?, join_date = ?, updated_at = now(), updated_by = 'API' where id = ?"
	updateTut := "update tutors set nomor_induk = ?, name = ?, tutor_type_id = ?, status = ? where id = ?"
	updateUsr := "update users set username = ? where id = ?"

	if tutor.Details.ID != 0 {
		//update
		_, err := tx.Exec(updateDet, &tutor.Details.EducationFront, &tutor.Details.EducationBack, &tutor.Details.Ktp, &tutor.Details.Sim, &tutor.Details.Npwp, &tutor.Details.GenderID, &tutor.Details.PobID, &tutor.Details.Dob, &tutor.Details.Phone, &tutor.Details.Email, &tutor.Details.StreetAddress, &tutor.Details.AddressID, &tutor.Details.InsSource, &tutor.Details.JoinDate, &tutor.Details.ID)
		if err != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = tutor.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}

		_, err2 := tx.Exec(updateTut, &tutor.NomorInduk, &tutor.Name, &tutor.TutorTypeID, &tutor.Status, &tutor.ID)
		if err2 != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = tutor.NomorInduk
			errors.SysMessage = err2.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}

		_, err3 := tx.Exec(updateUsr, &tutor.Details.Email, &tutor.UserID)
		if err3 != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = tutor.NomorInduk
			errors.SysMessage = err3.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}

	} else {
		//insert
		_, err := tx.Exec(insertDet, &tutor.Details.EducationFront, &tutor.Details.EducationBack, &tutor.Details.Ktp, &tutor.Details.Sim, &tutor.Details.Npwp, &tutor.Details.GenderID, &tutor.Details.PobID, &tutor.Details.Dob, &tutor.Details.Phone, &tutor.Details.Email, &tutor.Details.StreetAddress, &tutor.Details.AddressID, &tutor.Details.InsSource, &tutor.Details.JoinDate, &tutor.ID, &tutor.UserID)
		if err != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = tutor.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}
	}

	errors.Message = structs.Success
	errors.Data = tutor.NomorInduk
	errors.Code = http.StatusOK
	tx.Commit()
	return &errors
}

//CreateTutors is the func that will insert multiple tutors at once (complete).
//please note that the email in request parameter is the username coming from Login func
func (*repo) CreateTutors(tutor structs.Tutor) *structs.ErrorMessage {
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	//start checking insert
	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = tutor.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	sqlQuery := "insert into tutors (nomor_induk, name, tutor_type_id, user_id) values (?, ?, ?, ?)"
	res, err := tx.Exec(sqlQuery, &tutor.NomorInduk, &tutor.Name, &tutor.TutorTypeID, &tutor.UserID)
	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = tutor.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	lastID, err := res.LastInsertId()
	lastTutorID := int(lastID)
	if err != nil {
		tx.Rollback()
		errors.Message = structs.LastIDErr
		errors.Data = tutor.NomorInduk
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	//insert details
	if tutor.Details != nil {
		tutor.Details.TutorID = lastTutorID
		tutor.Details.UserID = tutor.UserID

		sqlQueryDetail := "insert into tutor_details (education_degree_front, education_degree_back, ktp, sim, npwp, gender_id, pob_id, dob, phone, email, street_address, address_id, institution_source_name, join_date, tutor_id, user_id ) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		res, err := tx.Exec(sqlQueryDetail, &tutor.Details.EducationFront, &tutor.Details.EducationBack, &tutor.Details.Ktp, &tutor.Details.Sim, &tutor.Details.Npwp, &tutor.Details.GenderID, &tutor.Details.PobID, &tutor.Details.Dob, &tutor.Details.Phone, &tutor.Details.Email, &tutor.Details.StreetAddress, &tutor.Details.AddressID, &tutor.Details.InsSource, &tutor.Details.JoinDate, &lastTutorID, &tutor.UserID)
		if err != nil {
			tx.Rollback()
			errors.Message = structs.QueryErr
			errors.Data = tutor.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			return &errors
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			errors.Message = structs.LastIDErr
			errors.Data = tutor.NomorInduk
			errors.SysMessage = err.Error()
			errors.Code = http.StatusInternalServerError
			log.Println(lastID)
			return &errors
		}
	}

	//insert educations
	k := 0
	if len(tutor.Education) != 0 {
		sqlQueryEdu := "insert into tutor_educations (univ_degree_id, univ_name, years, tutor_id) values (?, ?, ?, ?)"
		for range tutor.Education {
			tutor.Education[k].TutorID = lastTutorID

			_, err2 := tx.Exec(sqlQueryEdu, &tutor.Education[k].UnivDegreeID, &tutor.Education[k].UnivName, &tutor.Education[k].Years, &lastTutorID)
			if err2 != nil {
				tx.Rollback()
				errors.Message = structs.QueryErr
				errors.Data = tutor.NomorInduk
				errors.SysMessage = err2.Error()
				errors.Code = http.StatusInternalServerError
				return &errors
			}
			k++
		}
	}

	//insert certificates
	m := 0
	if len(tutor.Certificate) != 0 {
		sqlQueryCert := "insert into tutor_certificates (cert_name, cert_date, tutor_id) values (?, ?, ?)"
		for range tutor.Certificate {
			tutor.Certificate[m].TutorID = lastTutorID

			_, err2 := tx.Exec(sqlQueryCert, &tutor.Certificate[m].CertName, &tutor.Certificate[m].CertDate, &lastTutorID)
			if err2 != nil {
				tx.Rollback()
				errors.Message = structs.QueryErr
				errors.Data = tutor.NomorInduk
				errors.SysMessage = err.Error()
				errors.Code = http.StatusInternalServerError
				return &errors
			}
			m++
		}
	}

	//insert Experiences
	n := 0
	if len(tutor.Experience) != 0 {
		sqlQueryExp := "insert into tutor_experiences (exp_name, description, years, tutor_id) values (?, ?, ?, ?)"
		for range tutor.Experience {
			tutor.Experience[n].TutorID = lastTutorID

			_, err2 := tx.Exec(sqlQueryExp, &tutor.Experience[n].ExpName, &tutor.Experience[n].Description, &tutor.Experience[n].Years, &lastTutorID)
			if err2 != nil {
				tx.Rollback()
				errors.Message = structs.QueryErr
				errors.SysMessage = err2.Error()
				errors.Code = http.StatusInternalServerError
				return &errors
			}
			n++
		}
	}

	//insert Journal
	p := 0
	if len(tutor.Journal) != 0 {
		sqlQueryJour := "insert into tutor_journals (journal_name, publish_at, publish_date, tutor_id) values (?, ?, ?, ?)"
		for range tutor.Journal {
			tutor.Journal[p].TutorID = lastTutorID

			_, err2 := tx.Exec(sqlQueryJour, &tutor.Journal[p].JourName, &tutor.Journal[p].PublishAt, &tutor.Journal[p].PublishDate, &lastTutorID)
			if err2 != nil {
				tx.Rollback()
				errors.Message = structs.QueryErr
				errors.SysMessage = err.Error()
				errors.Code = http.StatusInternalServerError
				return &errors
			}
			p++
		}
	}

	//insert research
	a := 0
	if len(tutor.Research) != 0 {
		sqlQueryRes := "insert into tutor_researches (res_name, description, years, tutor_id) values (?, ?, ?, ?)"
		for range tutor.Research {
			tutor.Research[a].TutorID = lastTutorID

			_, err2 := tx.Exec(sqlQueryRes, &tutor.Research[a].ResName, &tutor.Research[a].Description, &tutor.Research[a].Years, &lastTutorID)
			if err2 != nil {
				tx.Rollback()
				errors.Message = structs.QueryErr
				errors.SysMessage = err2.Error()
				errors.Code = http.StatusInternalServerError
				return &errors
			}
			a++
		}
	}

	errors.Message = structs.Success
	errors.Code = http.StatusOK

	tx.Commit()
	return &errors

}

//UpdateCertificates is the func to create/update the certificates in tutor entity. please note that status = 0 (soft delete)
func (*repo) UpdateCertificates(cert structs.TutorCertificate) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	var queryStr string
	var refID int

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = cert.CertName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	insertCert := "insert into tutor_certificates (cert_name, cert_date, status, tutor_id) values (?, ?, ?, ?)"
	updateCert := "update tutor_certificates set cert_name = ?, cert_date = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	if cert.ID != 0 {
		queryStr = updateCert
		refID = cert.ID
	} else {
		queryStr = insertCert
		refID = cert.TutorID
	}

	_, err2 := tx.Exec(queryStr, &cert.CertName, &cert.CertDate, &cert.Status, &refID)
	if err2 != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = cert.CertName
		errors.SysMessage = err2.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	errors.Message = structs.Success
	errors.Data = cert.CertName
	errors.Code = http.StatusOK
	tx.Commit()
	return &errors
}

//UpdateExperiences is the func to create/update the experiences in tutor entity. please note that status = 0 (soft delete)
func (*repo) UpdateExperiences(exp structs.TutorExperience) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	var queryStr string
	var refID int

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = exp.ExpName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	insertExp := "insert into tutor_experiences (exp_name, description, years, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateExp := "update tutor_experiences set exp_name = ?, description = ?, years = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	if exp.ID != 0 {
		queryStr = updateExp
		refID = exp.ID
	} else {
		queryStr = insertExp
		refID = exp.TutorID
	}

	_, err2 := tx.Exec(queryStr, &exp.ExpName, &exp.Description, &exp.Years, &exp.Status, &refID)
	if err2 != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = exp.ExpName
		errors.SysMessage = err2.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	errors.Message = structs.Success
	errors.Data = exp.ExpName
	errors.Code = http.StatusOK
	tx.Commit()
	return &errors
}

//UpdateJournals is the func to create/update the journal in tutor entity. please note that status = 0 (soft delete)
func (*repo) UpdateJournals(jour structs.TutorJournal) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	var queryStr string
	var refID int

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = jour.JourName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	insertJour := "insert into tutor_journals (journal_name, publish_at, publish_date, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateJour := "update tutor_journals set journal_name = ?, publish_at = ?, publish_date = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	if jour.ID != 0 {
		queryStr = updateJour
		refID = jour.ID
	} else {
		queryStr = insertJour
		refID = jour.TutorID
	}

	_, err2 := tx.Exec(queryStr, &jour.JourName, &jour.PublishAt, &jour.PublishDate, &jour.Status, &refID)
	if err2 != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = jour.JourName
		errors.SysMessage = err2.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	errors.Message = structs.Success
	errors.Data = jour.JourName
	errors.Code = http.StatusOK
	tx.Commit()
	return &errors
}

//UpdateEducations is the func to create/update the education in tutor entity. please note that status = 0 (soft delete)
func (*repo) UpdateEducations(edu structs.TutorEducation) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	var queryStr string
	var refID int

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = edu.UnivName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	insertEdu := "insert into tutor_educations (univ_degree_id, univ_name, years, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateEdu := "update tutor_educations set univ_degree_id = ?, univ_name = ?, years = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	if edu.ID != 0 {
		queryStr = updateEdu
		refID = edu.ID
	} else {
		queryStr = insertEdu
		refID = edu.TutorID
	}

	_, err2 := tx.Exec(queryStr, &edu.UnivDegreeID, &edu.UnivName, &edu.Years, &edu.Status, &refID)
	if err2 != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = edu.UnivName
		errors.SysMessage = err2.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	errors.Message = structs.Success
	errors.Data = edu.UnivName
	errors.Code = http.StatusOK
	tx.Commit()
	return &errors
}

//UpdateResearches is the func to create/update the journal in tutor entity. please note that status = 0 (soft delete)
func (*repo) UpdateResearches(res structs.TutorResearch) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	var queryStr string
	var refID int

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = res.ResName
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	insertRsch := "insert into tutor_researches (res_name, description, years, status, tutor_id) values (?, ?, ?, ?, ?)"
	updateRsch := "update tutor_researches set res_name = ?, description = ?, years = ?, status = ?, updated_at = now(), updated_by = 'API' where id = ?"

	if res.ID != 0 {
		queryStr = updateRsch
		refID = res.ID
	} else {
		queryStr = insertRsch
		refID = res.ID
	}

	_, err2 := tx.Exec(queryStr, &res.ResName, &res.Description, &res.Years, &res.Status, &refID)
	if err2 != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = res.ResName
		errors.SysMessage = err2.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	errors.Message = structs.Success
	errors.Data = res.ResName
	errors.Code = http.StatusOK
	tx.Commit()
	return &errors
}

//CheckEmail is the func to check registered/updated email
func (*repo) CheckEmail(email string, usrID int) int {
	db := mysql.InitializeMySQL()
	sqlQueryCheck := "select count(1) from users where username = ? and id != ?"
	check := 0
	err := db.QueryRow(sqlQueryCheck, &email, &usrID).Scan(&check)

	if err != nil {
		check = 99
	}
	return check
}

//checkNomorInduk is the func to check registered/updated nomor induk
func (*repo) CheckNomorInduk(insID int, nmrInduk string, tutorID int) int {
	db := mysql.InitializeMySQL()
	sqlQueryCheck := "SELECT count(1) FROM tutors ttr inner join (select user_id from user_roles where institution_id = ?) ur on ttr.user_id = ur.user_id where ttr.nomor_induk = ? "
	check := 0
	if tutorID != 0 {
		sqlQueryCheck += "and ttr.id != ?"
		err := db.QueryRow(sqlQueryCheck, &insID, &nmrInduk, &tutorID).Scan(&check)
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
