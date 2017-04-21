package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

import "golang.org/x/crypto/bcrypt"

import "net/http"	

import "net/http/cgi"

//import "net/http/httputil"

import "fmt"

import "log"

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

var db *sql.DB
var err error



/*func signupPage(res http.ResponseWriter, req *http.Request) {

	fmt.Println("Hello")
	//req, err = cgi.Request()
	//if err != nil {
	//	return
	//}
	/*if req.Method != "POST" {
		http.ServeFile(res, req, "signup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	var user string

	err := db.QueryRow("SELECT Username FROM Users WHERE Username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}
		_, err = db.Exec("INSERT INTO Users(Username, Password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
} */

func errorResponse(code int, msg string) {
    fmt.Printf("Status:%d %s\r\n", code, msg)
    fmt.Printf("Content-Type: text/plain\r\n")
    fmt.Printf("\r\n")
    fmt.Printf("%s\r\n", msg)

}

func check(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func main() {

	//var req *http.Request
	//var err error
	//req, err = cgi.Request()
	//if err != nil {
	//	errorResponse(500, "cannot get cgi request" + err.Error())
	//	return
	//}
	//fmt.Printf("Content-Type: text/HTML\r\n")
    	//fmt.Printf("\r\n")
	//requestDump, err := httputil.DumpRequestOut(req, true)

   	//if err != nil {
     	//fmt.Println(err)
   	//}
    //fmt.Println(string(requestDump))

	//fmt.Println("Hello")
	//return

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	
	
	if err := cgi.Serve(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		header := res.Header()
		header.Set("Content-Type", "text/plain; charset=utf-8")
		
		if req.Method != "POST" {
			http.ServeFile(res, req, "signup.html")
			return
		}

		username := req.FormValue("username")
		password := req.FormValue("password")

		var user string

		err := db.QueryRow("SELECT Username FROM Users WHERE Username=?", username).Scan(&user)

		switch {
		case err == sql.ErrNoRows:
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.", 500)
				return
			}
			_, err = db.Exec("INSERT INTO Users(Username, Password) VALUES(?, ?)", username, hashedPassword)
			if err != nil {
				http.Error(res, "Server error, unable to create your account.", 500)
				return
			}

			res.Write([]byte("User created!"))
			return
		case err != nil:
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		default:
			http.Redirect(res, req, "/", 301)
		}
	})); err != nil {
		fmt.Println(err)
	}
	//http.HandleFunc("/signup", signupPage)
	//http.HandlerFunc(signupPage)
	//check(err, "cannot serve request")
	//http.ListenAndServe(":8080", nil)
}
