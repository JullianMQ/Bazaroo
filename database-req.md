# Database Requirements

## Customers
id: int primary key generated always as identity
username: varchar(100)
email: varchar(100)

## Payments


## Orders


## OrderDetails


## Products


## ProductsLines


## Employees


## Offices

## Database QUERIES

```postgresql
CREATE TABLE Addresses (
  addr_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  addr_line1 text NOT NULL,
  addr_line2 text NOT NULL,
  city text NOT NULL,
  state text NOT NULL,
  postal_code varchar(10) NOT NULL,
  country char(3) NOT NULL
);

CREATE TABLE Offices (
  office_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  phone_num text NOT NULL,
  addr_id int NOT NULL,
  FOREIGN KEY (addr_id) REFERENCES Addresses(addr_id)
);

CREATE TABLE Employees (
  emp_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  emp_fname text NOT NULL,
  emp_lname text NOT NULL,
  emp_email text UNIQUE NOT NULL,
  office_id int,
  job_title text NOT NULL,
  FOREIGN KEY(office_id) REFERENCES Offices(office_id)
);

CREATE TABLE Customers (
  cust_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  cust_fname text NOT NULL,
  cust_lname text NOT NULL,
  cust_email text UNIQUE NOT NULL,
  phone_num varchar(25),
  addr_id int NOT NULL,
  sales_rep_emp_id int,
  cred_limit decimal(10, 2) NOT NULL,
  FOREIGN KEY (addr_id) REFERENCES Addresses(addr_id),
  FOREIGN KEY (sales_rep_emp_id) REFERENCES Employees(emp_id)
);
 
CREATE TABLE Orders (
  ord_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  cust_id int,
  ord_date timestamp NOT NULL default now(),
  req_shipped_date date NOT NULL,
  comments text,
  FOREIGN KEY (cust_id) REFERENCES Customers(cust_id)
);

CREATE TABLE ProductLines (
  prod_line_name text PRIMARY KEY,
  prod_line_desc text
);

CREATE TABLE Vendors (
  vendor_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  vendor_name text NOT NULL UNIQUE,
  vendor_email text NOT NULL UNIQUE,
  vendor_phone_num text,
  addr_id int,
  FOREIGN KEY (addr_id) REFERENCES Addresses(addr_id)
);

CREATE TABLE Products (
  prod_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  prod_name text NOT NULL,
  prod_line_name text,
  prod_vendor_id int,
  prod_desc text,
  quan_in_stock int,
  buy_price decimal(10,2),
  msrp decimal(10,2),
  FOREIGN KEY (prod_line_name) REFERENCES ProductLines(prod_line_name),
  FOREIGN KEY (prod_vendor_id) REFERENCES Vendors(vendor_id)
);

CREATE TABLE OrderDetails (
  ord_id int,
  prod_id int,
  quan_ordered int NOT NULL,
  status text NOT NULL,
  price_each decimal(10,2),
  PRIMARY KEY (ord_id, prod_id),
  FOREIGN KEY (ord_id) REFERENCES Orders(ord_id),
  FOREIGN KEY (prod_id) REFERENCES Products(prod_id)
);

CREATE TABLE Payments (
  payment_id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  cust_id int,
  payment_date timestamp NOT NULL default now(),
  amount decimal(10,2) NOT NULL,
  payment_status text NOT NULL,
  ord_id int,
  FOREIGN KEY (cust_id) REFERENCES Customers(cust_id),
  FOREIGN KEY (ord_id) REFERENCES Orders(ord_id)
);
```

## Database Normalization Level

### 1NF
- It adheres to the standard that all the attributes in the database are atomic, meaning each value is indivisible
### 2NF
- It is in 1NF
- It adheres to the standard that all the non-key attributes depend on the whole primary key
### 3NF
- It is in 2NF
- It adheres to the standard that no non-key attribute depends on another non-key attribute and that there should be no transitive dependencies.
### BCNF
- It is in 3NF
- All non-key attributes depend solely on the primary key, every determinant must be a candidate key
### 4NF
- It is in BCNF
- There are no multi-valued dependencies, meaning that each fact is stored separately, without repeating or redundant data.
### 5NF
- It is in 4NF
- There are no join dependencies, meaning that we cannot break down any table further without losing information when joining them.

### Notes on how what the difference is between 3NF and BCNF
- If we were to add something like `sales_rep_emp_name` which would be determined by the `sales_rep_emp_id` it would violate the rule for BCNF as it is not a **candidate key** but a foreign key instead.

