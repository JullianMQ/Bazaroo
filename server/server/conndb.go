package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnDB() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	serviceURI := os.Getenv("AIVEN_DB_URI")
	conn, _ := url.Parse(serviceURI)
	conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"

	db, err = sql.Open("postgres", conn.String())

	if err != nil {
		log.Fatal(err)
	}
	// TODO: ADD A WAY TO CLOSE CONNECTION AFTER QUERY
	// defer db.Close()
}

func CreateSchema() {
	// addresses
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS addresses(
		addr_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		addr_line1 TEXT NOT NULL,
		addr_line2 TEXT,
		city TEXT NOT NULL,
		state TEXT NOT NULL,
		postal_code VARCHAR(10) NOT NULL,
		country CHAR(3) NOT NULL
	)`)
	if err != nil {
		panic(err)
	}

	// offices -> addresses
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS offices (
		office_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		phone_num TEXT NOT NULL,
		addr_id INT NOT NULL REFERENCES addresses(addr_id)
	)`)
	if err != nil {
		panic(err)
	}

	// employees -> offices
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS employees (
		emp_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		emp_fname TEXT NOT NULL,
		emp_lname TEXT NOT NULL,
		emp_email TEXT NOT NULL,
		office_id INT REFERENCES offices(office_id),
		job_title TEXT NOT NULL,
		emp_pass TEXT NOT NULL
	)`)
	if err != nil {
		panic(err)
	}

	// customers -> addresses, employees
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS customers (
		cust_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		cust_fname TEXT NOT NULL,
		cust_lname TEXT NOT NULL,
		cust_email TEXT NOT NULL UNIQUE,
		phone_num TEXT,
		addr_id INT REFERENCES addresses(addr_id),
		cred_limit NUMERIC(10,2) NOT NULL,
		cust_pass TEXT NOT NULL
	)`)
	if err != nil {
		panic(err)
	}

	// vendors -> addresses
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vendors (
		vendor_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		vendor_name TEXT NOT NULL,
		vendor_email TEXT NOT NULL UNIQUE,
		vendor_phone_num TEXT NOT NULL,
		addr_id INT REFERENCES addresses(addr_id)
	)`)
	if err != nil {
		panic(err)
	}

	// orders -> customers
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		ord_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		cust_id INT NOT NULL REFERENCES customers(cust_id),
		ord_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		req_shipped_date DATE NOT NULL,
		comments TEXT,
		rating INT
	)`)
	if err != nil {
		panic(err)
	}

	// payments -> customers, orders
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS payments (
		payment_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		cust_id INT NOT NULL REFERENCES customers(cust_id),
		payment_date TIMESTAMP,
		amount NUMERIC(10,2) NOT NULL,
		payment_status TEXT,
		ord_id INT NOT NULL REFERENCES orders(ord_id)
	)`)
	if err != nil {
		panic(err)
	}

	// product_lines
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS product_lines (
		prod_line_name TEXT PRIMARY KEY,
		prod_line_desc TEXT
	)`)
	if err != nil {
		panic(err)
	}

	// products -> product_lines
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
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

	// order_details -> orders, products
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS order_details (
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

func TestQuery() {
	// WARNING: TEST DATA FOR TEST QUERY SELECT
	_, err := db.Exec(`INSERT INTO addresses(
			addr_line1,
			city,
			state,
			postal_code,
			country
		)
		VALUES (
			'1234 A Avenue St.',
			'Angeles',
			'Pampanga',
			'2010',
			'PHL'
		)`)
	if err != nil {
		panic(err)
	}

	// WARNING: TEST GET THE DATA FROM INSERTED QUERY
	rows, err := db.Query(`SELECT
		addr_id,
		addr_line1,
		city,
		state,
		postal_code,
		country
		FROM addresses`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			addr_id     int
			addr_line1  string
			city        string
			state       string
			postal_code string
			country     string
		)
		if err := rows.Scan(
			&addr_id,
			&addr_line1,
			&city,
			&state,
			&postal_code,
			&country); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Address ID: %v \n", addr_id)
		fmt.Printf("Address Line1: %v \n", addr_line1)
		fmt.Printf("City: %v \n", city)
		fmt.Printf("State: %v \n", state)
		fmt.Printf("Postal Code: %v \n", postal_code)
		fmt.Printf("Country: %v \n", country)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func AddAddr(addr *AddrRequest) (int64, error) {
	res, err := db.Exec(`INSERT INTO addresses(
			addr_line1,
			addr_line2,
			city,
			state,
			postal_code,
			country
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)`,
		addr.Addr_line1,
		addr.Addr_line2,
		addr.City,
		addr.State,
		addr.Postal_code,
		strings.ToUpper(addr.Country),
	)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, nil
}

func AddOffice(office *OfficeRequest) (int64, error) {
	res, err := db.Exec(`INSERT INTO offices(
			phone_num,
			addr_id
		)
		VALUES (
			$1,
			$2
		)`,
		office.Phone_num,
		office.Addr_id,
	)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	return rows, nil
}

// TODO: GET ALL EMPLOYEES

func GetEmpById(id int64) (Employee, error) {
	var emp Employee
	err := db.QueryRow(`SELECT
		emp_id,
		emp_fname,
		emp_lname,
		emp_email,
		office_id,
		job_title
		FROM employees
		WHERE emp_id = $1`, id).Scan(
		&emp.Emp_id,
		&emp.Emp_fname,
		&emp.Emp_lname,
		&emp.Emp_email,
		&emp.Office_id,
		&emp.Job_title,
	)
	return emp, err
}

func AddEmp(emp *EmployeeRequest) (int64, error) {
	result, err := db.Exec(`INSERT INTO employees (
		emp_fname,
		emp_lname,
		emp_email,
		office_id,
		job_title
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)`,
		emp.Emp_fname,
		emp.Emp_lname,
		emp.Emp_email,
		emp.Office_id,
		emp.Job_title)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func AddVendor(vendor *VendorRequest) (int64, error) {
	result, err := db.Exec(`INSERT INTO vendors (vendor_name, vendor_email, vendor_phone_num, addr_id) VALUES ($1, $2, $3, $4)`,
		vendor.Vendor_name,
		vendor.Vendor_email,
		vendor.Vendor_phone_num,
		vendor.Addr_id)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func GetVendorById(id int64) (Vendor, error) {
	var vendor Vendor
	err := db.QueryRow(`SELECT
		vendor_id,
		vendor_name,
		vendor_email,
		vendor_phone_num,
		addr_id
		FROM vendors WHERE vendor_id = $1`,
		id).Scan(&vendor.Vendor_id, &vendor.Vendor_name, &vendor.Vendor_email, &vendor.Vendor_phone_num, &vendor.Addr_id)
	if err != nil {
		return vendor, err
	}
	return vendor, nil
}

func AddProductLine(productLine *ProductLineRequest) (int64, error) {
	rows, err := db.Exec(`INSERT INTO product_lines (
		prod_line_name,
		prod_line_desc
	) VALUES (
		$1,
		$2
	)`,
		productLine.Prod_line_name,
		productLine.Prod_line_desc)
	if err != nil {
		return 0, err
	}
	id, err := rows.RowsAffected()
	return id, err
}

func AddProduct(product *ProductRequest) (int64, error) {
	result, err := db.Exec(`INSERT INTO products (
		prod_name,
		prod_line_name,
		prod_vendor_id,
		prod_desc,
		prod_image,
		quan_in_stock,
		buy_price,
		msrp
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)`,
		product.Prod_name,
		product.Prod_line_name,
		product.Prod_vendor_id,
		product.Prod_desc,
		product.Prod_image,
		product.Quan_in_stock,
		product.Buy_price,
		product.Msrp)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func AddCustomer(customer *CustomerRequest) (int64, error) {
	result, err := db.Exec(`INSERT INTO customers (
		cust_fname,
		cust_lname,
		cust_email,
		phone_num,
		addr_id,
		sales_rep_emp_id,
		cred_limit
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)`,
		customer.Cust_fname,
		customer.Cust_lname,
		customer.Cust_email,
		customer.Phone_num,
		customer.Addr_id,
		customer.Sales_rep_emp_id,
		customer.Cred_limit)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func SignCustomer(customer *CustomerSignUp) (int64, error) {
	result, err := db.Exec(`INSERT INTO customers (
		cust_fname,
		cust_lname,
		cust_email,
		phone_num,
		cred_limit,
		cust_pass
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		md5($6)
	)`,
		customer.First_name,
		customer.Last_name,
		customer.Email,
		"",
		20000,
		customer.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func AddOrder(order *OrderRequest) (int64, error) {
	var err error
	result, err := db.Exec(`INSERT INTO orders (
		cust_id,
		ord_date,
		req_shipped_date,
		comments,
		rating
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)`,
		order.Cust_id,
		order.Ord_date,
		// Add 1 day and take out the hour, minute, second, and nanosecond
		order.Req_shipped_date.Truncate(24*time.Hour),
		order.Comments,
		order.Rating)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func AddPayment(payment *PaymentRequest) (int64, error) {
	result, err := db.Exec(`INSERT INTO payments (
		cust_id,
		payment_date,
		amount,
		payment_status,
		ord_id
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)`,
		payment.Cust_id,
		payment.Payment_date,
		payment.Amount,
		payment.Payment_status,
		payment.Ord_id)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func AddOrderDetail(orderDetail *OrderDetailRequest) (int64, error) {
	result, err := db.Exec(`INSERT INTO order_details (
		ord_id,
		prod_id,
		quan_ordered,
		price_each
	) VALUES (
		$1,
		$2,
		$3,
		$4
	)`,
		orderDetail.Ord_id,
		orderDetail.Prod_id,
		orderDetail.Quan_ordered,
		orderDetail.Price_each)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func LogInCustomer(clog *CustomerLogIn, cust *Customer) error {
	result, err := db.Query(`SELECT
		cust_id,
		cust_fname,
		cust_lname,
		cust_email,
		phone_num,
		addr_id,
		cred_limit
		FROM customers
		WHERE cust_email = $1 AND cust_pass = md5($2)`,
		clog.Email,
		clog.Password)
	if err != nil {
		return err
	}
	defer result.Close()

	if result.Next() == false {
		return errors.New("customer not found")
	}

	if err := result.Scan(
		&cust.Cust_id,
		&cust.Cust_fname,
		&cust.Cust_lname,
		&cust.Cust_email,
		&cust.Phone_num,
		&cust.Addr_id,
		&cust.Cred_limit); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SignEmployee(employee *EmployeeSignUp) (int64, error) {
	result, err := db.Exec(`INSERT INTO employees (
		emp_fname,
		emp_lname,
		emp_email,
		job_title,
		emp_pass
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		md5($5)
	)`,
		employee.First_name,
		employee.Last_name,
		employee.Email,
		employee.Job_title,
		employee.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func LogInEmployee(clog *EmployeeLogin, emp *Employee) error {
	result, err := db.Query(`SELECT
		emp_id,
		emp_fname,
		emp_lname,
		emp_email,
		job_title,
		office_id
		FROM employees
		WHERE emp_email = $1 AND emp_pass = md5($2)`,
		clog.Email,
		clog.Password)
	if err != nil {
		return err
	}
	defer result.Close()

	if result.Next() == false {
		return errors.New("employee not found")
	}

	if err := result.Scan(
		&emp.Emp_id,
		&emp.Emp_fname,
		&emp.Emp_lname,
		&emp.Emp_email,
		&emp.Job_title,
		&emp.Office_id); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
