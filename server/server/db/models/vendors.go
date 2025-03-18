package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateVendorsTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS vendors (
	vendor_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	vendor_name TEXT NOT NULL,
	vendor_email TEXT NOT NULL UNIQUE,
	vendor_phone_num TEXT NOT NULL,
	addr_id INT REFERENCES addresses(addr_id)
	)`)
	if err != nil {
		panic(err)
	}
}
