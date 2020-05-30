package api

import (
	"encoding/json"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"

	//blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//GetTutors in the db (all)
func GetTutors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tutors []structs.Tutor

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT id, nomor_induk, email FROM tutors"

	//var tutor structs.Tutor

	res, err := db.Query(sqlQuery)
	defer mysql.CloseRows(res)
	if err != nil {
		common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {
		tutor := structs.Tutor{}
		for res.Next() {
			res.Scan(&tutor.ID, &tutor.NomorInduk, &tutor.Email)
			tutors = append(tutors, tutor)
		}
		json.NewEncoder(w).Encode(tutors)
	}
}

//GetTutor in the db based on nomor_induk parameter
func GetTutor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //get params

	tutor := structs.Tutor{}
	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT id, nomor_induk, email FROM tutors where nomor_induk = ?"

	err := db.QueryRow(sqlQuery, params["tutorid"]).Scan(&tutor.ID, &tutor.NomorInduk, &tutor.Email)
	if err != nil {
		common.JSONError(w, structs.ErrNotFound, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tutor)

}

//CreateTutors (multiple)
func CreateTutors(w http.ResponseWriter, r *http.Request) {

}

//UpdateTutors (multiple)
func UpdateTutors(w http.ResponseWriter, r *http.Request) {

}

//DeleteTutors (multiple). this func will soft-delete the data
func DeleteTutors(w http.ResponseWriter, r *http.Request) {

}
