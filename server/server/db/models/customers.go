package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateCustomersTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS customers (
	cust_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	cust_fname TEXT NOT NULL,
	cust_lname TEXT NOT NULL,
	cust_email TEXT NOT NULL UNIQUE,
	phone_num TEXT,
	addr_id INT REFERENCES addresses(addr_id),
	cred_limit NUMERIC(10,2) NOT NULL,
	cust_pass TEXT NOT NULL
	)`)
	if err != nil {
		panic(err)
	}
}
