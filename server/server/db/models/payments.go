package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreatePaymentsTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS payments (
	payment_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	cust_id INT NOT NULL REFERENCES customers(cust_id),
	status TEXT DEFAULT 'in cart',
	payment_date TIMESTAMP,
	amount NUMERIC(10,2) NOT NULL,
	payment_status TEXT,
	ord_id INT NOT NULL REFERENCES orders(ord_id)
	)`)
	if err != nil {
		panic(err)
	}
}
