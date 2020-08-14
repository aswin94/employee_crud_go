package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ID         int
	Name       string
	Email      string
	Education  string
	Experience string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "admin12345!"
	dbName := "employees"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM emps ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		var id int
		var name, emailAddress, education, experience string
		err = selDB.Scan(&id, &name, &emailAddress, &education, &experience)
		if err != nil {
			panic(err.Error())
		}
		emp.ID = id
		emp.Name = name
		emp.Email = emailAddress
		emp.Education = education
		emp.Experience = experience
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nID := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM emps WHERE id=?", nID)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, email, education, experience string
		err = selDB.Scan(&id, &name, &email, &education, &experience)
		if err != nil {
			panic(err.Error())
		}
		emp.ID = id
		emp.Name = name
		emp.Email = email
		emp.Education = education
		emp.Experience = experience
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nID := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM emps WHERE id=?", nID)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for selDB.Next() {
		var id int
		var name, email, education, experience string
		err = selDB.Scan(&id, &name, &email, &education, &experience)
		if err != nil {
			panic(err.Error())
		}
		emp.ID = id
		emp.Name = name
		emp.Email = email
		emp.Education = education
		emp.Experience = experience
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		education := r.FormValue("education")
		experience := r.FormValue("experience")
		insForm, err := db.Prepare("INSERT INTO emps(name, email, education, experience) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email, education, experience)
		log.Println("INSERT: Name: " + name + " | Email: " + email)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		education := r.FormValue("education")
		experience := r.FormValue("experience")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE emps SET name=?, email=?, education=?, experience=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email, education, experience, id)
		log.Println("UPDATE: Name: " + name + " | Email: " + email)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM emps WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
