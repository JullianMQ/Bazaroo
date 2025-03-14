package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"regexp"
	"slices"
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

func ContainsZero(v []int) bool {
	return slices.Contains(v, 0)
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

func isEmailValid(e string) bool {
	if _, err := mail.ParseAddress(e); err != nil {
		return false
	}
	return true
}

func isEmailInDb(e string) bool {
	rows, err := db.Query(`SELECT
		emp_email
		FROM employees
		WHERE emp_email = $1`, e)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var emp_email string
		if err := rows.Scan(&emp_email); err != nil {
			log.Println(err)
		}
		if emp_email == e {
			return true
		}
	}
	return false
}

func isPhoneValid(p string) bool {
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(p)
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
		panic(err)
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
	Emp_id    int    `json:"emp_id"`
	Emp_fname string `json:"emp_fname"`
	Emp_lname string `json:"emp_lname"`
	Emp_email string `json:"emp_email"`
	Office_id int    `json:"office_id"`
	Job_title string `json:"job_title"`
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
			office_id int
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

	if isEmailInDb(emp.Emp_email) {
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
	Vendor_id int `json:"vendor_id"`
	Vendor_name string `json:"vendor_name"`
	Vendor_email string `json:"vendor_email"`
	Vendor_phone_num string `json:"vendor_phone_num"`
	Addr_id int `json:"addr_id"`
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
		vendor := Vendor {
			Vendor_id: vendor_id,
			Vendor_name: vendor_name,
			Vendor_email: vendor_email,
			Vendor_phone_num: vendor_phone_num,
			Addr_id: addr_id,
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
	Vendor_name string `json:"vendor_name"`
	Vendor_email string `json:"vendor_email"`
	Vendor_phone_num string `json:"vendor_phone_num"`
	Addr_id int `json:"addr_id"`
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

	if isEmailInDb(vendor.Vendor_email) {
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
	Prod_line_name string `json:"prod_line_name"`
	Prod_line_desc string `json:"prod_line_desc"`
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
			prod_line_desc string
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
		productLine.Prod_line_name,
		productLine.Prod_line_desc}) {
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
	Prod_id int `json:"prod_id"`
	Prod_name string `json:"prod_name"`
	Prod_line_name string `json:"prod_line_name"`
	Prod_vendor_id int `json:"prod_vendor_id"`
	Prod_desc string `json:"prod_desc"`
	Quan_in_stock int `json:"quan_in_stock"`
	Buy_price float64 `json:"buy_price"`
	Msrp float64 `json:"msrp"`
}

func GetProducts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := db.Query(`SELECT
		prod_id,
		prod_name,
		prod_line_name,
		prod_vendor_id,
		prod_desc,
		quan_in_stock,
		buy_price,
		msrp
		FROM products`)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var (
			prod_id int
			prod_name string
			prod_line_name string
			prod_vendor_id int
			prod_desc string
			quan_in_stock int
			buy_price float64
			msrp float64
		)
		if err := rows.Scan(
			&prod_id,
			&prod_name,
			&prod_line_name,
			&prod_vendor_id,
			&prod_desc,
			&quan_in_stock,
			&buy_price,
			&msrp); err != nil {
			log.Println(err)
			return
		}
		products = append(products, Product{
			Prod_id: prod_id,
			Prod_name: prod_name,
			Prod_line_name: prod_line_name,
			Prod_vendor_id: prod_vendor_id,
			Prod_desc: prod_desc,
			Quan_in_stock: quan_in_stock,
			Buy_price: buy_price,
			Msrp: msrp,
		})
	}
	json.NewEncoder(res).Encode(products)
}
