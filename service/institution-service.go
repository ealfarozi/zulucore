package service

import (
	"net/http"

	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/structs"
	"gopkg.in/go-playground/validator.v9"
)

type InstitutionService interface {
	GetInstitutions(page string, limit string) (*[]structs.Institution, *structs.ErrorMessage)
	GetInstitution(insID string, insCode string, page string, limit string) (*structs.Institution, *structs.ErrorMessage)
	CreateInstitutions(ins structs.Institution) *structs.ErrorMessage
	ValidateInstitution(ins *structs.Institution) *structs.ErrorMessage
}

type insService struct{}

var (
	insRepo interfaces.InstitutionRepository
)

func NewInstitutionService(repository interfaces.InstitutionRepository) InstitutionService {
	insRepo = repository
	return &insService{}
}

func (*insService) CreateInstitutions(ins structs.Institution) *structs.ErrorMessage {
	return insRepo.CreateInstitutions(ins)
}

func (*insService) GetInstitution(insID string, insCode string, page string, limit string) (*structs.Institution, *structs.ErrorMessage) {
	return insRepo.GetInstitution(insID, insCode, page, limit)
}

func (*insService) GetInstitutions(page string, limit string) (*[]structs.Institution, *structs.ErrorMessage) {
	return insRepo.GetInstitutions(page, limit)
}

func (*insService) ValidateInstitution(ins *structs.Institution) *structs.ErrorMessage {
	var errors structs.ErrorMessage
	v := validator.New()
	err := v.Struct(ins)
	if err != nil {
		errors.Message = structs.Validate
		errors.Data = ins.Name
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}
	return nil
}
