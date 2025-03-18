package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateAddressesTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS addresses(
	addr_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	addr_line1 TEXT NOT NULL,
	addr_line2 TEXT,
	city TEXT NOT NULL,
	state TEXT NOT NULL,
	postal_code VARCHAR(10) NOT NULL,
	country CHAR(3) NOT NULL
	)`)
	if err != nil {
		panic(err)
	}
}
