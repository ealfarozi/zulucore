package structs

//Tutor struct
type Tutor struct {
	ID         int32  `json:"id,omitempty"`
	NomorInduk string `json:"nomor_induk" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	TutorType  string `json:"type" validate:"required"`
}
