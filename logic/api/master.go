package api

import (
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

		res, err := tx.Exec(sqlQuery, &insts[j].Code, &insts[j].Name, &insts[j].Street)
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

	db := mysql.InitializeMySQL()

	sqlQuery := "select id, code, name, street_address, street_map_id, bill_address, bill_map_id, pic_name, pic_phone, expired_at, status from institutions"

	res, err := db.Query(sqlQuery)
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {
		institution := structs.Institution{}
		for res.Next() {
			res.Scan(&institution.ID, &institution.Code, &institution.Name, &institution.Street, &institution.FullAddress, &institution.BillStreet, &institution.BillFullAddress, &institution.PICName, &institution.PICPhone, &institution.ExpireAt, &institution.Status)
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
	//var addr structs.Address

	params := mux.Vars(r)
	institution := structs.Institution{}
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

//GetRoles is func to fulfill the dropbox in FE
func GetRoles(w http.ResponseWriter, r *http.Request) {

}
