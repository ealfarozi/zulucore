package common

import (
	"encoding/json"
	"net/http"

	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
)

//JSONError is func to return JSON error format
func JSONError(w http.ResponseWriter, message string, sysMessage string, code int) {
	var errstr structs.ErrorMessage
	errstr.Message = message
	errstr.SysMessage = sysMessage
	errstr.Code = code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errstr)
}

func JSONErr(w http.ResponseWriter, errStr *structs.ErrorMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errStr.Code)
	json.NewEncoder(w).Encode(errStr)
}

func JSONErrs(w http.ResponseWriter, errStr *[]structs.ErrorMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(errStr)
}

//GetAddressOnly is func to get all address via SQL not from http request
func GetAddressOnly(addID int) structs.Address {
	var addr structs.Address

	db := mysql.InitializeMySQL()
	sqlQueryAddr := "select id, province_id, province_name, city_id, city_name, kecamatan_id, kecamatan_name, kelurahan_id, kelurahan_name, zipcode from address_map where id = ?"
	res, err := db.Query(sqlQueryAddr, addID)
	defer mysql.CloseRows(res)

	if err != nil {
		return addr
	}

	for res.Next() {
		res.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
	}
	return addr
}

//CheckNomorInduk is the func to check registered/updated nomor induk
func CheckNomorInduk(insID int, nmrInduk string, tutorID int) int {
	db := mysql.InitializeMySQL()
	sqlQueryCheck := "SELECT count(1) FROM tutors ttr inner join (select user_id from user_roles where institution_id = ?) ur on ttr.user_id = ur.user_id where ttr.nomor_induk = ? "
	check := 0
	if tutorID != 0 {
		sqlQueryCheck += "and ttr.id != ?"
		err := db.QueryRow(sqlQueryCheck, &insID, &nmrInduk, &tutorID).Scan(&check)
		if err != nil {
			check = 99
		}
	} else {
		err := db.QueryRow(sqlQueryCheck, &insID, &nmrInduk).Scan(&check)
		if err != nil {
			check = 99
		}
	}

	return check
}

func CheckNomorIndukStd(insID int, nmrInduk string, stdID int) int {
	db := mysql.InitializeMySQL()
	sqlQueryCheck := "SELECT count(1) FROM students std inner join (select user_id from user_roles where institution_id = ?) ur on std.user_id = ur.user_id where std.nomor_induk = ? "
	check := 0
	if stdID != 0 {
		sqlQueryCheck += "and std.id != ?"
		err := db.QueryRow(sqlQueryCheck, &insID, &nmrInduk, &stdID).Scan(&check)
		if err != nil {
			check = 99
		}
	} else {
		err := db.QueryRow(sqlQueryCheck, &insID, &nmrInduk).Scan(&check)
		if err != nil {
			check = 99
		}
	}

	return check
}

//CheckEmail is the func to check registered/updated email
func CheckEmail(email string, usrID int) int {
	db := mysql.InitializeMySQL()
	sqlQueryCheck := "select count(1) from users where username = ? and id != ?"
	check := 0
	err := db.QueryRow(sqlQueryCheck, &email, &usrID).Scan(&check)

	if err != nil {
		check = 99
	}
	return check
}
