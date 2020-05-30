package api

import (
	"encoding/json"
	"net/http"

	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
	"gopkg.in/go-playground/validator.v9"
)

//CreateUsers is the func for creating multiple user(tutor, student) at once
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	//create user, user role, tutor/siswa/admin
	w.Header().Set("Content-Type", "application/json")
}

//CreateUserLogin is the func for creating the username/email that will be used for logging in
func CreateUserLogin(w http.ResponseWriter, r *http.Request) {
	//create username/email with hashed default password: DOB - YYYYMMDD
	w.Header().Set("Content-Type", "application/json")

	var usr []structs.User
	var errstr structs.ErrorMessage

	_ = json.NewDecoder(r.Body).Decode(&usr)
	db := mysql.InitializeMySQL()
	sqlQuery := "INSERT INTO users (username, password) values (?,?)"
	sqlQueryRole := "INSERT INTO user_roles (user_id, role_id, institution_id) values (?,?,?)"
	tx, err := db.Begin()
	j := 0
	for range usr {
		if err != nil {
			tx.Rollback()
			common.JSONError(w, usr[j].Username, err.Error(), http.StatusInternalServerError)
			return
		}

		v := validator.New()
		err := v.Struct(usr[j])
		if err != nil {
			tx.Rollback()
			common.JSONError(w, structs.Validate, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := tx.Exec(sqlQuery, &usr[j].Username, hashAndSalt([]byte(usr[j].Password)))
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

		if &usr[j].RoleID != nil {
			_, err := tx.Exec(sqlQueryRole, &lastInsertedID, &usr[j].RoleID, &usr[j].InstitutionID)
			if err != nil {
				tx.Rollback()
				common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		errstr.Message = structs.Success
		errstr.Code = http.StatusOK
		j++
	}
	tx.Commit()
	json.NewEncoder(w).Encode(errstr)
}
