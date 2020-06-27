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

type insLogic struct{}

var (
	insService service.InstitutionService
	insRepo    interfaces.InstitutionRepository = repositories.NewInstitutionRepository()
)

//InstitutionLogic is the interface for institutions
type InstitutionLogic interface {
	GetInstitutions(w http.ResponseWriter, r *http.Request)
	GetInstitution(w http.ResponseWriter, r *http.Request)
	CreateInstitutions(w http.ResponseWriter, r *http.Request)
}

//NewInstitutionLogic is the func to calling the construtor of institution interface
func NewInstitutionLogic(service service.InstitutionService) InstitutionLogic {
	insService = service
	return &insLogic{}
}

//CreateTutors is the func that will insert multiple tutors at once (complete).
//please note that the email in request parameter is the username coming from Login func
func (*insLogic) CreateInstitutions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var insts []structs.Institution
	var errs []structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&insts)

	j := 0
	for range insts {
		errStr := insService.ValidateInstitution(&insts[j])
		if errStr != nil {
			errs = append(errs, *errStr)
		}

		errStr = insService.CreateInstitutions(insts[j])
		if errStr.Code != http.StatusOK {
			errs = append(errs, *errStr)
		} else {
			errs = append(errs, structs.ErrorMessage{Data: insts[j].Name, Message: structs.Success, SysMessage: "", Code: http.StatusOK})
		}
		j++

	}
	common.JSONErrs(w, &errs)
	return
}

//GetInstitution in the db (all)
func (*insLogic) GetInstitution(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ins, errStr := insService.GetInstitution(r.FormValue("insId"), r.FormValue("insCode"))

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(ins)

}

//GetInstitutions in the db (all)
func (*insLogic) GetInstitutions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ins, errStr := insService.GetInstitutions()

	if errStr != nil {
		common.JSONErr(w, errStr)
		return
	}
	json.NewEncoder(w).Encode(ins)

}
