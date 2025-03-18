package utils

import (
	"database/sql"
	"log"
	"net/mail"
	"regexp"
	"slices"

	"github.com/JullianMQ/Bazaroo/server/db"
)

func ContainsEmpty(s []string) bool {
	return slices.Contains(s, "")
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

func IsEmailValid(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

func IsEmailInDb(email string, ut string) bool {
	var rows *sql.Rows
	var err error

	switch ut {
	case "customer":
		rows, err = db.GetCustByEmailQuery(email)
	case "employee":
		rows, err = db.GetEmpByEmailQuery(email)
	case "vendor":
		rows, err = db.GetVendorByEmailQuery(email)
	default:
		panic("Invalid user type")
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user_email string
		if err := rows.Scan(&user_email); err != nil {
			log.Println(err)
		}
		if user_email == email {
			return true
		}
	}
	return false
}

func IsAddrIdInDb(addr_id int) bool {
	rows, err := db.GetAddrByIdQuery(addr_id)
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

func IsProdLineInDb(prod_line string) bool {
	rows, err := db.GetProdLineByNameQuery(prod_line)
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

func IsPhoneValid(phone_num string) bool {
	return regexp.MustCompile(`^09\d{9}$`).MatchString(phone_num)
}

func CustomerIdInDb(cust_id int) bool {
	rows, err := db.GetCustomerIdInDbQuery(cust_id)
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
