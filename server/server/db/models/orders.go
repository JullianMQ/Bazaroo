package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateOrdersTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS orders (
	ord_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	cust_id INT NOT NULL REFERENCES customers(cust_id),
	ord_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	req_shipped_date DATE NOT NULL,
	comments TEXT,
	rating INT
	)`)
	if err != nil {
		panic(err)
	}
}
