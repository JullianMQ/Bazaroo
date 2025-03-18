package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/JullianMQ/Bazaroo/server/db"
	"github.com/JullianMQ/Bazaroo/server/utils"
)

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
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.GetAddrQuery()
	if err != nil {
		utils.ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
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
			log.Println(err)
			return
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
	json.NewEncoder(res).Encode(addr)
}


