package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetRoot(res http.ResponseWriter, req *http.Request) {
	myMess := "Message"
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(myMess)
}

func GetEmps(res http.ResponseWriter, req *http.Request) {
	mess := "emps"
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(mess)
}

type MyMess struct {
	Message string `json:"Great"`
}

// TODO: GET EMP BY ID
func GetEmpId(res http.ResponseWriter, req *http.Request) {
	mess := req.URL.Query().Get("id")
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(MyMess{
		Message: mess,
	})
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
	Addr_id     int    `json:"addr_id"`
	Addr_line1  string `json:"addr_line1"`
	Addr_line2  sql.NullString `json:"addr_line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	Postal_code string `json:"postal_code"`
	Country     string `json:"country"`
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

type ErrorResponse struct {
	Error string `json:"error"`
}

type OkResponse struct {
	Message string `json:"message"`
}

func PostAddr(res http.ResponseWriter, req *http.Request) {
	var addr *AddrRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&addr)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "Could not decode request body",
		})
		panic(err)
	}

	if addr.Addr_line1 == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "addr_line1 is required",
		})
		return
	}

	if addr.City == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "city is required",
		})
		return
	}

	if addr.State == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "state is required",
		})
		return
	}

	if addr.Postal_code == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "postal_code is required",
		})
		return
	}

	if len(addr.Postal_code) > 10 {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "postal_code must be 10 characters or less",
		})
		return
	}

	if addr.Country == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "country is required",
		})
		return
	}

	if len(addr.Country) != 3 {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "country must be 3 characters, of iso 3166-1 alpha-3 code",
		})
		return
	}

	rows, err := AddAddr(addr)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: fmt.Sprintf("Length of country is %d", len(addr.Country)),
		})
		log.Fatal(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Address added successfully. %d rows affected", rows),
	})
}
