package main

import (
	"database/sql"
	"net/http"
	//"log"
	"fmt"
	"text/template"
	_"github.com/go-sql-driver/mysql"
)

func dbConnect()(conexion *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := ""
	DbName := "godb"

	conexion, err := sql.Open(Driver, User + ":" + Password + "@tcp(127.0.0.1)/" + DbName)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	fmt.Println("Server runing...")
	http.ListenAndServe(":8087", nil)
}

type Employer struct {
	Id int
	Name string
	Email string
}

func Home(w http.ResponseWriter, r *http.Request) {

	isConected := dbConnect()
	rows, err := isConected.Query("SELECT * FROM employer")
	if err != nil {
		panic(err.Error())
	}

	employer := Employer{}
	arrayEmployer := []Employer{}

	for rows.Next() {
		var id int
		var name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		employer.Id = id
		employer.Name = name
		employer.Email = email

		arrayEmployer = append(arrayEmployer, employer)
	}

	templates.ExecuteTemplate(w, "home", arrayEmployer)
}

func Create(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		isConected := dbConnect()
		insertRows, err := isConected.Prepare("INSERT INTO employer(name, email) VALUES(?, ?)")

		if err != nil { panic(err.Error()) }
		insertRows.Exec(name, email)

		http.Redirect(w, r, "/", 301)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idEmployer := r.URL.Query().Get("id")

	isConected := dbConnect()
	row, err := isConected.Query("SELECT * FROM employer WHERE id = ?", idEmployer)

	employer := Employer{}
	for row.Next() {
		var id int
		var name, email string
		err = row.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		employer.Id = id
		employer.Name = name
		employer.Email = email
	}

	templates.ExecuteTemplate(w, "edit", employer)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		isConected := dbConnect()
		updateRow, err := isConected.Prepare("UPDATE employer SET name = ?, email = ? WHERE id = ?")

		if err != nil { panic(err.Error()) }
		updateRow.Exec(name, email, id)

		http.Redirect(w, r, "/", 301)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idEmployer := r.URL.Query().Get("id")

	isConected := dbConnect()
	deleteRow, err := isConected.Prepare("DELETE FROM employer WHERE id = ?")

	if err != nil { panic(err.Error()) }

	deleteRow.Exec(idEmployer)
	http.Redirect(w, r, "/", 301)
}