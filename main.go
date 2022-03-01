package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

type Event struct {
	Id       int
	Name     string
	Location string
}

const (
	host     = "ec2-3-89-214-80.compute-1.amazonaws.com"
	port     = 5432
	user     = "zcmuimobqcsutt"
	password = "424cfd55622a152e7264132b96ca12b3897f1615bdbab13e013543e4ea3e421e"
	dbname   = "d7b94fk31tosjs"
)

func dbConn() (db *sql.DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Event ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	emp := Event{}
	res := []Event{}
	for selDB.Next() {
		var id int
		var name, location string
		err = selDB.Scan(&id, &name, &location)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Location = location
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Event WHERE id=?", nId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	emp := Event{}
	for selDB.Next() {
		var id int
		var name, location string
		err = selDB.Scan(&id, &name, &location)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Location = location
	}
	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM Event WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := Event{}
	for selDB.Next() {
		var id int
		var name, location string
		err = selDB.Scan(&id, &name, &location)
		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.Name = name
		emp.Location = location
	}
	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		location := r.FormValue("location")
		insForm, err := db.Prepare("INSERT INTO Event(name, location) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, location)
		log.Println("INSERT: Name: " + name + " | Location: " + location)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		location := r.FormValue("location")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE Event SET name=?, location=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, location, id)
		log.Println("UPDATE: Name: " + name + " | Location: " + location)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Event WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	http.Redirect(w, r, "/", 301)
	defer db.Close()
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
