package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetRoot(res http.ResponseWriter, req *http.Request) {
	myMess := "Message"
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(myMess)
}

// =====================================ADDRESSES========================================
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

	rows, err := GetAddrQuery()
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

func GetAddrID(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := req.URL.Query().Get("id")
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, try again later."))
		log.Println(err)
		return
	}

	addr, err := GetAddrByIDQuery(id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not get address, try again later."))
		log.Println(err)
		return
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

	inserted_id, err := AddAddrQuery(addr)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not add address, try again later."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("%d", inserted_id),
	})
}

func PutAddr(res http.ResponseWriter, req *http.Request) {
	var addr *AddrRequest
	var addr_id string
	addr_id = req.URL.Query().Get("id")
	addr_id_int, err := strconv.ParseInt(addr_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id into int64, try again later."))
		log.Println(err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewDecoder(req.Body).Decode(&addr)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
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

	rows, err := EditAddrQuery(addr, addr_id_int)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error editing address, make sure addr_id is in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Address edited successfully. %d rows affected", rows),
	})
}

func DeleteAddr(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	addr_id := req.URL.Query().Get("id")
	addr_id_int, err := strconv.ParseInt(addr_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id into int64, try again later."))
		log.Println(err)
		return
	}

	rows, err := DeleteAddrByIdQuery(addr_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not delete address, try again later."))
		log.Println(err)
		return
	}

	if rows == 0 {
		json.NewEncoder(res).Encode(OkResponse{
			Message: "Address not found",
		})
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Address deleted successfully. %d rows affected", rows),
	})
}
// ========================================ADDRESSES========================================

// =====================================OFFICE========================================
type Office struct {
	Office_id int    `json:"office_id"`
	Phone_num string `json:"phone_num"`
	Addr_id   int    `json:"addr_id"`
}

type OfficeRequest struct {
	Phone_num string `json:"phone_num"`
	Addr_id   int    `json:"addr_id"`
}

func GetOffices(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := GetOfficesQuery()
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
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

	office_id, err := AddOfficeQuery(office)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding office, make sure addr_id is in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("%d", office_id),
	})
}
// ========================================OFFICE========================================

// =====================================EMPLOYEES=====================================
type Employee struct {
	Emp_id    int           `json:"emp_id"`
	Emp_fname string        `json:"emp_fname"`
	Emp_lname string        `json:"emp_lname"`
	Emp_email string        `json:"emp_email"`
	Office_id sql.NullInt64 `json:"office_id"`
	Job_title string        `json:"job_title"`
}

type EmployeeOffice struct {
	Office_id int `json:"office_id"`
}

type EmployeeRequest struct {
	Emp_fname string `json:"emp_fname"`
	Emp_lname string `json:"emp_lname"`
	Emp_email string `json:"emp_email"`
	Office_id int    `json:"office_id"`
	Job_title string `json:"job_title"`
}

type EmployeeSignUp struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Job_title  string `json:"job_title"`
}

type EmployeeLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	emp, err := GetEmpByIdQuery(emp_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not get employee by id, check if id is correct."))
		log.Println(err)
		return
	}

	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("%v", emp.Office_id.Valid),
	})
}

func GetEmps(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := GetEmpsQuery()
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
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

func PutEmpOffice(res http.ResponseWriter, req *http.Request) {
	var empOffice *EmployeeOffice
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&empOffice)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	emp_id := req.URL.Query().Get("id")
	emp_id_int, err := strconv.ParseInt(emp_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, try again later."))
		log.Println(err)
		return
	}

	if !isEmpIdInDb(int(emp_id_int)) {
		ErrorRes(res, http.StatusBadRequest,
			"emp_id is invalid")
		log.Println(err)
		return
	}

	rows, err := EditEmpOfficeQuery(emp_id_int, empOffice.Office_id)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error editing employee office, make sure office_id is in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Employee office edited successfully. %d rows affected", rows),
	})
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

	_, err = AddEmpQuery(emp)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("CHECK IF OFFICE ID IN DATABASE: %s", err))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Employee added successfully."),
	})
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

	_, err = SignEmployeeQuery(employeeSignUp)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding employee, try again later."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Employee added successfully."),
	})
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

	err = LogInEmployeeQuery(employeeLogin, employee)
	if err != nil {
		ErrorRes(res, http.StatusUnauthorized,
			fmt.Sprintf("Incorrect email or password!"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("%d", employee.Emp_id),
	})
}
// =====================================EMPLOYEES=====================================

// =====================================VENDOR=====================================
type Vendor struct {
	Vendor_id        int    `json:"vendor_id"`
	Vendor_name      string `json:"vendor_name"`
	Vendor_email     string `json:"vendor_email"`
	Vendor_phone_num string `json:"vendor_phone_num"`
	Addr_id          int    `json:"addr_id"`
}

type VendorRequest struct {
	Vendor_name      string `json:"vendor_name"`
	Vendor_email     string `json:"vendor_email"`
	Vendor_phone_num string `json:"vendor_phone_num"`
	Addr_id          int    `json:"addr_id"`
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

	vendor, err := GetVendorByIdQuery(vendor_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Error getting vendor")
		log.Println(err)
		return
	}

	json.NewEncoder(res).Encode(vendor)
}

func GetVendors(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := GetVendorsQuery()
	if err != nil {
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

	rows, err := AddVendorQuery(vendor)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding vendor, make sure addr_id is in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Vendor added successfully. %d rows affected", rows),
	})
}
// =====================================VENDOR=====================================

// =====================================PRODUCT LINE=====================================
type ProductLine struct {
	Prod_line_name string         `json:"prod_line_name"`
	Prod_line_desc sql.NullString `json:"prod_line_desc"`
}

type ProductLineRequest struct {
	Prod_line_name string `json:"prod_line_name"`
	Prod_line_desc string `json:"prod_line_desc"`
}

func GetProductLine(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := GetProductLinesQuery()
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

	rows, err := AddProductLineQuery(productLine)
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
// =====================================PRODUCT LINE=====================================

// =====================================PRODUCT=====================================
type ProductById struct {
	Prod_id        int            `json:"prod_id"`
	Prod_name      string         `json:"prod_name"`
	Prod_line_name string         `json:"prod_line_name"`
	Prod_desc      sql.NullString `json:"prod_desc"`
	Prod_image     sql.NullString `json:"prod_image"`
	Quan_in_stock  int            `json:"quan_in_stock"`
	Buy_price      float64        `json:"buy_price"`
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

type ProductRequest struct {
	Prod_name      string  `json:"prod_name"`
	Prod_line_name string  `json:"prod_line_name"`
	Prod_vendor_id int     `json:"prod_vendor_id"`
	Prod_desc      string  `json:"prod_desc"`
	Prod_image     string  `json:"prod_image"`
	Quan_in_stock  int     `json:"quan_in_stock"`
	Buy_price      float64 `json:"buy_price"`
	Msrp           float64 `json:"msrp"`
	Office_id      float64 `json:"office_id"`
}

type BoughtProdById struct {
	Quan_bought int `json:"quan_bought"`
}

func GetProductById(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	prod_id := req.URL.Query().Get("id")
	prod_id_int, err := strconv.ParseInt(prod_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id into int64, try again later."))
		log.Println(err)
		return
	}
	rows, err := GetProdByIdQuery(prod_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not get product by id, check if id is correct."))
		log.Println(err)
		return
	}

	json.NewEncoder(res).Encode(ProductById{
		Prod_id:        rows.Prod_id,
		Prod_name:      rows.Prod_name,
		Prod_line_name: rows.Prod_line_name,
		Prod_desc:      rows.Prod_desc,
		Prod_image:     rows.Prod_image,
		Quan_in_stock:  rows.Quan_in_stock,
		Buy_price:      rows.Buy_price,
	})
}

func GetProducts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := GetProductsQuery()
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
			Prod_image: sql.NullString{
				String: fmt.Sprintf("/v1/images/%s", prod_image.String),
				Valid:  prod_image.Valid},
			Quan_in_stock: quan_in_stock,
			Buy_price:     buy_price,
			Msrp:          msrp,
		})
	}
	json.NewEncoder(res).Encode(products)
}

func PutBoughtProductById(res http.ResponseWriter, req *http.Request) {
	product := &BoughtProdById{}
	res.Header().Set("Content-Type", "application/json")
	prod_id := req.URL.Query().Get("id")
	prod_id_int, err := strconv.ParseInt(prod_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id into int64, try again later."))
		log.Println(err)
		return
	}
	err = json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	if ContainsZero([]any{
		product.Quan_bought}) {
		ErrorRes(res, http.StatusBadRequest,
			"quan_bought cannot be zero")
		return
	}

	_, err = PutBoughtProdByIdQuery(prod_id_int, product)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error buying product, make sure bought quantity is lower or equal to quantity in stock."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Product edited successfully"),
	})
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

	rows, err := AddProductQuery(product)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding product, make sure prod_vendor_id is in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Product added successfully. %d rows affected", rows),
	})
}
// =====================================PRODUCT=====================================

// =====================================CUSTOMERS========================================
type Customer struct {
	Cust_id    int           `json:"cust_id"`
	Cust_fname string        `json:"cust_fname"`
	Cust_lname string        `json:"cust_lname"`
	Cust_email string        `json:"cust_email"`
	Phone_num  string        `json:"phone_num"`
	Addr_id    sql.NullInt64 `json:"addr_id"`
	Cred_limit float64       `json:"cred_limit"`
}

type CustomerAcc struct {
	Cust_id    int           `json:"cust_id"`
	Cust_fname string        `json:"cust_fname"`
	Cust_lname string        `json:"cust_lname"`
	Cust_email string        `json:"cust_email"`
	Addr_id    sql.NullInt64 `json:"addr_id"`
	Phone_num  string        `json:"phone_num"`
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

type CustomerSignUp struct {
	First_name string  `json:"first_name"`
	Last_name  string  `json:"last_name"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Cred_limit float64 `json:"cred_limit"`
}

type CustomerLogIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CustAddr struct {
	Addr_id int `json:"addr_id"`
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

	cust, err := GetCustByIdQuery(cust_id_int)
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
		Addr_id:    cust.Addr_id,
		Phone_num:  cust.Phone_num,
	})
}

func GetCustomers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := GetCustomersQuery()
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

	_, err = AddCustomerQuery(customer)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding customer, make sure addr_id is in the database"))
		log.Println(err)
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Customer added successfully."),
	})
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

	_, err = SignCustomerQuery(customerSignUp)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding customer, try again later."))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Customer added successfully."),
	})
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

	err = LogInCustomerQuery(customerLogIn, customer)
	if err != nil {
		ErrorRes(res, http.StatusUnauthorized,
			fmt.Sprintf("Incorrect email or password!"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("%d", customer.Cust_id),
	})
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

	rows, err := AddCustAddrQuery(cust_id_int, addCustAddr.Addr_id)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding customer address, make sure addr_id is in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Customer address added successfully. %d rows affected", rows),
	})
}

// ========================================CUSTOMERS========================================

// =====================================ORDERS=====================================
type OrderByCustId struct {
	Ord_id   int    `json:"ord_id"`
	Status   string `json:"status"`
	Comments string `json:"comments"`
	Rating   int    `json:"rating"`
}

type Order struct {
	Ord_id   int            `json:"ord_id"`
	Cust_id  int            `json:"cust_id"`
	Status   string         `json:"status"`
	Comments sql.NullString `json:"comments"`
	Rating   int            `json:"rating"`
}

type OrderRequest struct {
	Cust_id  int    `json:"cust_id"`
	Status   string `json:"status"`
	Comments string `json:"comments"`
	Rating   int    `json:"rating"`
}

type OrdersInCart struct {
	Ord_id       int            `json:"ord_id"`
	Prod_name    string         `json:"prod_name"`
	Prod_id      int            `json:"prod_id"`
	Prod_image   sql.NullString `json:"prod_image"`
	Price        float64        `json:"price"`
	Quan_ordered int            `json:"quan_ordered"`
}

type PaidOrder struct {
	Ord_id       int            `json:"ord_id"`
	Prod_name    string         `json:"prod_name"`
	Prod_id      int            `json:"prod_id"`
	Prod_image   sql.NullString `json:"prod_image"`
	Price        float64        `json:"price"`
	Quan_ordered int            `json:"quan_ordered"`
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

	rows, err := GetOrderByCustIdQuery(cust_id_int)
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
			ord_id   int
			cust_id  int
			status   string
			comments string
			rating   int
		)
		if err := rows.Scan(
			&ord_id,
			&cust_id,
			&status,
			&comments,
			&rating); err != nil {
			log.Println(err)
			return
		}
		orders = append(orders, OrderByCustId{
			Ord_id:   ord_id,
			Status:   status,
			Comments: comments,
			Rating:   rating,
		})
	}
	json.NewEncoder(res).Encode(orders)
}

func GetOrders(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	rows, err := GetOrdersQuery()
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
			ord_id   int
			cust_id  int
			status   string
			comments sql.NullString
			rating   int
		)
		if err := rows.Scan(
			&ord_id,
			&cust_id,
			&status,
			&comments,
			&rating); err != nil {
			log.Println(err)
			return
		}
		orders = append(orders, Order{
			Ord_id:   ord_id,
			Cust_id:  cust_id,
			Status:   status,
			Comments: comments,
			Rating:   rating,
		})
	}
	json.NewEncoder(res).Encode(orders)
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

	if ContainsZero([]any{
		order.Cust_id}) {
		ErrorRes(res, http.StatusBadRequest,
			"rating cannot be zero")
		return
	}

	if order.Rating < 0 || order.Rating > 5 {
		ErrorRes(res, http.StatusBadRequest,
			"rating must be between 0 and 5")
		return
	}

	rows, err := AddOrderQuery(order)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding order, make sure cust_id is in the database"))
		log.Println(err)
		return
	}
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Order added successfully. %d rows affected", rows),
	})
}

func GetOrderInCart(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.ParseInt(cust_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Could not parse id into int64."))
		log.Println(err)
		return
	}

	rows, err := OrderInCartQuery(cust_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not get order in cart, try again later")
		log.Println(err)
		return
	}
	defer rows.Close()

	orderInCart := []OrdersInCart{}
	for rows.Next() {
		var (
			ord_id     int
			prod_name  string
			prod_id    int
			prod_image sql.NullString
			quantity   int
			price      float64
		)
		if err := rows.Scan(
			&ord_id,
			&prod_name,
			&prod_id,
			&prod_image,
			&quantity,
			&price,
		); err != nil {
			log.Println(err)
		}
		orderInCart = append(orderInCart, OrdersInCart{
			Ord_id:    ord_id,
			Prod_name: prod_name,
			Prod_id:   prod_id,
			Prod_image: sql.NullString{
				String: fmt.Sprintf("/v1/images/%s", prod_image.String),
				Valid:  prod_image.Valid,
			},
			Price:        price,
			Quan_ordered: quantity,
		})
	}
	json.NewEncoder(res).Encode(orderInCart)
}

func GetOrderInPaid(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.ParseInt(cust_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Could not parse id into int64."))
		log.Println(err)
		return
	}

	rows, err := OrderInPaidQuery(cust_id_int)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not get order in paid, try again later")
		log.Println(err)
		return
	}
	defer rows.Close()

	orderPaid := []PaidOrder{}
	for rows.Next() {
		var (
			ord_id     int
			prod_name  string
			prod_id    int
			prod_image sql.NullString
			quantity   int
			price      float64
		)
		if err := rows.Scan(
			&ord_id,
			&prod_name,
			&prod_id,
			&prod_image,
			&quantity,
			&price,
		); err != nil {
			log.Println(err)
		}
		orderPaid = append(orderPaid, PaidOrder{
			Ord_id:    ord_id,
			Prod_name: prod_name,
			Prod_id:   prod_id,
			Prod_image: sql.NullString{
				String: fmt.Sprintf("/v1/images/%s", prod_image.String),
				Valid:  prod_image.Valid},
			Price:        price * float64(quantity),
			Quan_ordered: quantity,
		})
	}
	json.NewEncoder(res).Encode(orderPaid)
}

func CheckOutCart(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.Atoi(cust_id)
	if err != nil {
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: "Invalid Customer ID",
		})
		return
	}

	id, err := CheckOutCartQuery(int64(cust_id_int))
	if err != nil {
		json.NewEncoder(res).Encode(ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Order ID: %d", id),
	})
}
// =====================================ORDERS=====================================

// =====================================PAYMENTS=====================================
type Payment struct {
	Payment_id     int       `json:"payment_id"`
	Cust_id        int       `json:"cust_id"`
	Payment_date   time.Time `json:"payment_date"`
	Amount         float64   `json:"amount"`
	Payment_status string    `json:"payment_status"`
	Ord_id         int       `json:"ord_id"`
}

type PaymentRequest struct {
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

	rows, err := GetPaymentsByCustIdQuery(cust_id_int)
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

	rows, err := AddPaymentQuery(payment)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding payment, make sure cust_id is in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Payment added successfully. %d rows affected", rows),
	})
}
// =====================================PAYMENTS=====================================

// =====================================ORDER DETAILS=====================================
type OrderDetail struct {
	Ord_id       int `json:"ord_id"`
	Prod_id      int `json:"prod_id"`
	Quan_ordered int `json:"quan_ordered"`
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

	rows, err := GetOrderDetailsByOrderIdQuery(order_id_int)
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
		)
		if err := rows.Scan(
			&ord_id,
			&prod_id,
			&quan_ordered); err != nil {
			log.Println(err)
			return
		}
		orderDetails = append(orderDetails, OrderDetail{
			Ord_id:       ord_id,
			Prod_id:      prod_id,
			Quan_ordered: quan_ordered,
		})
	}
	json.NewEncoder(res).Encode(orderDetails)
}

type CartOrderDetail struct {
	Prod_id      int `json:"prod_id"`
	Quan_ordered int `json:"quan_ordered"`
}

func PostAddToCart(res http.ResponseWriter, req *http.Request) {
	var orderdetail CartOrderDetail
	res.Header().Set("Content-Type", "application/json")
	cust_id := req.URL.Query().Get("id")
	cust_id_int, err := strconv.ParseInt(cust_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest, "Invalid Customer ID")
		log.Println(err)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&orderdetail)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest, "Invalid request body")
		log.Println(err)
		return
	}

	suc, err := AddToCartQuery(cust_id_int, orderdetail)
	if err != nil {
		ErrorRes(res, http.StatusBadRequest, fmt.Sprintf("Product or Customer ID invalid %v", err))
		log.Println(err)
		return
	}

	json.NewEncoder(res).Encode(OkResponse{
		Message: suc,
	})
}

type QuantityOrderDetail struct {
	Quan_ordered int `json:"quan_ordered"`
}

func PutOrderDetailQuantity(res http.ResponseWriter, req *http.Request) {
	var quantityOrderDetail *QuantityOrderDetail
	res.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(req.Body).Decode(&quantityOrderDetail)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			"Could not decode request body")
		log.Println(err)
		return
	}

	order_id := req.URL.Query().Get("ord_id")
	order_id_int, err := strconv.ParseInt(order_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, make sure that it is included in the query."))
		log.Println(err)
		return
	}

	prod_id := req.URL.Query().Get("prod_id")
	prod_id_int, err := strconv.ParseInt(prod_id, 10, 64)
	if err != nil {
		ErrorRes(res, http.StatusInternalServerError,
			fmt.Sprintf("Could not parse id, make sure that it is included in the query."))
		log.Println(err)
		return
	}

	rows, err := EditOrderDetailQuantityQuery(order_id_int, prod_id_int, quantityOrderDetail.Quan_ordered)
	if rows == 0 {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error editing order detail quantity, make sure order_id and prod_id is in database"))
		log.Println(err)
		return
	}

	if err != nil {
		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error editing order detail quantity, make sure order_id and prod_id is in database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("%d", rows),
	})
}

type OrderDetailRequest struct {
	Ord_id       int `json:"ord_id"`
	Prod_id      int `json:"prod_id"`
	Quan_ordered int `json:"quan_ordered"`
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

	rows, err := AddOrderDetailQuery(orderDetail)
	if err != nil {

		ErrorRes(res, http.StatusBadRequest,
			fmt.Sprintf("Error adding order detail, make sure ord_id and prod_id are in the database"))
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(OkResponse{
		Message: fmt.Sprintf("Order detail added successfully. %d rows affected", rows),
	})
}
// =====================================ORDER DETAILS=====================================


// EMPLOYEES
// TODO: ADD A PUT ROUTE FOR EDITING EMPLOYEES
// TODO: ADD A DELETE ROUTE FOR EDITING PRODUCTS

// PRODUCTS
// TODO: ADD A PUT ROUTE FOR EDITING PRODUCTS
