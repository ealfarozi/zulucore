package structs

//Student struct
type Student struct {
	ID          int             `json:"id,omitempty"`
	NomorInduk  string          `json:"nomor_induk" validate:"required"`
	Name        string          `json:"name" validate:"required"`
	DegreeID    int             `json:"degree_id" validate:"required"`
	StudentType int             `json:"type_id" validate:"required"`
	CurrID      int             `json:"curriculum_id" validate:"required"`
	UserID      int             `json:"user_id"`
	Status      int             `json:"status,omitempty"`
	InsID       int             `json:"institution_id,omitempty" validate:"required"`
	Scores      *Scores         `json:"scores,omitempty"`
	Details     *StudentDetails `json:"details,omitempty"`
	Parents     []Parents       `json:"parent,omitempty"`
}

//Scores struct
type Scores struct {
	ID         int     `json:"id,omitempty"`
	GradeID    int     `json:"grade_id"`
	Score      float32 `json:"score"`
	SubjectID  int     `json:"subject_id"`
	SemesterID int     `json:"semester_id"`
	StudentID  int     `json:"student_id"`
	Status     int     `json:"status,omitempty"`
}

//StudentDetails struct
type StudentDetails struct {
	ID            int     `json:"id,omitempty"`
	KkNO          string  `json:"kk_no,omitempty"`
	Ktp           string  `json:"ktp,omitempty"`
	Sim           string  `json:"sim,omitempty"`
	Npwp          string  `json:"npwp,omitempty"`
	GenderID      int     `json:"gender_id" validate:"required"`
	PobID         int     `json:"pob_id" validate:"required"`
	Dob           string  `json:"dob" validate:"required"`
	Phone         string  `json:"phone" validate:"min=8,max=20,startswith=08"`
	Email         string  `json:"email"  validate:"email"`
	InsSource     string  `json:"institution_name,omitempty"`
	JoinDate      string  `json:"join_date,omitempty"`
	StudentID     int     `json:"student_id,omitempty"`
	TutorID       int     `json:"tutor_id,omitempty"`
	UserID        int     `json:"user_id,omitempty"`
	StreetAddress string  `json:"street_address,omitempty"`
	AddressID     int     `json:"address_map_id,omitempty"`
	AddressDetail Address `json:"address_detail,omitempty"`
}

//Parent struct
type Parents struct {
	ID            int             `json:"id,omitempty"`
	Name          string          `json:"name" validate:"required"`
	Phone         string          `json:"phone,omitempty" validate:"min=8,max=20,startswith=08"`
	Email         string          `json:"email,omitempty"  validate:"email"`
	GenderID      int             `json:"gender_id" validate:"required"`
	UserID        int             `json:"user_id,omitempty" validate:"required"`
	ProfessionID  int             `json:"profession_id,omitempty" validate:"required"`
	StreetAddress string          `json:"street_address,omitempty"`
	AddressID     int             `json:"address_map_id,omitempty"`
	AddressDetail Address         `json:"address_detail,omitempty"`
	Students      []StudentParent `json:"students,omitempty"`
}

//Student-Parent struct
type StudentParent struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	InsID   int    `json:"institution_id,omitempty"`
	InsName string `json:"institution_name,omitempty"`
}
