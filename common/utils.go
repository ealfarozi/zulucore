package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func SetPageLimit(page string, limit string) string {
	var offset int
	if page == "" {
		return " limit 100 offset 0"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return ""
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return ""
	}

	offset = (pageInt - 1) * limitInt
	ret := fmt.Sprintf(" limit %d offset %d", limitInt, offset)
	return ret

}
