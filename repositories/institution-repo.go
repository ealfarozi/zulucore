package repositories

import (
	"net/http"
	"strconv"

	"github.com/ealfarozi/zulucore/interfaces"
	"github.com/ealfarozi/zulucore/repositories/mysql"
	"github.com/ealfarozi/zulucore/structs"
)

type insRepo struct{}

//NewInstitutionRepository is the constructor for Institution repository
func NewInstitutionRepository() interfaces.InstitutionRepository {
	return &insRepo{}
}

//CreateInstitutions is the func for creating the institutions
func (*insRepo) CreateInstitutions(insts structs.Institution) *structs.ErrorMessage {
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()
	tx, err := db.Begin()

	//start checking insert
	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = insts.Name
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	sqlQuery := "insert into institutions (code, name, street_address, street_map_id, bill_address, bill_map_id, pic_name, pic_phone, expired_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	res, err := tx.Exec(sqlQuery, &insts.Code, &insts.Name, &insts.Street, &insts.MapID, &insts.BillStreet, &insts.BillMapID, &insts.PICName, &insts.PICPhone, &insts.ExpireAt)
	if err != nil {
		tx.Rollback()
		errors.Message = structs.QueryErr
		errors.Data = insts.Name
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	lastID, err := res.LastInsertId()
	lastInsID := int(lastID)
	if err != nil {
		tx.Rollback()
		errors.Message = structs.LastIDErr
		errors.Data = insts.Name
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return &errors
	}

	errors.Message = structs.Success
	errors.Data = strconv.Itoa(lastInsID)
	errors.Code = http.StatusOK

	tx.Commit()
	return &errors
}

//GetInstitution is func to fulfill the dropbox in FE
func (*insRepo) GetInstitution(insID string, insCode string) (*structs.Institution, *structs.ErrorMessage) {
	var prm string
	var FullMapID int
	var BillMapID int

	var ins structs.Institution
	var errors structs.ErrorMessage
	var addr structs.Address

	db := mysql.InitializeMySQL()

	sqlQuery := "select id, code, name, street_address, street_map_id, bill_address, bill_map_id, pic_name, pic_phone, expired_at, status from institutions where "

	if insID != "" {
		sqlQuery += "id = ?"
		prm = insID
	}
	if insCode != "" {
		sqlQuery += "code = ?"
		prm = insCode
	}

	err := db.QueryRow(sqlQuery, prm).Scan(&ins.ID, &ins.Code, &ins.Name, &ins.Street, &FullMapID, &ins.BillStreet, &BillMapID, &ins.PICName, &ins.PICPhone, &ins.ExpireAt, &ins.Status)
	if err != nil {
		errors.Message = structs.ErrNotFound
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

	//get Address
	sqlQueryAddress := "select id, province_id, province_name, city_id, city_name, kecamatan_id, kecamatan_name, kelurahan_id, kelurahan_name, zipcode from address_map where id = ?"
	res, err := db.Query(sqlQueryAddress, FullMapID)
	res2, err2 := db.Query(sqlQueryAddress, BillMapID)
	defer mysql.CloseRows(res)
	defer mysql.CloseRows(res2)
	if err != nil || err2 != nil {
		errors.Message = structs.QueryErr
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}
	for res.Next() {
		res.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
		ins.FullAddress = addr
	}
	for res2.Next() {
		res2.Scan(&addr.ID, &addr.ProvinceID, &addr.ProvinceName, &addr.CityID, &addr.CityName, &addr.KecamatanID, &addr.KecamatanName, &addr.KelurahanID, &addr.KelurahanName, &addr.ZipCode)
		ins.BillFullAddress = addr
	}

	return &ins, nil
}

//GetInstitutions is func to fulfill the dropbox in FE
func (*insRepo) GetInstitutions() (*[]structs.Institution, *structs.ErrorMessage) {

	var institutions []structs.Institution
	var errors structs.ErrorMessage

	db := mysql.InitializeMySQL()

	sqlQuery := "select id, code, name, street_address, street_map_id, bill_address, bill_map_id, pic_name, pic_phone, expired_at, status from institutions"

	res, err := db.Query(sqlQuery)
	defer mysql.CloseRows(res)
	if err != nil {
		errors.Message = structs.QueryErr
		errors.SysMessage = err.Error()
		errors.Code = http.StatusInternalServerError
		return nil, &errors
	}

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
				errors.Message = structs.QueryErr
				errors.SysMessage = err.Error()
				errors.Code = http.StatusInternalServerError
				return nil, &errors
			}
			if err2 != nil {
				errors.Message = structs.QueryErr
				errors.SysMessage = err2.Error()
				errors.Code = http.StatusInternalServerError
				return nil, &errors
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
	return &institutions, nil
}
