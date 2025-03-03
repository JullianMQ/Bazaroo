--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15 (Ubuntu 14.15-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.15 (Ubuntu 14.15-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.addresses (
    addr_id integer NOT NULL,
    addr_line1 text NOT NULL,
    addr_line2 text NOT NULL,
    city text NOT NULL,
    state text NOT NULL,
    postal_code character varying(10) NOT NULL,
    country character(3) NOT NULL
);


ALTER TABLE public.addresses OWNER TO root;

--
-- Name: addresses_addr_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.addresses ALTER COLUMN addr_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.addresses_addr_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: customers; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.customers (
    cust_id integer NOT NULL,
    cust_fname text NOT NULL,
    cust_lname text NOT NULL,
    cust_email text NOT NULL,
    phone_num character varying(25),
    addr_id integer NOT NULL,
    sales_rep_emp_id integer,
    cred_limit numeric(10,2) NOT NULL
);


ALTER TABLE public.customers OWNER TO root;

--
-- Name: customers_cust_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.customers ALTER COLUMN cust_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.customers_cust_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: employees; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.employees (
    emp_id integer NOT NULL,
    emp_fname text NOT NULL,
    emp_lname text NOT NULL,
    emp_email text NOT NULL,
    office_id integer,
    job_title text NOT NULL
);


ALTER TABLE public.employees OWNER TO root;

--
-- Name: employees_emp_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.employees ALTER COLUMN emp_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.employees_emp_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: offices; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.offices (
    office_id integer NOT NULL,
    phone_num text NOT NULL,
    addr_id integer NOT NULL
);


ALTER TABLE public.offices OWNER TO root;

--
-- Name: offices_office_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.offices ALTER COLUMN office_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.offices_office_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: orderdetails; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.orderdetails (
    ord_id integer NOT NULL,
    prod_id integer NOT NULL,
    quan_ordered integer NOT NULL,
    status text NOT NULL,
    price_each numeric(10,2)
);


ALTER TABLE public.orderdetails OWNER TO root;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.orders (
    ord_id integer NOT NULL,
    cust_id integer,
    ord_date timestamp without time zone DEFAULT now() NOT NULL,
    req_shipped_date date NOT NULL,
    comments text
);


ALTER TABLE public.orders OWNER TO root;

--
-- Name: orders_ord_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.orders ALTER COLUMN ord_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.orders_ord_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: payments; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.payments (
    payment_id integer NOT NULL,
    cust_id integer,
    payment_date timestamp without time zone DEFAULT now() NOT NULL,
    amount numeric(10,2) NOT NULL,
    payment_status text NOT NULL,
    ord_id integer
);


ALTER TABLE public.payments OWNER TO root;

--
-- Name: payments_payment_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.payments ALTER COLUMN payment_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.payments_payment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: productlines; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.productlines (
    prod_line_name text NOT NULL,
    prod_line_desc text
);


ALTER TABLE public.productlines OWNER TO root;

--
-- Name: products; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.products (
    prod_id integer NOT NULL,
    prod_name text NOT NULL,
    prod_line_name text,
    prod_vendor_id integer,
    prod_desc text,
    quan_in_stock integer,
    buy_price numeric(10,2),
    msrp numeric(10,2)
);


ALTER TABLE public.products OWNER TO root;

--
-- Name: products_prod_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.products ALTER COLUMN prod_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.products_prod_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: vendors; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.vendors (
    vendor_id integer NOT NULL,
    vendor_name text NOT NULL,
    vendor_email text NOT NULL,
    vendor_phone_num text,
    addr_id integer
);


ALTER TABLE public.vendors OWNER TO root;

--
-- Name: vendors_vendor_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

ALTER TABLE public.vendors ALTER COLUMN vendor_id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.vendors_vendor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.addresses (addr_id, addr_line1, addr_line2, city, state, postal_code, country) FROM stdin;
\.


--
-- Data for Name: customers; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.customers (cust_id, cust_fname, cust_lname, cust_email, phone_num, addr_id, sales_rep_emp_id, cred_limit) FROM stdin;
\.


--
-- Data for Name: employees; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.employees (emp_id, emp_fname, emp_lname, emp_email, office_id, job_title) FROM stdin;
\.


--
-- Data for Name: offices; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.offices (office_id, phone_num, addr_id) FROM stdin;
\.


--
-- Data for Name: orderdetails; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.orderdetails (ord_id, prod_id, quan_ordered, status, price_each) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.orders (ord_id, cust_id, ord_date, req_shipped_date, comments) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.payments (payment_id, cust_id, payment_date, amount, payment_status, ord_id) FROM stdin;
\.


--
-- Data for Name: productlines; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.productlines (prod_line_name, prod_line_desc) FROM stdin;
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.products (prod_id, prod_name, prod_line_name, prod_vendor_id, prod_desc, quan_in_stock, buy_price, msrp) FROM stdin;
\.


--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.vendors (vendor_id, vendor_name, vendor_email, vendor_phone_num, addr_id) FROM stdin;
\.


--
-- Name: addresses_addr_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.addresses_addr_id_seq', 1, false);


--
-- Name: customers_cust_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.customers_cust_id_seq', 1, false);


--
-- Name: employees_emp_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.employees_emp_id_seq', 1, false);


--
-- Name: offices_office_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.offices_office_id_seq', 1, false);


--
-- Name: orders_ord_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.orders_ord_id_seq', 1, false);


--
-- Name: payments_payment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.payments_payment_id_seq', 1, false);


--
-- Name: products_prod_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.products_prod_id_seq', 1, false);


--
-- Name: vendors_vendor_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.vendors_vendor_id_seq', 1, false);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (addr_id);


--
-- Name: customers customers_cust_email_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_cust_email_key UNIQUE (cust_email);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (cust_id);


--
-- Name: employees employees_emp_email_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_emp_email_key UNIQUE (emp_email);


--
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (emp_id);


--
-- Name: offices offices_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.offices
    ADD CONSTRAINT offices_pkey PRIMARY KEY (office_id);


--
-- Name: orderdetails orderdetails_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orderdetails
    ADD CONSTRAINT orderdetails_pkey PRIMARY KEY (ord_id, prod_id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (ord_id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (payment_id);


--
-- Name: productlines productlines_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.productlines
    ADD CONSTRAINT productlines_pkey PRIMARY KEY (prod_line_name);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (prod_id);


--
-- Name: vendors vendors_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_pkey PRIMARY KEY (vendor_id);


--
-- Name: vendors vendors_vendor_email_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_vendor_email_key UNIQUE (vendor_email);


--
-- Name: vendors vendors_vendor_name_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_vendor_name_key UNIQUE (vendor_name);


--
-- Name: customers customers_addr_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_addr_id_fkey FOREIGN KEY (addr_id) REFERENCES public.addresses(addr_id);


--
-- Name: customers customers_sales_rep_emp_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_sales_rep_emp_id_fkey FOREIGN KEY (sales_rep_emp_id) REFERENCES public.employees(emp_id);


--
-- Name: employees employees_office_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.employees
    ADD CONSTRAINT employees_office_id_fkey FOREIGN KEY (office_id) REFERENCES public.offices(office_id);


--
-- Name: offices offices_addr_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.offices
    ADD CONSTRAINT offices_addr_id_fkey FOREIGN KEY (addr_id) REFERENCES public.addresses(addr_id);


--
-- Name: orderdetails orderdetails_ord_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orderdetails
    ADD CONSTRAINT orderdetails_ord_id_fkey FOREIGN KEY (ord_id) REFERENCES public.orders(ord_id);


--
-- Name: orderdetails orderdetails_prod_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orderdetails
    ADD CONSTRAINT orderdetails_prod_id_fkey FOREIGN KEY (prod_id) REFERENCES public.products(prod_id);


--
-- Name: orders orders_cust_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_cust_id_fkey FOREIGN KEY (cust_id) REFERENCES public.customers(cust_id);


--
-- Name: payments payments_cust_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_cust_id_fkey FOREIGN KEY (cust_id) REFERENCES public.customers(cust_id);


--
-- Name: payments payments_ord_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_ord_id_fkey FOREIGN KEY (ord_id) REFERENCES public.orders(ord_id);


--
-- Name: products products_prod_line_name_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_prod_line_name_fkey FOREIGN KEY (prod_line_name) REFERENCES public.productlines(prod_line_name);


--
-- Name: products products_prod_vendor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_prod_vendor_id_fkey FOREIGN KEY (prod_vendor_id) REFERENCES public.vendors(vendor_id);


--
-- Name: vendors vendors_addr_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_addr_id_fkey FOREIGN KEY (addr_id) REFERENCES public.addresses(addr_id);


--
-- PostgreSQL database dump complete
--

