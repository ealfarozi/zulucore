package interfaces

import "github.com/ealfarozi/zulucore/structs"

type InstitutionRepository interface {
	GetInstitutions() (*[]structs.Institution, *structs.ErrorMessage)
	GetInstitution(insID string, insCode string) (*structs.Institution, *structs.ErrorMessage)
	CreateInstitutions(insts structs.Institution) *structs.ErrorMessage
}
