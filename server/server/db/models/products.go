package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateProductsTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS products (
	prod_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	prod_name TEXT NOT NULL,
	prod_line_name TEXT NOT NULL REFERENCES product_lines(prod_line_name),
	prod_vendor_id INT REFERENCES vendors(vendor_id),
	office_id INT NOT NULL REFERENCES offices(office_id),
	prod_desc TEXT,
	prod_image TEXT,
	quan_in_stock INT,
	buy_price NUMERIC(10,2),
	msrp NUMERIC(10,2)
	)`)
	if err != nil {
		panic(err)
	}
}
