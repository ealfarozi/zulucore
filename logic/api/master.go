package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
)

//GetReferences is func to get any refs in references table
func GetReferences(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var references []structs.Reference
	params := r.FormValue("groupName")

	db := mysql.InitializeMySQL()

	sqlQuery := "select sub_id, name, status from `references` where group_name = ? order by sub_id asc"

	res, err := db.Query(sqlQuery, params)
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {
		reference := structs.Reference{}
		for res.Next() {
			res.Scan(&reference.ID, &reference.Values, &reference.Status)
			references = append(references, reference)
		}
		json.NewEncoder(w).Encode(references)
	}
}

//GetAddress is func to get any refs in references table
func GetAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var sqlQuery, prm string
	var allAddr []structs.Address

	db := mysql.InitializeMySQL()

	if r.FormValue("all") != "" {
		prm = r.FormValue("all")
		if prm == "*" {
			sqlQuery = "select id, province_id, province_name, city_id, city_name, kecamatan_id, kecamatan_name, kelurahan_id, kelurahan_name, zipcode from address_map"
		} else {
			sqlQuery = "select id, province_id, province_name, city_id, city_name, kecamatan_id, kecamatan_name, kelurahan_id, kelurahan_name, zipcode from address_map where kelurahan_name like ?"
		}
	}

	if r.FormValue("province") != "" {
		prm = r.FormValue("province")
		if prm == "*" {
			sqlQuery = "select id, name from provinces"
		} else {
			sqlQuery = "select id, name from provinces where name like ?"
		}
	}
	if r.FormValue("city") != "" {
		prm = r.FormValue("city")
		if prm == "*" {
			sqlQuery = "select id, name, province_id from cities"
		} else {
			sqlQuery = "select id, name, province_id from cities where name like ?"
		}
	}
	if r.FormValue("kecamatan") != "" {
		prm = r.FormValue("kecamatan")
		if prm == "*" {
			sqlQuery = "select id, name, city_id from kecamatan"
		} else {
			sqlQuery = "select id, name, city_id from kecamatan where name like ?"
		}
	}
	if r.FormValue("kelurahan") != "" {
		prm = r.FormValue("kelurahan")
		if prm == "*" {
			sqlQuery = "select id, name, kecamatan_id, zipcode from kelurahan"
		} else {
			sqlQuery = "select id, name, kecamatan_id, zipcode from kelurahan where name like ?"
		}
	}
	var res *sql.Rows
	var err error

	if prm == "*" {
		res, err = db.Query(sqlQuery)
	} else {
		prm = "%" + prm + "%"
		res, err = db.Query(sqlQuery, prm)
	}

	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {
		addr := structs.Address{}
		if r.FormValue("all") != "" {
			for res.Next() {
				res.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
				allAddr = append(allAddr, addr)
			}
		}

		if r.FormValue("province") != "" {
			for res.Next() {
				res.Scan(&addr.ProvinceID, &addr.ProvinceName)
				allAddr = append(allAddr, addr)
			}
		}
		if r.FormValue("city") != "" {
			for res.Next() {
				res.Scan(&addr.CityID, &addr.CityName, &addr.ProvinceID)
				allAddr = append(allAddr, addr)
			}
		}
		if r.FormValue("kecamatan") != "" {
			for res.Next() {
				res.Scan(&addr.KecamatanID, &addr.KecamatanName, &addr.CityID)
				allAddr = append(allAddr, addr)
			}
		}
		if r.FormValue("kelurahan") != "" {
			for res.Next() {
				res.Scan(&addr.KelurahanID, &addr.KelurahanName, &addr.KecamatanID, &addr.ZipCode)
				allAddr = append(allAddr, addr)
			}
		}
	}
	json.NewEncoder(w).Encode(allAddr)
}

//GetRoles is func to fulfill the dropbox in FE
func GetRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var roles []structs.Role

	db := mysql.InitializeMySQL()

	sqlQuery := "select id, name from roles where id != 1"

	res, err := db.Query(sqlQuery)
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {

		for res.Next() {
			role := structs.Role{}
			res.Scan(&role.RoleID, &role.RoleName)
			roles = append(roles, role)
		}
		json.NewEncoder(w).Encode(roles)
	}
}
