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

func isLoggedIn() {

}

func login(userName string, password string, res http.ResponseWriter, req *http.Request) error {

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
		return errors.New("Wrong password.")
	}

	minutes := 1

	expiration := time.Now().Local().Add(time.Duration(minutes) * time.Second)

	rand.Seed(time.Now().UTC().UnixNano())

	randNum := rand.Intn(100000000)

	sessionID, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(randNum)), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	loginCookie := http.Cookie{Name: "sessionID", Value: string(sessionID[:]), Expires: expiration}

	http.SetCookie(res, &loginCookie)

	return nil
}

func logout() {

}

func main() {
	dummy()
	return

	fmt.Printf("Status: 500 Bad\r\n")
	fmt.Printf("Content-Type: text/plain\r\n")
	fmt.Printf("\r\n")
	fmt.Printf("If problem persists, please contact system admin.\r\n")
	return
}

func dummy() {

	//fmt.Printf("Status: 500 Bad\r\n")
	//fmt.Printf("Content-Type: text/plain\r\n")
	//fmt.Printf("\r\n")
	//fmt.Printf("If problem persisthtrhehts, please contact system admin.\r\n")
	//return

	if err := cgi.Serve(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		header := res.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")
		//res.Header().Add("Content-Type", "text/html\r\n")
		res.Write([]byte("yes"))

		//http.Redirect(res, req, "/login.html", 301)

		return

		if req.Method != "POST" {
			http.Redirect(res, req, "/login.html", 301)
			return
		}

		username := req.FormValue("username")
		password := req.FormValue("password")

		err := login(username, password, res, req)

		if err != nil {
			res.Write([]byte(err.Error()))
		} else {
			res.Write([]byte("yes"))
		}

	})); err != nil {
		fmt.Println("not yes")
	}

}
