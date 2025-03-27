package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"slices"
	"regexp"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type OkResponse struct {
	Message string `json:"message"`
}

func ContainsEmpty(s []string) bool {
	return slices.Contains(s, "")
}

func isPhoneValid(phone_num string) bool {
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(phone_num)
}

func ContainsZero(v any) bool {
	switch values := v.(type) {
	case []int:
		return containsInt(values, 0)
	case []float64:
		return containsFloat(values, 0.0)
	case []any:
		for _, value := range values {
			if value == 0 {
				return true
			}
		}
		for _, value := range values {
			if value == 0.0 {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func containsInt(v []int, target int) bool {
	return slices.Contains(v, target)
}

func containsFloat(v []float64, target float64) bool {
	return slices.Contains(v, target)
}

func ErrorRes(res http.ResponseWriter, status int, mess string) {
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(ErrorResponse{
		Error: mess,
	})
}

func isEmailValid(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}


// WITH QUERIES
func isAddrIdInDb(addr_id int) bool {
	rows, err := db.Query(`SELECT
		addr_id
		FROM addresses
		WHERE addr_id = $1`, addr_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var addr_id int
		if err := rows.Scan(&addr_id); err != nil {
			log.Println(err)
		}
		if addr_id == addr_id {
			return true
		}
	}
	return false
}

func isEmailInDb(email string, ut string) bool {
	var rows *sql.Rows
	var err error

	switch ut {
	case "customer":
		rows, err = db.Query(`SELECT
			cust_email
			FROM customers
			WHERE cust_email = $1`, email)
	case "employee":
		rows, err = db.Query(`SELECT
			emp_email
			FROM employees
			WHERE emp_email = $1`, email)
	case "vendor":
		rows, err = db.Query(`SELECT
			vendor_email
			FROM vendors
			WHERE vendor_email = $1`, email)
	default:
		panic("Invalid user type")
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var emp_email string
		if err := rows.Scan(&emp_email); err != nil {
			log.Println(err)
		}
		if emp_email == email {
			return true
		}
	}
	return false
}

func isEmpIdInDb(emp_id int) bool {
	rows, err := db.Query(`SELECT
		emp_id
		FROM employees
		WHERE emp_id = $1`, emp_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var emp_id int
		if err := rows.Scan(&emp_id); err != nil {
			log.Println(err)
		}
		if emp_id == emp_id {
			return true
		}
	}
	return false
}

func isProdLineInDb(prod_line string) bool {
	rows, err := db.Query(`SELECT
		prod_line_name
		FROM product_lines
		WHERE prod_line_name = $1`, prod_line)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var prod_line_name string
		if err := rows.Scan(&prod_line_name); err != nil {
			log.Println(err)
		}
		if prod_line_name == prod_line {
			return true
		}
	}
	return false
}

func isProdIdInDb(prod_id int) bool {
	rows, err := db.Query(`SELECT
		prod_id
		FROM products
		WHERE prod_id = $1`, prod_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var prod_id int
		if err := rows.Scan(&prod_id); err != nil {
			log.Println(err)
		}
		if prod_id == prod_id {
			return true
		}
	}
	return false
}

func CustomerIdInDb(cust_id int) bool {
	rows, err := db.Query(`SELECT
		cust_id
		FROM customers
		WHERE cust_id = $1`, cust_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cust_id int
		if err := rows.Scan(&cust_id); err != nil {
			log.Println(err)
		}
		if cust_id == cust_id {
			return true
		}
	}
	return false
}

func isOrderIdInDb(order_id int) bool {
	rows, err := db.Query(`SELECT
		ord_id
		FROM orders
		WHERE ord_id = $1`, order_id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ord_id int
		if err := rows.Scan(&ord_id); err != nil {
			log.Println(err)
		}
		if ord_id == order_id {
			return true
		}
	}
	return false
}
