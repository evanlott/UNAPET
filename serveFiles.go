package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello from GO!") // write data to response
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("www/login.html")
		t.Execute(w, nil)
	} else {
		// open database connection
		db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
		if err != nil {
			panic("No connection")
		}

		defer db.Close()

		r.ParseForm()

		username := r.Form["username"]
		password := r.Form["password"]

		// TODO : maybe change this to prepare to "prevent SQL injection"
		// check password
		rows, err := db.Query("select Username from Users where Username=? and Password=?", username[0], password[0])
		if err != nil {
			panic("DB error")
		}

		var result string

		// TODO : query returning no results causes rows.Next() to throw segmentation fault
		rows.Next()
		rows.Scan(&result)

		if result == username[0] {
			fmt.Fprintf(w, "Login good.")
		} else {
			fmt.Fprintf(w, "Login bad.")
		}
	}
}

func main() {
	http.HandleFunc("/", sayhelloName) // setting router rule
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
