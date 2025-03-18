package db

import (
	"database/sql"
)

func GetAddrQuery() (*sql.Rows, error) {
	return DB.Query(`SELECT
		addr_id,
		addr_line1,
		addr_line2,
		city,
		state,
		postal_code,
		country
		FROM addresses`)
}

func GetAddrByIdQuery(addr_id int) (*sql.Rows, error) {
	return DB.Query(`SELECT
		addr_id,
		FROM addresses
		WHERE addr_id = $1`, addr_id)
}

func GetCustByEmailQuery(email string) (*sql.Rows, error) {
	return DB.Query(`SELECT
		cust_email,
		FROM customers
		WHERE cust_email = $1`, email)
}

func GetEmpByEmailQuery(email string) (*sql.Rows, error) {
	return DB.Query(`SELECT
		emp_email,
		FROM employees
		WHERE emp_email = $1`, email)
}

func GetVendorByEmailQuery(email string) (*sql.Rows, error) {
	return DB.Query(`SELECT
		vendor_email,
		FROM vendors
		WHERE vendor_email = $1`, email)
}

func GetProdLineByNameQuery(prod_line string) (*sql.Rows, error) {
	return DB.Query(`SELECT
		prod_line_name
		FROM product_lines
		WHERE prod_line_name = $1`, prod_line)
}

func GetCustomerIdInDbQuery(cust_id int) (*sql.Rows, error) {
	return DB.Query(`SELECT
		cust_id
		FROM customers
		WHERE cust_id = $1`, cust_id)
}
