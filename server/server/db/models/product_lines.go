package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateProductLinesTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS product_lines (
	prod_line_name TEXT PRIMARY KEY,
	prod_line_desc TEXT
	)`)
	if err != nil {
		panic(err)
	}
}
