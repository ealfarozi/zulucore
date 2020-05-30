package structs

import (
	"time"
)

//User for UserLogin
type User struct {
	Username      string `json:"username" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8,max=100"`
	RoleID        int    `json:"role_id,omitempty"`
	InstitutionID int    `json:"institution_id,omitempty"`
}

//Auth struct for responding back the json
type Auth struct {
	Message      string        `json:"message"`
	Data         UserData      `json:"data"`
	Roles        []Role        `json:"roles"`
	Institutions []Institution `json:"institutions"`
}

//Role struct for nested auth
type Role struct {
	RoleID   int32  `json:"role_id" db:"role_id"`
	RoleName string `json:"role_name" db:"role_name"`
}

//UserData is parent struct for auth
type UserData struct {
	AccessToken string    `json:"access_token"`
	ExpireAt    time.Time `json:"expire_at"`
	UserID      int32     `json:"user_id" db:"id"`
	Username    string    `json:"username" db:"username"`
}

//Institution is the struct to get institutions
type Institution struct {
	ID          int64     `json:"id,omitempty"`
	Code        string    `json:"code" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Street      string    `json:"street_address,omitempty" validate:"required"`
	FullAddress []Address `json:"full_address,omitempty"`
}

//Address is the struct to get Map address
type Address struct {
	ID            int    `json:"id"`
	ProvinceID    int    `json:"province_id"`
	ProvinceName  string `json:"province_name,omitempty"`
	CityID        int    `json:"city_id"`
	CityName      string `json:"city_name,omitempty"`
	KecamatanID   int    `json:"kecamatan_id"`
	KecamatanName string `json:"kecamatan_name,omitempty"`
	KelurahanID   int    `json:"kelurahan_id"`
	KelurahanName string `json:"kelurahan_name,omitempty"`
	ZipCodeID     int    `json:"zipcode_id"`
	ZipCode       string `json:"zipcode,omitempty"`
}
