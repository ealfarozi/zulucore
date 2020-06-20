package repositories

import (
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
)

type repo struct{}

func NewMysqlRepository() interfaces.TutorRepository {
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

func (*repo) GetTutors(insID string) (*[]structs.Tutor, *structs.ErrorMessage) {
	var tutors []structs.Tutor
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT ttr.id, ttr.nomor_induk, ttr.name, ttr.tutor_type_id, ttr.user_id, ttr.status FROM tutors ttr inner join (select user_id from user_roles where institution_id = ?) ur on ttr.user_id = ur.user_id"

	//var tutor structs.Tutor

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
