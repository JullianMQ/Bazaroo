package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateOrderDetailsTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS order_details (
		ord_id INT,
		prod_id INT,
		PRIMARY KEY (ord_id, prod_id),
		CONSTRAINT ord_id_fk FOREIGN KEY (ord_id) REFERENCES orders(ord_id),
		CONSTRAINT prod_id_fk FOREIGN KEY (prod_id) REFERENCES products(prod_id),
		quan_ordered INT NOT NULL,
		price_each NUMERIC(10,2) NOT NULL
	)`)
	if err != nil {
		panic(err)
	}
}
