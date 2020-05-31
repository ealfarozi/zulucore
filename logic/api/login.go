package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ealfarozi/zulucore/common"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
	"golang.org/x/crypto/bcrypt"
)

//@todo need to put in env
var expMins time.Duration = 60
var myKey = []byte(mysql.ViperEnvVariable("JWT_KEY"))

//Login and generate the JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ath structs.Auth
	var rol structs.Role
	var ins structs.Institution
	var usr structs.User
	var passw string

	_ = json.NewDecoder(r.Body).Decode(&usr)

	db := mysql.InitializeMySQL()

	sqlQuery := "SELECT id, username, password FROM users where username = ?"
	err := db.QueryRow(sqlQuery, usr.Username).Scan(&ath.Data.UserID, &ath.Data.Username, &passw)

	if err != nil {
		common.JSONError(w, structs.UserNotFound, err.Error(), http.StatusInternalServerError)
		return
	}
	if comparePasswords(passw, []byte(usr.Password)) {
		tokenString, err := GenerateJWT()
		if err != nil {
			common.JSONError(w, structs.GenTokenErr, err.Error(), http.StatusInternalServerError)
			return
		}

		//get user roles
		sqlQueryRole := "select ro.id, ro.name from user_roles ur inner join roles ro on ur.role_id = ro.id and ur.user_id = ?"
		res, err := db.Query(sqlQueryRole, ath.Data.UserID)
		defer mysql.CloseRows(res)
		if err != nil {
			common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
			return
		}
		for res.Next() {
			res.Scan(&rol.RoleID, &rol.RoleName)
			ath.Roles = append(ath.Roles, rol)
		}

		//get user institutions
		sqlQueryIns := "select ins.id, ins.code, ins.name from user_roles ur inner join institutions ins on ur.institution_id = ins.id and ur.user_id = ?"
		res2, err := db.Query(sqlQueryIns, ath.Data.UserID)
		defer mysql.CloseRows(res2)
		if err != nil {
			common.JSONError(w, structs.QueryErr, err.Error(), http.StatusInternalServerError)
			return
		}
		for res2.Next() {
			res2.Scan(&ins.ID, &ins.Code, &ins.Name)
			ath.Institutions = append(ath.Institutions, ins)
		}

		ath.Message = structs.Success
		ath.Data.AccessToken = tokenString
		ath.Data.ExpireAt = time.Now().Add(time.Minute * expMins)

		json.NewEncoder(w).Encode(ath)
	} else {
		common.JSONError(w, structs.IncorrectPass, "", http.StatusInternalServerError)
		return
	}
}

//IsAuthorized is the func for validating the JWT token
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}

				return myKey, nil
			})

			if err != nil {
				common.JSONError(w, structs.Unauthorized, err.Error(), http.StatusUnauthorized)
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			common.JSONError(w, structs.AuthNotFound, "", http.StatusUnauthorized)
			return
		}
	})
}

//GenerateJWT is func to generate the token
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	//claims["user"] = "Erlangga"
	//@todo need to define the expired time
	claims["exp"] = time.Now().Add(time.Minute * expMins).Unix()

	tokenString, err := token.SignedString(myKey)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenString, nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
