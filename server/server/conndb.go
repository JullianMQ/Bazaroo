package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

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
		status TEXT DEFAULT 'in cart',
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
		quan_ordered INT NOT NULL
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
	var id int64

	err := db.QueryRow(`INSERT INTO addresses(
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
		) RETURNING addr_id`,
		addr.Addr_line1,
		addr.Addr_line2,
		addr.City,
		addr.State,
		addr.Postal_code,
		strings.ToUpper(addr.Country),
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetAddrByID(id int64) (AddrResponse, error) {
	var addr AddrResponse
	err := db.QueryRow(`SELECT * FROM addresses WHERE addr_id = $1`, id).Scan(
		&addr.Addr_id,
		&addr.Addr_line1,
		&addr.Addr_line2,
		&addr.City,
		&addr.State,
		&addr.Postal_code,
		&addr.Country,
	)
	return addr, err
}

func EditAddr(addr *AddrRequest, addr_id int64) (int64, error) {
	result, err := db.Exec(`UPDATE addresses SET
		addr_line1 = $1,
		addr_line2 = $2,
		city = $3,
		state = $4,
		postal_code = $5,
		country = $6
		WHERE addr_id = $7`,
		addr.Addr_line1,
		addr.Addr_line2,
		addr.City,
		addr.State,
		addr.Postal_code,
		strings.ToUpper(addr.Country),
		addr_id)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func DeleteAddrById(addr_id int64) (int64, error) {
	result, err := db.Exec(`DELETE FROM addresses WHERE addr_id = $1`, addr_id)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func AddOffice(office *OfficeRequest) (int64, error) {
	var officeId int64

	err := db.QueryRow(`INSERT INTO offices(
			phone_num,
			addr_id
		)
		VALUES (
			$1,
			$2
		) RETURNING office_id`,
		office.Phone_num,
		office.Addr_id,
	).Scan(&officeId)
	if err != nil {
		return 0, err
	}
	return officeId, nil
}

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

func EditEmpOffice(emp_id int64, office_id int) (int64, error) {
	result, err := db.Exec(`UPDATE employees SET office_id = $1 WHERE emp_id = $2`,
		office_id,
		emp_id)
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
		msrp,
		office_id
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9
	)`,
		product.Prod_name,
		product.Prod_line_name,
		product.Prod_vendor_id,
		product.Prod_desc,
		product.Prod_image,
		product.Quan_in_stock,
		product.Buy_price,
		product.Msrp,
		product.Office_id)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func GetProdById(id int64) (Product, error) {
	var prod Product
	err := db.QueryRow(`SELECT
		prod_id,
		prod_name,
		prod_line_name,
		prod_desc,
		prod_image,
		quan_in_stock,
		buy_price
		FROM products WHERE prod_id = $1`, id).Scan(
		&prod.Prod_id,
		&prod.Prod_name,
		&prod.Prod_line_name,
		&prod.Prod_desc,
		&prod.Prod_image,
		&prod.Quan_in_stock,
		&prod.Buy_price)
	if err != nil {
		return prod, err
	}
	return prod, nil
}

func PutBoughtProdById(id int64, product *BoughtProdById) (int64, error) {
	prod, err := GetProdById(id)
	newQuan := prod.Quan_in_stock - product.Quan_bought

	if product.Quan_bought > prod.Quan_in_stock {
		return 0, errors.New("Quan_bought cannot be greater than Quan_in_stock")
	}

	result, err := db.Exec(`UPDATE products SET
		quan_in_stock = $2
		WHERE prod_id = $1`,
		id, newQuan)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	return rows, err
}

func GetCustById(id int64) (Customer, error) {
	var cust Customer
	err := db.QueryRow(`SELECT
		cust_id,
		cust_fname,
		cust_lname,
		cust_email,
		phone_num,
		addr_id,
		cred_limit
		FROM customers WHERE cust_id = $1`,
		id).Scan(&cust.Cust_id, &cust.Cust_fname, &cust.Cust_lname, &cust.Cust_email, &cust.Phone_num, &cust.Addr_id, &cust.Cred_limit)
	if err != nil {
		return cust, err
	}
	return cust, nil
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

func AddCustAddr(cust_id int64, addr_id int) (int64, error) {
	var result sql.Result
	var err error

	if addr_id == 0 {
		result, err = db.Exec(`UPDATE customers SET
			addr_id = $2
			WHERE cust_id = $1`,
			cust_id,
			nil)
		return 0, err
	}

	result, err = db.Exec(`UPDATE customers SET 
		addr_id = $2
		WHERE cust_id = $1`,
		cust_id,
		addr_id)
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
	var id int64
	err = db.QueryRow(`INSERT INTO orders (
		cust_id,
		status,
		comments,
		rating
	) VALUES (
		$1,
		$2,
		$3,
		$4
	) RETURNING ord_id`,
		order.Cust_id,
		order.Status,
		order.Comments,
		order.Rating).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func CheckOutCartQuery(cust_id int64) (int64, error) {
	var id int64
	err := db.QueryRow(`UPDATE orders
		SET status = 'paid'
		WHERE cust_id = $1
		AND status = 'in cart' RETURNING ord_id`, cust_id).Scan(&id)

	if err != nil {
		return 0, err
	}
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
		quan_ordered
	) VALUES (
		$1,
		$2,
		$3
	)`,
		orderDetail.Ord_id,
		orderDetail.Prod_id,
		orderDetail.Quan_ordered)
	if err != nil {
		return 0, err
	}
	id, err := result.RowsAffected()
	return id, err
}

func CheckIfOrderPending(cust_id int64) int {
	var ord_id int
	err := db.QueryRow(`SELECT ord_id FROM orders WHERE cust_id = $1 AND status = 'in cart'`, cust_id).Scan(&ord_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Println("Database error:", err)
		return 0
	}
	return ord_id
}

func AddToCart(cust_id int64, orderdetail CartOrderDetail) (string, error) {
	ord_id := CheckIfOrderPending(cust_id)

	if ord_id == 0 {
		newOrderID, err := AddOrder(&OrderRequest{
			Cust_id:  int(cust_id),
			Status:   "in cart",
			Comments: "",
			Rating:   0,
		})
		if err != nil {
			return "", errors.New(fmt.Sprintf("Error adding to cart: %v", err))
		}
		ord_id = int(newOrderID)
	}

	res, err := db.Exec(`INSERT INTO order_details (ord_id, prod_id, quan_ordered) VALUES ($1, $2, $3)`,
		ord_id, orderdetail.Prod_id, orderdetail.Quan_ordered)
	if err != nil {
		return "", err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Successfully added product. %d rows affected", rows), nil
}

func EditOrderDetailQuantity(order_id int64, prod_id int64, quan_ordered int) (int64, error) {
	result, err := db.Exec(`UPDATE order_details
		SET quan_ordered = $1
		WHERE ord_id = $2 AND prod_id = $3`, quan_ordered, order_id, prod_id)

	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	return rows, err
}

func OrderInCart(cust_id int64) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT
		o.ord_id,
		p.prod_name as name,
		p.prod_id,
		p.prod_image,
		od.quan_ordered as quantity,
		p.buy_price as price
		FROM order_details od
		JOIN orders o ON od.ord_id = o.ord_id
		JOIN products p ON od.prod_id = p.prod_id
		WHERE o.status = 'in cart' AND o.cust_id = $1;`, cust_id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func OrderInPaid(cust_id int64) (*sql.Rows, error) {
	rows, err := db.Query(`SELECT
		o.ord_id,
		p.prod_name as name,
		p.prod_id,
		p.prod_image,
		od.quan_ordered as quantity,
		p.buy_price as price
		FROM order_details od
		JOIN orders o ON od.ord_id = o.ord_id
		JOIN products p ON od.prod_id = p.prod_id
		WHERE o.status = 'paid' AND o.cust_id = $1;`, cust_id)
	if err != nil {
		return nil, err
	}
	return rows, nil
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
