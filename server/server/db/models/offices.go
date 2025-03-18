package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateOfficesTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS offices (
	office_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	phone_num TEXT NOT NULL,
	addr_id INT NOT NULL REFERENCES addresses(addr_id)
	)`)
	if err != nil {
		panic(err)
	}
}
