package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"slices"
	"strconv"
	"time"
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

func GetRoot(res http.ResponseWriter, req *http.Request) {
	myMess := "Message"
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(myMess)
}

func isEmailValid(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

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

func isPhoneValid(phone_num string) bool {
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(phone_num)
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

type AddrRequest struct {
	Addr_line1  string `json:"addr_line1"`
	Addr_line2  string `json:"addr_line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	Postal_code string `json:"postal_code"`
	Country     string `json:"country"`
}

type AddrResponse struct {
	Addr_id     int            `json:"addr_id"`
	Addr_line1  string         `json:"addr_line1"`
	Addr_line2  sql.NullString `json:"addr_line2"`
	City        string         `json:"city"`
	State       string         `json:"state"`
	Postal_code string         `json:"postal_code"`
	Country     string         `json:"country"`
}

func GetAddr(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		addr_id,
		addr_line1,
		addr_line2,
		city,
		state,
		postal_code,
		country
		FROM addresses`)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
	}
	defer rows.Close()

	var addr []AddrResponse
	for rows.Next() {
		var (
			addr_id     int
			addr_line1  string
			addr_line2  sql.NullString
			city        string
			state       string
			postal_code string
			country     string
		)
		if err := rows.Scan(
			&addr_id,
			&addr_line1,
			&addr_line2,
			&city,
			&state,
			&postal_code,
			&country); err != nil {
			log.Println(err)
		}
		addr = append(addr, AddrResponse{
			Addr_id:     addr_id,
			Addr_line1:  addr_line1,
			Addr_line2:  addr_line2,
			City:        city,
			State:       state,
			Postal_code: postal_code,
			Country:     country,
		})
	}
	json.NewEncoder(res).Encode(addr)
}

func PostAddr(res http.ResponseWriter, req *http.Request) {
	var addr *AddrRequest
	res.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(req.Body).Decode(&addr)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		return
	}

	if ContainsEmpty([]string{
		addr.Addr_line1,
		addr.City,
		addr.State,
		addr.Postal_code,
		addr.Country}) {
		ErrorRes(res, http.StatusBadRequest,
			"addr_line1, city, state, postal_code, country cannot be empty")
		return
	}

	if len(addr.Postal_code) > 10 {
		ErrorRes(res, http.StatusBadRequest,
			"postal_code must be 10 characters or less")
		return
	}

	if len(addr.Country) != 3 {
		ErrorRes(res, http.StatusBadRequest,
			"country must be 3 characters, of iso 3166-1 alpha-3 code")
		return
	}

	rows, err := AddAddr(addr)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not add address, try again later."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Address added successfully. %d rows affected", rows),
	})
}

type Office struct {
	Office_id int    `json:"office_id"`
	Phone_num string `json:"phone_num"`
	Addr_id   int    `json:"addr_id"`
}

func GetOffices(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		office_id,
		phone_num,
		addr_id
		FROM offices`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var offices []Office
	for rows.Next() {
		var (
			office_id int
			phone_num string
			addr_id   int
		)
		if err := rows.Scan(
			&office_id,
			&phone_num,
			&addr_id); err != nil {
			ErrorRes(res, http.StatusInternalServerError,
				fmt.Sprintf("Could not get offices"))
			log.Println(err)
			return
		}
		offices = append(offices, Office{
			Office_id: office_id,
			Phone_num: phone_num,
			Addr_id:   addr_id,
		})
	}
	json.NewEncoder(res).Encode(offices)
}

type OfficeRequest struct {
	Phone_num string `json:"phone_num"`
	Addr_id   int    `json:"addr_id"`
}

func PostOffice(res http.ResponseWriter, req *http.Request) {
	var office *OfficeRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&office)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
	}

	if ContainsEmpty([]string{
		office.Phone_num}) {
		ErrorRes(res, http.StatusBadRequest,
			"phone_num cannot be empty")
		return
	}

	if ContainsZero([]int{
		office.Addr_id}) {
		ErrorRes(res, http.StatusBadRequest,
			"addr_id is required")
		return
	}

	if len(office.Phone_num) != 11 {
		ErrorRes(res, http.StatusBadRequest,
			"phone_num must be 11 characters")
		return
	}

	rows, err := AddOffice(office)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding office, make sure addr_id is also in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Office added successfully. %d rows affected", rows),
	})
}

type Employee struct {
	Emp_id    int           `json:"emp_id"`
	Emp_fname string        `json:"emp_fname"`
	Emp_lname string        `json:"emp_lname"`
	Emp_email string        `json:"emp_email"`
	Office_id sql.NullInt64 `json:"office_id"`
	Job_title string        `json:"job_title"`
}

func GetEmps(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		emp_id,
		emp_fname,
		emp_lname,
		emp_email,
		office_id,
		job_title
		FROM employees`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var emps []Employee
	for rows.Next() {
		var (
			emp_id    int
			emp_fname string
			emp_lname string
			emp_email string
			office_id sql.NullInt64
			job_title string
		)
		if err := rows.Scan(
			&emp_id,
			&emp_fname,
			&emp_lname,
			&emp_email,
			&office_id,
			&job_title); err != nil {
			ErrorRes(res, http.StatusInternalServerError,
				fmt.Sprintf("Could not get employees"))
			log.Println(err)
			return
		}
		emps = append(emps, Employee{
			Emp_id:    emp_id,
			Emp_fname: emp_fname,
			Emp_lname: emp_lname,
			Emp_email: emp_email,
			Office_id: office_id,
			Job_title: job_title,
		})
	}
	json.NewEncoder(res).Encode(emps)
}

func GetEmpId(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	emp_id := req.URL.Query().Get("id")
	emp_id_int, err := strconv.ParseInt(emp_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id into int64, try again later."))
		log.Println(err)
		return
	}

	emp, err := GetEmpById(emp_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not get employee by id, check if id is correct."))
		log.Println(err)
		return
	}

	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("%v", emp.Emp_fname),
	})
}

type EmployeeRequest struct {
	Emp_fname string `json:"emp_fname"`
	Emp_lname string `json:"emp_lname"`
	Emp_email string `json:"emp_email"`
	Office_id int    `json:"office_id"`
	Job_title string `json:"job_title"`
}

func PostEmp(res http.ResponseWriter, req *http.Request) {
	var emp *EmployeeRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&emp)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		emp.Emp_fname,
		emp.Emp_lname,
		emp.Emp_email,
		emp.Job_title}) {
		ErrorRes(res, http.StatusBadRequest,
			"emp_fname, emp_lname, emp_email, job_title cannot be empty")
		return
	}

	if ContainsZero([]int{
		emp.Office_id}) {
		ErrorRes(res, http.StatusBadRequest,
			"office_id is required")
		return
	}

	if !isEmailValid(emp.Emp_email) {
		ErrorRes(res, http.StatusBadRequest,
			"email is invalid")
		return
	}

	if isEmailInDb(emp.Emp_email, "employee") {
		ErrorRes(res, http.StatusBadRequest,
			"email is already in use")
		return
	}

	rows, err := AddEmp(emp)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("CHECK IF OFFICE ID IN DATABASE: %s", err))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Employee added successfully. %d rows affected", rows),
	})
}

type Vendor struct {
	Vendor_id        int    `json:"vendor_id"`
	Vendor_name      string `json:"vendor_name"`
	Vendor_email     string `json:"vendor_email"`
	Vendor_phone_num string `json:"vendor_phone_num"`
	Addr_id          int    `json:"addr_id"`
}

func GetVendors(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		vendor_id,
		vendor_name,
		vendor_email,
		vendor_phone_num,
		addr_id
		FROM vendors`)
	if err != nil {
		log.Println(fmt.Sprintf("CHECK IF VENDOR ID IN DATABASE: %s", err))
		log.Println(err)
		return
	}
	defer rows.Close()

	var vendors []Vendor

	for rows.Next() {
		var vendor_id int
		var vendor_name string
		var vendor_email string
		var vendor_phone_num string
		var addr_id int
		err := rows.Scan(&vendor_id, &vendor_name, &vendor_email, &vendor_phone_num, &addr_id)
		if err != nil {
			log.Println(fmt.Sprintf("CHECK IF VENDOR ID IN DATABASE: %s", err))
			log.Println(err)
			return
		}
		vendor := Vendor{
			Vendor_id:        vendor_id,
			Vendor_name:      vendor_name,
			Vendor_email:     vendor_email,
			Vendor_phone_num: vendor_phone_num,
			Addr_id:          addr_id,
		}
		vendors = append(vendors, vendor)
	}
	json.NewEncoder(res).Encode(vendors)
}

func GetVendorId(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vendor_id := req.URL.Query().Get("id")
	vendor_id_int, err := strconv.ParseInt(vendor_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			"Invalid vendor id")
		return
	}

	vendor, err := GetVendorById(vendor_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Error getting vendor")
		log.Println(err)
		return
	}

	json.NewEncoder(res).Encode(vendor)
}

type VendorRequest struct {
	Vendor_name      string `json:"vendor_name"`
	Vendor_email     string `json:"vendor_email"`
	Vendor_phone_num string `json:"vendor_phone_num"`
	Addr_id          int    `json:"addr_id"`
}

func PostVendor(res http.ResponseWriter, req *http.Request) {
	var vendor *VendorRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&vendor)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Error decoding request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		vendor.Vendor_name,
		vendor.Vendor_email,
		vendor.Vendor_phone_num}) {
		ErrorRes(res, http.StatusBadRequest,
			"vendor_name, vendor_email, vendor_phone_num cannot be empty")
		return
	}

	if ContainsZero([]int{
		vendor.Addr_id}) {
		ErrorRes(res, http.StatusBadRequest,
			"addr_id is required")
		return
	}

	if !isEmailValid(vendor.Vendor_email) {
		ErrorRes(res, http.StatusBadRequest,
			"vendor_email is invalid")
		return
	}

	if isEmailInDb(vendor.Vendor_email, "vendor") {
		ErrorRes(res, http.StatusBadRequest,
			"vendor_email is already in use")
		return
	}

	if !isPhoneValid(vendor.Vendor_phone_num) {
		ErrorRes(res, http.StatusBadRequest,
			"Vendor phone number is invalid")
		return
	}

	rows, err := AddVendor(vendor)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding vendor, make sure addr_id is also in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Vendor added successfully. %d rows affected", rows),
	})
}

type ProductLine struct {
	Prod_line_name string         `json:"prod_line_name"`
	Prod_line_desc sql.NullString `json:"prod_line_desc"`
}

func GetProductLine(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		prod_line_name,
		prod_line_desc
		FROM product_lines`)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var productLines []ProductLine
	for rows.Next() {
		var (
			prod_line_name string
			prod_line_desc sql.NullString
		)
		if err := rows.Scan(
			&prod_line_name,
			&prod_line_desc); err != nil {
			log.Println(err)
			return
		}
		productLines = append(productLines, ProductLine{
			Prod_line_name: prod_line_name,
			Prod_line_desc: prod_line_desc,
		})
	}
	json.NewEncoder(res).Encode(productLines)
}

type ProductLineRequest struct {
	Prod_line_name string `json:"prod_line_name"`
	Prod_line_desc string `json:"prod_line_desc"`
}

func PostProductLine(res http.ResponseWriter, req *http.Request) {
	var productLine *ProductLineRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&productLine)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		productLine.Prod_line_name}) {
		ErrorRes(res, http.StatusBadRequest,

			"prod_line_name, prod_line_desc cannot be empty")
		return
	}

	rows, err := AddProductLine(productLine)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding product line, try again later"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Product line added successfully. %d rows affected", rows),
	})
}

type Product struct {
	Prod_id        int            `json:"prod_id"`
	Prod_name      string         `json:"prod_name"`
	Prod_line_name string         `json:"prod_line_name"`
	Prod_vendor_id int            `json:"prod_vendor_id"`
	Prod_desc      sql.NullString `json:"prod_desc"`
	Prod_image     sql.NullString `json:"prod_image"`
	Quan_in_stock  int            `json:"quan_in_stock"`
	Buy_price      float64        `json:"buy_price"`
	Msrp           float64        `json:"msrp"`
}

func GetProducts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		prod_id,
		prod_name,
		prod_line_name,
		prod_vendor_id,
		prod_desc,
		prod_image,
		quan_in_stock,
		buy_price,
		msrp
		FROM products`)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var (
			prod_id        int
			prod_name      string
			prod_line_name string
			prod_vendor_id int
			prod_desc      sql.NullString
			prod_image     sql.NullString
			quan_in_stock  int
			buy_price      float64
			msrp           float64
		)
		if err := rows.Scan(
			&prod_id,
			&prod_name,
			&prod_line_name,
			&prod_vendor_id,
			&prod_desc,
			&prod_image,
			&quan_in_stock,
			&buy_price,
			&msrp); err != nil {
			log.Println(err)
			return
		}
		products = append(products, Product{
			Prod_id:        prod_id,
			Prod_name:      prod_name,
			Prod_line_name: prod_line_name,
			Prod_vendor_id: prod_vendor_id,
			Prod_desc:      prod_desc,
			Prod_image:     prod_image,
			Quan_in_stock:  quan_in_stock,
			Buy_price:      buy_price,
			Msrp:           msrp,
		})
	}
	json.NewEncoder(res).Encode(products)
}

type ProductRequest struct {
	Prod_name      string  `json:"prod_name"`
	Prod_line_name string  `json:"prod_line_name"`
	Prod_vendor_id int     `json:"prod_vendor_id"`
	Prod_desc      string  `json:"prod_desc"`
	Prod_image     string  `json:"prod_image"`
	Quan_in_stock  int     `json:"quan_in_stock"`
	Buy_price      float64 `json:"buy_price"`
	Msrp           float64 `json:"msrp"`
}

func PostProduct(res http.ResponseWriter, req *http.Request) {
	var product *ProductRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		product.Prod_name,
		product.Prod_line_name}) {
		ErrorRes(res, http.StatusBadRequest,
			"prod_name, prod_line_name cannot be empty")
		return
	}

	if ContainsZero([]any{
		product.Quan_in_stock,
		product.Buy_price,
		product.Msrp}) {
		ErrorRes(res, http.StatusBadRequest,
			"quan_in_stock, buy_price, msrp cannot be zero")
		return
	}

	if !isProdLineInDb(product.Prod_line_name) {
		ErrorRes(res, http.StatusBadRequest,
			"prod_line_name not in database")
		return
	}

	rows, err := AddProduct(product)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding product, make sure prod_vendor_id is also in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Product added successfully. %d rows affected", rows),
	})
}

type Customer struct {
	Cust_id    int           `json:"cust_id"`
	Cust_fname string        `json:"cust_fname"`
	Cust_lname string        `json:"cust_lname"`
	Cust_email string        `json:"cust_email"`
	Phone_num  string        `json:"phone_num"`
	Addr_id    sql.NullInt64 `json:"addr_id"`
	Cred_limit float64       `json:"cred_limit"`
}

func GetCustomers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		cust_id,
		cust_fname,
		cust_lname,
		cust_email,
		phone_num,
		addr_id,
		cred_limit
		FROM customers`)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not get customers, try again later")
		log.Println(err)
		return
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var (
			cust_id    int
			cust_fname string
			cust_lname string
			cust_email string
			phone_num  string
			addr_id    sql.NullInt64
			cred_limit float64
		)
		if err := rows.Scan(
			&cust_id,
			&cust_fname,
			&cust_lname,
			&cust_email,
			&phone_num,
			&addr_id,
			&cred_limit); err != nil {
			log.Println(err)
			return
		}
		customers = append(customers, Customer{
			Cust_id:    cust_id,
			Cust_fname: cust_fname,
			Cust_lname: cust_lname,
			Cust_email: cust_email,
			Phone_num:  phone_num,
			Addr_id:    addr_id,
			Cred_limit: cred_limit,
		})
	}
	json.NewEncoder(res).Encode(customers)
}

type CustomerAcc struct {
	Cust_id    int           `json:"cust_id"`
	Cust_fname string        `json:"cust_fname"`
	Cust_lname string        `json:"cust_lname"`
	Cust_email string        `json:"cust_email"`
	Phone_num  string        `json:"phone_num"`
}

func GetCustomerById(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.ParseInt(cust_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, try again later."))
		log.Println(err)
		return
	}

	cust, err := GetCustById(cust_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not get customer by id, check if id is correct."))
		log.Println(err)
		return
	}

	json.NewEncoder(res).Encode(CustomerAcc{
		Cust_id:    cust.Cust_id,
		Cust_fname: cust.Cust_fname,
		Cust_lname: cust.Cust_lname,
		Cust_email: cust.Cust_email,
		Phone_num:  cust.Phone_num,
	})
}

type CustomerRequest struct {
	Cust_fname       string  `json:"cust_fname"`
	Cust_lname       string  `json:"cust_lname"`
	Cust_email       string  `json:"cust_email"`
	Phone_num        string  `json:"phone_num"`
	Addr_id          int     `json:"addr_id"`
	Sales_rep_emp_id int     `json:"sales_rep_emp_id"`
	Cred_limit       float64 `json:"cred_limit"`
}

func PostCustomer(res http.ResponseWriter, req *http.Request) {
	var customer *CustomerRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&customer)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		customer.Cust_fname,
		customer.Cust_lname,
		customer.Cust_email}) {
		ErrorRes(res, http.StatusBadRequest,
			"cust_fname, cust_lname, cust_email cannot be empty")
		return
	}

	if ContainsZero([]any{
		customer.Cred_limit}) {
		ErrorRes(res, http.StatusBadRequest,
			"cred_limit cannot be zero")
		return
	}

	if !isEmailValid(customer.Cust_email) {
		ErrorRes(res, http.StatusBadRequest,
			"cust_email is invalid")
		return
	}

	if isEmailInDb(customer.Cust_email, "customer") {
		ErrorRes(res, http.StatusBadRequest,
			"cust_email is already in use")
		return
	}

	if customer.Addr_id != 0 {
		if !isAddrIdInDb(customer.Addr_id) {
			ErrorRes(res, http.StatusBadRequest,
				"addr_id is invalid")
			return
		}
	}

	if customer.Phone_num != "" {
		if !isPhoneValid(customer.Phone_num) {
			ErrorRes(res, http.StatusBadRequest,
				"phone_num is invalid")
			return
		}
	}

	if !isEmpIdInDb(customer.Sales_rep_emp_id) {
		ErrorRes(res, http.StatusBadRequest,
			"sales_rep_emp_id is invalid")
		return
	}

	rows, err := AddCustomer(customer)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding customer, make sure addr_id is also in database"))
		log.Println(err)
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Customer added successfully. %d rows affected", rows),
	})
}

type Order struct {
	Ord_id           int            `json:"ord_id"`
	Cust_id          int            `json:"cust_id"`
	Ord_date         time.Time      `json:"ord_date"`
	Req_shipped_date time.Time      `json:"req_shipped_date"`
	Comments         sql.NullString `json:"comments"`
	Rating           int            `json:"rating"`
}

type OrderByCustId struct {
	Ord_id           int       `json:"ord_id"`
	Ord_date         time.Time `json:"ord_date"`
	Req_shipped_date time.Time `json:"req_shipped_date"`
	Comments         string    `json:"comments"`
	Rating           int       `json:"rating"`
}

func GetOrders(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		ord_id,
		cust_id,
		ord_date,
		req_shipped_date,
		comments,
		rating
		FROM orders`)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not get orders, try again later")
		log.Println(err)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var (
			ord_id           int
			cust_id          int
			ord_date         time.Time
			req_shipped_date time.Time
			comments         sql.NullString
			rating           int
		)
		if err := rows.Scan(
			&ord_id,
			&cust_id,
			&ord_date,
			&req_shipped_date,
			&comments,
			&rating); err != nil {
			log.Println(err)
			return
		}
		orders = append(orders, Order{
			Ord_id:   ord_id,
			Cust_id:  cust_id,
			Ord_date: ord_date,
			// take out the hour, minute, second, and nanosecond
			Req_shipped_date: req_shipped_date.Truncate(24 * time.Hour),
			Comments:         comments,
			Rating:           rating,
		})
	}
	json.NewEncoder(res).Encode(orders)
}

func GetOrderByCustId(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.ParseInt(cust_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id into int64, try again later."))
		log.Println(err)
		return
	}

	rows, err := db.Query(`SELECT
		ord_id,
		cust_id,
		ord_date,
		req_shipped_date,
		comments,
		rating
		FROM orders
		WHERE cust_id = $1`, cust_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not get orders, try again later")
		log.Println(err)
		return
	}
	defer rows.Close()

	var orders []OrderByCustId
	for rows.Next() {
		var (
			ord_id           int
			cust_id          int
			ord_date         time.Time
			req_shipped_date time.Time
			comments         string
			rating           int
		)
		if err := rows.Scan(
			&ord_id,
			&cust_id,
			&ord_date,
			&req_shipped_date,
			&comments,
			&rating); err != nil {
			log.Println(err)
			return
		}
		orders = append(orders, OrderByCustId{
			Ord_id:   ord_id,
			Ord_date: ord_date,
			// take out the hour, minute, second, and nanosecond
			Req_shipped_date: req_shipped_date.Truncate(24 * time.Hour),
			Comments:         comments,
			Rating:           rating,
		})
	}
	json.NewEncoder(res).Encode(orders)
}

type OrderRequest struct {
	Cust_id          int       `json:"cust_id"`
	Ord_date         time.Time `json:"ord_date"`
	Req_shipped_date time.Time `json:"req_shipped_date"`
	Comments         string    `json:"comments"`
	Rating           int       `json:"rating"`
}

func PostOrder(res http.ResponseWriter, req *http.Request) {
	var order *OrderRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&order)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if !CustomerIdInDb(order.Cust_id) {
		ErrorRes(res, http.StatusBadRequest,
			"cust_id is not in db")
		return
	}

	if order.Ord_date.IsZero() {
		ErrorRes(res, http.StatusBadRequest,
			"ord_date cannot be zero")
		return
	}

	if order.Req_shipped_date.IsZero() {
		ErrorRes(res, http.StatusBadRequest,
			"req_shipped_date cannot be zero")
		return
	}

	if ContainsZero([]any{
		order.Cust_id}) {
		ErrorRes(res, http.StatusBadRequest,
			"rating cannot be zero")
		return
	}

	if order.Rating < 1 || order.Rating > 5 {
		ErrorRes(res, http.StatusBadRequest,
			"rating must be between 1 and 5")
		return
	}

	rows, err := AddOrder(order)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding order, make sure cust_id is also in database"))
		log.Println(err)
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Order added successfully. %d rows affected", rows),
	})
}

type Payment struct {
	Payment_id     int       `json:"payment_id"`
	Cust_id        int       `json:"cust_id"`
	Payment_date   time.Time `json:"payment_date"`
	Amount         float64   `json:"amount"`
	Payment_status string    `json:"payment_status"`
	Ord_id         int       `json:"ord_id"`
}

func GetPaymentsByCustId(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.ParseInt(cust_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, make sure that it is included in the query."))
		log.Println(err)
		return
	}

	rows, err := db.Query(`SELECT
		payment_id,
		cust_id,
		payment_date,
		amount,
		payment_status,
		ord_id
		FROM payments
		WHERE cust_id = $1`, cust_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not get payments, try again later")
		log.Println(err)
		return
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var (
			payment_id     int
			cust_id        int
			payment_date   time.Time
			amount         float64
			payment_status string
			ord_id         int
		)
		if err := rows.Scan(
			&payment_id,
			&cust_id,
			&payment_date,
			&amount,
			&payment_status,
			&ord_id); err != nil {
			log.Println(err)
			return
		}
		payments = append(payments, Payment{
			Payment_id:     payment_id,
			Cust_id:        cust_id,
			Payment_date:   payment_date,
			Amount:         amount,
			Payment_status: payment_status,
			Ord_id:         ord_id,
		})
	}
	json.NewEncoder(res).Encode(payments)
}

type PaymentRequest struct {
	Cust_id        int       `json:"cust_id"`
	Payment_date   time.Time `json:"payment_date"`
	Amount         float64   `json:"amount"`
	Payment_status string    `json:"payment_status"`
	Ord_id         int       `json:"ord_id"`
}

func PostPayment(res http.ResponseWriter, req *http.Request) {
	var payment *PaymentRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&payment)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		payment.Payment_status,
	}) {
		ErrorRes(res, http.StatusBadRequest,
			"payment_status cannot be empty")
		return
	}

	if ContainsZero([]any{
		payment.Cust_id,
		payment.Amount,
		payment.Ord_id}) {
		ErrorRes(res, http.StatusBadRequest,
			"cust_id is required")
		return
	}

	if payment.Payment_date.IsZero() {
		ErrorRes(res, http.StatusBadRequest,
			"payment_date cannot be zero")
		return
	}

	if payment.Amount <= 0 {
		ErrorRes(res, http.StatusBadRequest,
			"amount must be greater than zero")
		return
	}

	if !CustomerIdInDb(payment.Cust_id) {
		ErrorRes(res, http.StatusBadRequest,
			"cust_id is invalid")
		return
	}

	if !isOrderIdInDb(payment.Ord_id) {
		ErrorRes(res, http.StatusBadRequest,
			"ord_id is invalid")
		return
	}

	rows, err := AddPayment(payment)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding payment, make sure cust_id is also in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Payment added successfully. %d rows affected", rows),
	})
}

type OrderDetail struct {
	Ord_id       int     `json:"ord_id"`
	Prod_id      int     `json:"prod_id"`
	Quan_ordered int     `json:"quan_ordered"`
	Price_each   float64 `json:"price_each"`
}

func GetOrderDetailsByOrderId(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	order_id := req.URL.Query().Get("id")
	order_id_int, err := strconv.ParseInt(order_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, make sure that it is included in the query."))
		log.Println(err)
		return
	}

	rows, err := db.Query(`SELECT
		ord_id,
		prod_id,
		quan_ordered,
		price_each
		FROM order_details
		WHERE ord_id = $1`, order_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not get order details, try again later")
		log.Println(err)
		return
	}
	defer rows.Close()

	var orderDetails []OrderDetail
	for rows.Next() {
		var (
			ord_id       int
			prod_id      int
			quan_ordered int
			price_each   float64
		)
		if err := rows.Scan(
			&ord_id,
			&prod_id,
			&quan_ordered,
			&price_each); err != nil {
			log.Println(err)
			return
		}
		orderDetails = append(orderDetails, OrderDetail{
			Ord_id:       ord_id,
			Prod_id:      prod_id,
			Quan_ordered: quan_ordered,
			Price_each:   price_each,
		})
	}
	json.NewEncoder(res).Encode(orderDetails)
}

type OrderDetailRequest struct {
	Ord_id       int     `json:"ord_id"`
	Prod_id      int     `json:"prod_id"`
	Quan_ordered int     `json:"quan_ordered"`
	Price_each   float64 `json:"price_each"`
}

func PostOrderDetail(res http.ResponseWriter, req *http.Request) {
	var orderDetail *OrderDetailRequest
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&orderDetail)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsZero([]any{
		orderDetail.Prod_id,
		orderDetail.Quan_ordered,
		orderDetail.Price_each,
		orderDetail.Ord_id}) {
		ErrorRes(res, http.StatusBadRequest,
			"quan_ordered, price_each, ord_id cannot be zero")
		return
	}

	if !isOrderIdInDb(orderDetail.Ord_id) {
		ErrorRes(res, http.StatusBadRequest,
			"ord_id is invalid")
		return
	}

	if !isProdIdInDb(orderDetail.Prod_id) {
		ErrorRes(res, http.StatusBadRequest,
			"prod_id is invalid")
		return
	}

	rows, err := AddOrderDetail(orderDetail)
	if err != nil {
		fmt.Println(orderDetail.Ord_id)
		fmt.Println(orderDetail.Prod_id)

		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding order detail, make sure ord_id and prod_id are also in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Order detail added successfully. %d rows affected", rows),
	})
}

// CUSTOMERS
type CustomerSignUp struct {
	First_name string  `json:"first_name"`
	Last_name  string  `json:"last_name"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Cred_limit float64 `json:"cred_limit"`
}

func PostCustomerSignUp(res http.ResponseWriter, req *http.Request) {
	var customerSignUp *CustomerSignUp
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&customerSignUp)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		customerSignUp.First_name,
		customerSignUp.Last_name,
		customerSignUp.Email,
		customerSignUp.Password}) {
		ErrorRes(res, http.StatusBadRequest,
			"first_name, last_name, email, password cannot be empty")
		return
	}

	if len(customerSignUp.Password) < 8 {
		ErrorRes(res, http.StatusBadRequest,
			"password must be at least 8 characters")
		return
	}

	if isEmailInDb(customerSignUp.Email, "customer") {
		ErrorRes(res, http.StatusBadRequest,
			"email is already in use")
		return
	}

	if !isEmailValid(customerSignUp.Email) {
		ErrorRes(res, http.StatusBadRequest,
			"email is invalid")
		return
	}

	rows, err := SignCustomer(customerSignUp)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding customer, try again later."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Customer added successfully. %d rows affected", rows),
	})
}

type CustomerLogIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func PostCustomerLogIn(res http.ResponseWriter, req *http.Request) {
	var customerLogIn *CustomerLogIn
	customer := &Customer{}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&customerLogIn)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		customerLogIn.Email,
		customerLogIn.Password}) {
		ErrorRes(res, http.StatusBadRequest,
			"email, password cannot be empty")
		return
	}

	err = LogInCustomer(customerLogIn, customer)
	if err != nil {
		ErrorRes(res, http.StatusUnauthorized,
			fmt.Sprintf("Incorrect email or password!"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Customer %s %s (%s) logged in successfully", customer.Cust_fname, customer.Cust_lname, customer.Cust_email),
	})
}

type CustAddr struct {
	Addr_id int `json:"addr_id"`
}

func PutCustAddr(res http.ResponseWriter, req *http.Request) {
	var addCustAddr *CustAddr
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&addCustAddr)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if addCustAddr.Addr_id != 0 {
		if !isAddrIdInDb(addCustAddr.Addr_id) {
			ErrorRes(res, http.StatusBadRequest,
				"addr_id is invalid")
			return
		}
	}

	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.ParseInt(cust_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, try again later."))
		log.Println(err)
		return
	}

	rows, err := AddCustAddr(cust_id_int, addCustAddr.Addr_id)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding customer address, make sure addr_id is also in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Customer address added successfully. %d rows affected", rows),
	})
}

// EMPLOYEES
type EmployeeSignUp struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Job_title  string `json:"job_title"`
}

func PostEmployeeSignUp(res http.ResponseWriter, req *http.Request) {
	var employeeSignUp *EmployeeSignUp
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&employeeSignUp)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		employeeSignUp.First_name,
		employeeSignUp.Last_name,
		employeeSignUp.Email,
		employeeSignUp.Password,
		employeeSignUp.Job_title}) {
		ErrorRes(res, http.StatusBadRequest,
			"first_name, last_name, email, password, job_title cannot be empty")
		return
	}

	if len(employeeSignUp.Password) < 8 {
		ErrorRes(res, http.StatusBadRequest,
			"password must be at least 8 characters")
		return
	}

	if isEmailInDb(employeeSignUp.Email, "employee") {
		ErrorRes(res, http.StatusBadRequest,
			"email is already in use")
		return
	}

	if !isEmailValid(employeeSignUp.Email) {
		ErrorRes(res, http.StatusBadRequest,
			"email is invalid")
		return
	}

	rows, err := SignEmployee(employeeSignUp)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding employee, try again later."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Employee added successfully. %d rows affected", rows),
	})
}

type EmployeeLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func PostEmpLogin(res http.ResponseWriter, req *http.Request) {
	var employeeLogin *EmployeeLogin
	employee := &Employee{}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&employeeLogin)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsEmpty([]string{
		employeeLogin.Email,
		employeeLogin.Password}) {
		ErrorRes(res, http.StatusBadRequest,
			"email, password cannot be empty")
		return
	}

	err = LogInEmployee(employeeLogin, employee)
	if err != nil {
		ErrorRes(res, http.StatusUnauthorized,
			fmt.Sprintf("Incorrect email or password!"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Employee %s %s (%s) logged in successfully", employee.Emp_fname, employee.Emp_lname, employee.Emp_email),
	})
}

// EMPLOYEES
// TODO: ADD A PUT ROUTE FOR EDITING EMPLOYEES
// TODO: ADD A DELETE ROUTE FOR EDITING PRODUCTS

// CUSTOMERS
// TODO: ADD A PUT ROUTE TO EDIT CUSTOMER ADDRESS

// ADDRESSES
// TODO: ADD A PUT ROUTE FOR EDITING ADDRESSES
// TODO: ADD A DELETE ROUTE

// PRODUCTS
// TODO: ADD A PUT ROUTE FOR EDITING PRODUCTS
