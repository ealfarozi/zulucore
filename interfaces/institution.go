package interfaces

import "github.com/ealfarozi/zulucore/structs"

type InstitutionRepository interface {
	GetInstitutions(page string, limit string) (*[]structs.Institution, *structs.ErrorMessage)
	GetInstitution(insID string, insCode string, page string, limit string) (*structs.Institution, *structs.ErrorMessage)
	CreateInstitutions(insts structs.Institution) *structs.ErrorMessage
}
