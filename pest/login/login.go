package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cgi"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

// Nathan
func isLoggedIn(sessionID string) (bool, string, error) {

	var name string
	var id string
	var expires string

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return false, id, err
	}

	err = db.QueryRow("SELECT UserName, SessionID, Expires FROM ActiveSessions WHERE SessionID=?", sessionID).Scan(&name, &id, &expires)

	if err != nil {
		return false, id, errors.New("You are not currently logged in.")
	}

	if sessionID != id {
		return false, id, errors.New("Session token mismatch.")
	}

	// compare time.NOW() with session expire date time

	return true, name, nil
}

// Nathan, Abdullah, Brad
func login(userName string, password string, res http.ResponseWriter, req *http.Request) error {

	// check if num login attempts > max attempts alloed
	// if it is, send a random password

	var databasePassword string

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT Password FROM Users WHERE Username=?", userName).Scan(&databasePassword)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))

	if err != nil {
		// increment num login attempts
		return errors.New("The password entered is incorrect.")
	}

	minutes := 5

	expiration := time.Now().Local().Add(time.Duration(minutes) * time.Second)

	rand.Seed(time.Now().UTC().UnixNano())

	randNum := rand.Intn(100000000)

	sessionID, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(randNum)), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	loginCookie := http.Cookie{Name: "sessionID", Value: string(sessionID[:])} //, Expires: expiration}

	http.SetCookie(res, &loginCookie)

	result, err := db.Exec("INSERT INTO ActiveSessions VALUES (?, ?, ?)", sessionID, userName, expiration)

	if err != nil {
		return errors.New("Error starting session.")
	}

	rowsAffected, err := result.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Sessions start failed.")
	}

	return nil
}

// Nathan
func logout(userName string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	_, err = db.Exec("delete from ActiveSessions where UserName=?", userName)

	if err != nil {
		return errors.New("Course failed to delete from the database.")
	}

	return nil
}

func errorResponse(msg string) {
	fmt.Printf("Status: 500 Bad\r\n")
	fmt.Printf("Content-Type: text/plain\r\n")
	fmt.Printf("\r\n")
	fmt.Printf("%s\r\n", msg)
	fmt.Printf("If problem persists, please contact system admin.\r\n")
}

// Nathan
func main() {

	if err := cgi.Serve(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		header := res.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		if req.Method != "POST" {
			http.Redirect(res, req, "/Nathan/login.html", 301)
			return
		}

		username := req.FormValue("username")
		password := req.FormValue("password")

		cookie, err := req.Cookie("sessionID")

		if err != nil {
			res.Write([]byte(err.Error()))
		}

		id := cookie.Value

		//fmt.Fprint(w, cookie)

		if req.FormValue("action") == "login" {
			err = login(username, password, res, req)
		} else {
			auth, _, _ := isLoggedIn(id)
			if auth == true {
				res.Write([]byte("is logged in"))
			} else {
				res.Write([]byte("not logged in"))
			}

		}

		if err != nil {
			res.Write([]byte(err.Error()))
		} else {
			res.Write([]byte("yes"))
			//http.Redirect(res, req, "/Nathan/login.html", 301)
		}

	})); err != nil {
		errorResponse("Server error occurred.")
	}

}
