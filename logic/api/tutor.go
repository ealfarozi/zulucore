package api

import (
	"encoding/json"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/repositories"
	"github.com/ealfarozi/zulucore/service"
	"github.com/ealfarozi/zulucore/structs"

	//blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type logic struct{}

var (
	tutorService service.TutorService
	repo         interfaces.TutorRepository = repositories.NewTutorRepository()
)

//TutorLogic is the interface of httpRequest for Tutor
type TutorLogic interface {
	GetTutors(w http.ResponseWriter, r *http.Request)
	GetTutorDetails(w http.ResponseWriter, r *http.Request)
	GetTutor(w http.ResponseWriter, r *http.Request)
	UpdateTutorDetails(w http.ResponseWriter, r *http.Request)
	UpdateExperiences(w http.ResponseWriter, r *http.Request)
	UpdateEducations(w http.ResponseWriter, r *http.Request)
	UpdateCertificates(w http.ResponseWriter, r *http.Request)
	UpdateJournals(w http.ResponseWriter, r *http.Request)
	UpdateResearches(w http.ResponseWriter, r *http.Request)
	CreateTutors(w http.ResponseWriter, r *http.Request)
}

//NewTutorLogic is the func to calling the constructor of tutor interface
func NewTutorLogic(service service.TutorService) TutorLogic {
	tutorService = service
	return &logic{}
}

//GetTutors in the db (all)
func (*logic) GetTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tutors, errStr := tutorService.GetTutors(r.FormValue("institution_id"), r.FormValue("_page"), r.FormValue("_limit"))

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

	tutors, errStr := tutorService.GetTutor(r.FormValue("nomor_induk"), r.FormValue("name"), r.FormValue("institution_id"), r.FormValue("_page"), r.FormValue("_limit"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(tutors)

}

//UpdateEducations is the func to create/update the education in tutor entity. please note that status = 0 (soft delete)
func (*logic) UpdateEducations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var edus []structs.TutorEducation
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&edus)

	j := 0
	for range edus {
		errStr := tutorService.ValidateEdu(&edus[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		errStr = tutorService.UpdateEducations(edus[j])
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		} else {
			errs = append(errs, structs.ErrorMessage{Data: edus[j].UnivName, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
		j++
	}

	json.NewEncoder(w).Encode(errs)
}

//UpdateCertificates is the func to create/update the certificates in tutor entity. please note that status = 0 (soft delete)
func (*logic) UpdateCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var certs []structs.TutorCertificate
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&certs)

	j := 0
	for range certs {
		errStr := tutorService.ValidateCert(&certs[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		errStr = tutorService.UpdateCertificates(certs[j])
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		} else {
			errs = append(errs, structs.ErrorMessage{Data: certs[j].CertName, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
		j++
	}

	json.NewEncoder(w).Encode(errs)
}

//UpdateExperiences is the func to create/update the experiences in tutor entity. please note that status = 0 (soft delete)
func (*logic) UpdateExperiences(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var exps []structs.TutorExperience
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&exps)

	j := 0
	for range exps {
		errStr := tutorService.ValidateExp(&exps[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		errStr = tutorService.UpdateExperiences(exps[j])
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		} else {
			errs = append(errs, structs.ErrorMessage{Data: exps[j].ExpName, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
		j++
	}

	json.NewEncoder(w).Encode(errs)
}

//UpdateJournals is the func to create/update the journal in tutor entity. please note that status = 0 (soft delete)
func (*logic) UpdateJournals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jours []structs.TutorJournal
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&jours)

	j := 0
	for range jours {
		errStr := tutorService.ValidateJour(&jours[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		errStr = tutorService.UpdateJournals(jours[j])
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		} else {
			errs = append(errs, structs.ErrorMessage{Data: jours[j].JourName, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
		j++
	}

	json.NewEncoder(w).Encode(errs)
}

//UpdateResearches is the func to create/update the journal in tutor entity. please note that status = 0 (soft delete)
func (*logic) UpdateResearches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rsch []structs.TutorResearch
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&rsch)

	j := 0
	for range rsch {
		errStr := tutorService.ValidateRes(&rsch[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		errStr = tutorService.UpdateResearches(rsch[j])
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		} else {
			errs = append(errs, structs.ErrorMessage{Data: rsch[j].ResName, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
		j++
	}

	json.NewEncoder(w).Encode(errs)
}

//UpdateTutorDetails is the func to create/update the tutor detail (ONLY) on Frontend side for tutor entity. The update will includes nomor_induk and tutor_name as well.
//Please note that Tutor.status = 0 (soft delete). In order to create a new tutor please refer to CreateTutors func.
//Email field should be coming from Login func.
//Both of ID (tutor and tutor_details) are needed in this function
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
		} else {
			errs = append(errs, structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
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

	for j := range tutors {
		tutor, errStr := tutorService.Validate(&tutors[j])
		if errStr != nil {
			errs = append(errs, *errStr)
			continue
		}

		checkNomorInduk := tutorService.CheckNomorInduk(tutors[j].InsID, tutors[j].NomorInduk, 0)
		if checkNomorInduk != 0 {
			errStr := structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.NomorInd, SysMessage: "", Code: http.StatusInternalServerError}
			errs = append(errs, errStr)
			continue
		}

		if tutors[j].Details != nil {
			checkEmail := tutorService.CheckEmail(tutors[j].Details.Email, tutors[j].UserID)
			if checkEmail != 0 {
				errStr := structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.Email, SysMessage: "", Code: http.StatusInternalServerError}
				errs = append(errs, errStr)
				continue
			}
		}

		errStr = tutorService.CreateTutors(*tutor)
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
			continue
		} else {
			errs = append(errs, structs.ErrorMessage{Data: tutors[j].NomorInduk, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
	}
	common.JSONErrs(w, &errs)
	return

}
