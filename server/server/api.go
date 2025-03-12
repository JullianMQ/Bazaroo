package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type OkResponse struct {
	Message string `json:"message"`
}

func ErrorRes(res http.ResponseWriter, status int, mess string) {
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(ErrorResponse{
		Error: mess,
	})
}


func GetRoot(res http.ResponseWriter, req *http.Request) {
	myMess := "Message"
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(myMess)
}

type AddrRequest struct {
	Addr_line1  string `json:"addr_line1"`
	Addr_line2  string `json:"addr_line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	Postal_code string `json:"postal_code"`
	Country     string `json:"country"`
}

type AddrResponse struct {
	Addr_id     int            `json:"addr_id"`
	Addr_line1  string         `json:"addr_line1"`
	Addr_line2  sql.NullString `json:"addr_line2"`
	City        string         `json:"city"`
	State       string         `json:"state"`
	Postal_code string         `json:"postal_code"`
	Country     string         `json:"country"`
}

func GetAddr(res http.ResponseWriter, req *http.Request) {
	rows, err := db.Query(`SELECT
		addr_id,
		addr_line1,
		addr_line2,
		city,
		state,
		postal_code,
		country
		FROM addresses`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var addr []AddrResponse
	for rows.Next() {
		var (
			addr_id     int
			addr_line1  string
			addr_line2  sql.NullString
			city        string
			state       string
			postal_code string
			country     string
		)
		if err := rows.Scan(
			&addr_id,
			&addr_line1,
			&addr_line2,
			&city,
			&state,
			&postal_code,
			&country); err != nil {
			log.Fatal(err)
		}
		addr = append(addr, AddrResponse{
			Addr_id:     addr_id,
			Addr_line1:  addr_line1,
			Addr_line2:  addr_line2,
			City:        city,
			State:       state,
			Postal_code: postal_code,
			Country:     country,
		})
	}
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(addr)
}

func PostAddr(res http.ResponseWriter, req *http.Request) {
	var addr *AddrRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&addr)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
	}

	if addr.Addr_line1 == "" {
		ErrorRes(res, http.StatusBadRequest,
			"addr_line1 is required")
		return
	}

	if addr.City == "" {
		ErrorRes(res, http.StatusBadRequest,
			"city is required")
		return
	}

	if addr.State == "" {
		ErrorRes(res, http.StatusBadRequest,
			"state is required")
		return
	}

	if addr.Postal_code == "" {
		ErrorRes(res, http.StatusBadRequest,
			"postal_code is required")
		return
	}

	if len(addr.Postal_code) > 10 {
		ErrorRes(res, http.StatusBadRequest,
			"postal_code must be 10 characters or less")
		return
	}

	if addr.Country == "" {
		ErrorRes(res, http.StatusBadRequest,
			"country is required")
		return
	}

	if len(addr.Country) != 3 {
		ErrorRes(res, http.StatusBadRequest,
			"country must be 3 characters, of iso 3166-1 alpha-3 code")
		return
	}

	rows, err := AddAddr(addr)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not add address, try again later."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Address added successfully. %d rows affected", rows),
	})
}

type Office struct {
	Office_id int    `json:"office_id"`
	Phone_num string `json:"phone_num"`
	Addr_id   int `json:"addr_id"`
}

func GetOffices(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		office_id,
		phone_num,
		addr_id
		FROM offices`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var offices []Office
	for rows.Next() {
		var (
			office_id int
			phone_num string
			addr_id   int
		)
		if err := rows.Scan(
			&office_id,
			&phone_num,
			&addr_id); err != nil {
			log.Fatal(err)
		}
		offices = append(offices, Office{
			Office_id: office_id,
			Phone_num: phone_num,
			Addr_id:   addr_id,
		})
	}
	json.NewEncoder(res).Encode(offices)
}

type OfficeRequest struct {
	Phone_num string `json:"phone_num"`
	Addr_id   int    `json:"addr_id"`
}

func PostOffice(res http.ResponseWriter, req *http.Request) {
	var office *OfficeRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&office)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
	}

	if office.Phone_num == "" {
		ErrorRes(res, http.StatusBadRequest,
			"phone_num is required")
		return
	}

	if len(office.Phone_num) != 11 {
		ErrorRes(res, http.StatusBadRequest,
			"phone_num must be 11 characters")
		return
	}

	if office.Addr_id == 0 {
		ErrorRes(res, http.StatusBadRequest, 
			"addr_id is required")
		return
	}

	rows, err := AddOffice(office)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest, 
			fmt.Sprintf("Error adding office, make sure addr_id is also in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Office added successfully. %d rows affected", rows),
	})
}
