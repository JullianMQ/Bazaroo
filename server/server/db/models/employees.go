package models

import "github.com/JullianMQ/Bazaroo/server/db"

func CreateEmployeesTable() {
	_, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS employees (
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
}
