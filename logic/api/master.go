package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

//CreateInstitutions is the func for creating the institutions
func CreateInstitutions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var insts []structs.Institution
	var errstr structs.ErrorMessage
	_ = json.NewDecoder(r.Body).Decode(&insts)
	db := mysql.InitializeMySQL()
	sqlQuery := "insert into institutions (code, name, street_address, street_map_id, bill_address, bill_map_id, pic_name, pic_phone, expired_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	tx, err := db.Begin()
	j := 0
	for range insts {
		if err != nil {
			tx.Rollback()
			common.JSONError(w, insts[j].Name, err.Error(), http.StatusInternalServerError)
			return
		}

		v := validator.New()
		err := v.Struct(insts[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := tx.Exec(sqlQuery, &insts[j].Code, &insts[j].Name, &insts[j].Street, &insts[j].MapID, &insts[j].BillStreet, &insts[j].BillMapID, &insts[j].PICName, &insts[j].PICPhone, &insts[j].ExpireAt)

		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
			return
		}

		lastInsertedID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.LastIDErr, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("success insert into institution with ID:", lastInsertedID)
		errstr.Message = structs.Success
		errstr.Code = http.StatusOK
		j++
	}

	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}

//GetInstitutions is func to fulfill the dropbox in FE
func GetInstitutions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var institutions []structs.Institution
	//var FullMapID, BillMapID int
	db := mysql.InitializeMySQL()

	sqlQuery := "select id, code, name, street_address, street_map_id, bill_address, bill_map_id, pic_name, pic_phone, expired_at, status from institutions"

	res, err := db.Query(sqlQuery)
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {

		for res.Next() {
			institution := structs.Institution{}
			addr := structs.Address{}
			res.Scan(&institution.ID, &institution.Code, &institution.Name, &institution.Street, &institution.MapID, &institution.BillStreet, &institution.BillMapID, &institution.PICName, &institution.PICPhone, &institution.ExpireAt, &institution.Status)

			if institution.MapID != 0 && institution.BillMapID != 0 {

				//get Address
				sqlQueryAddress := "select id, province_id, province_name, city_id, city_name, kecamatan_id, kecamatan_name, kelurahan_id, kelurahan_name, zipcode from address_map where id = ?"
				res, err := db.Query(sqlQueryAddress, institution.MapID)
				res2, err2 := db.Query(sqlQueryAddress, institution.BillMapID)
				defer mysql.CloseRows(res)
				defer mysql.CloseRows(res2)
				if err != nil {
					common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
					return
				}
				if err2 != nil {
					common.JSONError(w, structs.QueryErr, err2.Error(), http.StatusInternalServerError)
					return
				}
				for res.Next() {
					res.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
					institution.FullAddress = addr
				}
				for res2.Next() {
					res2.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
					institution.BillFullAddress = addr
				}
			}
			institutions = append(institutions, institution)
		}
		json.NewEncoder(w).Encode(institutions)
	}
}

//GetInstitution is func to fulfill the dropbox in FE
func GetInstitution(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var prm string
	var FullMapID int
	var BillMapID int

	institution := structs.Institution{}
	addr := structs.Address{}

	params := mux.Vars(r)
	db := mysql.InitializeMySQL()

	sqlQuery := "select id, code, name, street_address, street_map_id, bill_address, bill_map_id, pic_name, pic_phone, expired_at, status from institutions where "

	if r.FormValue("insId") != "" {
		sqlQuery += "id = ?"
		prm = r.FormValue("insId")
	}
	if r.FormValue("insCode") != "" {
		sqlQuery += "code = ?"
		prm = r.FormValue("insCode")
		log.Println(params)
	}

	err := db.QueryRow(sqlQuery, prm).Scan(&institution.ID, &institution.Code, &institution.Name, &institution.Street, &FullMapID, &institution.BillStreet, &BillMapID, &institution.PICName, &institution.PICPhone, &institution.ExpireAt, &institution.Status)
	if err != nil {
		common.JSONError(w, structs.ErrNotFound, err.Error(), http.StatusInternalServerError)
		return
	}

	//get Address
	sqlQueryAddress := "select id, province_id, province_name, city_id, city_name, kecamatan_id, kecamatan_name, kelurahan_id, kelurahan_name, zipcode from address_map where id = ?"
	res, err := db.Query(sqlQueryAddress, FullMapID)
	res2, err2 := db.Query(sqlQueryAddress, BillMapID)
	defer mysql.CloseRows(res)
	defer mysql.CloseRows(res2)
	if err != nil || err2 != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}
	for res.Next() {
		res.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
		institution.FullAddress = addr
	}
	for res2.Next() {
		res2.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
		institution.BillFullAddress = addr
	}

	json.NewEncoder(w).Encode(institution)
}

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
