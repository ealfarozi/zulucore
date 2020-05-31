package structs

//Tutor struct
type Tutor struct {
	ID          int32              `json:"id,omitempty"`
	NomorInduk  string             `json:"nomor_induk" validate:"required"`
	Name        string             `json:"name" validate:"required"`
	TutorTypeID int                `json:"type" validate:"required"`
	UserID      int32              `json:"user_id" validate:"required"`
	Status      int32              `json:"status"`
	Details     *TutorDetails      `json:"details,omitempty"`
	Education   []TutorEducation   `json:"education,omitempty"`
	Certificate []TutorCertificate `json:"certificate,omitempty"`
	Experience  []TutorExperience  `json:"experience,omitempty"`
	Research    []TutorResearch    `json:"research,omitempty"`
	Journal     []TutorJournal     `json:"journal,omitempty"`
}

//TutorDetails is the struct to get the details of a tutor
type TutorDetails struct {
	ID             int     `json:"id,omitempty"`
	EducationFront string  `json:"degree_front,omitempty"`
	EducationBack  string  `json:"degree_back,omitempty"`
	Ktp            string  `json:"ktp,omitempty"`
	Sim            string  `json:"sim,omitempty"`
	Npwp           string  `json:"npwp,omitempty"`
	GenderID       int     `json:"gender_id" validate:"required"`
	PobID          int     `json:"pob_id" validate:"required"`
	Dob            string  `json:"dob" validate:"required"`
	Phone          string  `json:"phone" validate:"required,min=8,max=20,startswith=08"`
	Email          string  `json:"email"  validate:"required,email"`
	StreetAddress  string  `json:"street_address,omitempty"`
	AddressID      int     `json:"address_map_id,omitempty"`
	AddressDetail  Address `json:"address_detail,omitempty"`
	InsSource      string  `json:"institution_name,omitempty"`
	JoinDate       string  `json:"join_date,omitempty"`
	TutorID        int32   `json:"tutor_id" validate:"required"`
	UserID         int32   `json:"user_id" validate:"required"`
}

//TutorEducation is the struct to get the education details of a tutor
type TutorEducation struct {
	ID           int    `json:"id,omitempty"`
	UnivDegreeID int    `json:"univ_degree_id" validate:"required"`
	UnivName     string `json:"univ_name" validate:"required"`
	Years        int    `json:"years" validate:"required"`
	TutorID      int32  `json:"tutor_id" validate:"required"`
}

//TutorCertificate is the struct to get the certificate details of a tutor
type TutorCertificate struct {
	ID       int    `json:"id,omitempty"`
	CertName int    `json:"cert_name" validate:"required"`
	CertDate string `json:"cert_date" validate:"required"`
	TutorID  int32  `json:"tutor_id" validate:"required"`
}

//TutorExperience is the struct to get the experience details of a tutor
type TutorExperience struct {
	ID          int    `json:"id,omitempty"`
	ExpName     int    `json:"exp_name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Years       int    `json:"years" validate:"required"`
	TutorID     int32  `json:"tutor_id" validate:"required"`
}

//TutorResearch is the struct to get the research list of a tutor
type TutorResearch struct {
	ID          int    `json:"id,omitempty"`
	ResName     int    `json:"research_name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Years       int    `json:"years" validate:"required"`
	TutorID     int32  `json:"tutor_id" validate:"required"`
}

//TutorJournal is the struct to get the journal list of a tutor
type TutorJournal struct {
	ID          int    `json:"id,omitempty"`
	JourName    string `json:"journal_name" validate:"required"`
	PublishAt   string `json:"publish_at" validate:"required"`
	PublishDate string `json:"publish_date" validate:"required"`
	TutorID     int32  `json:"tutor_id" validate:"required"`
}
