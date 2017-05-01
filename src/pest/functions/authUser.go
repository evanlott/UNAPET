package functions

import (
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//---------------------------------------------------------------------------
//Inputs: username
//Outputs: This returns true if a user is logged in. It returns false if a 
//	user is not logged in. It returns an error if an error occurs.
//Written By: Nathan Huckaba
//Purpose: This function will be used to determine if a user is logged in
//	or not. 
//---------------------------------------------------------------------------
func isLoggedIn(username string) (bool, error) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return false, err
	}

	rows, err := db.Query("SELECT * FROM ActiveSessions WHERE UserName=?", username)

	if err != nil {
		return false, errors.New("DB error.")
	}

	if rows.Next() != true {
		return false, errors.New("You are not currently logged in.")
	}

	// compare time.NOW() with session expire date time

	return true, nil
}

//---------------------------------------------------------------------------
//Inputs: username
//Outputs: This returns true if a user is logged in. It returns false if a 
//	user is not logged in. It returns an error if an error occurs.
//Written By: Nathan Huckaba, Abdullah Karaman, Brad Lanford
//Purpose: This function will be used to determine if a user is logged in
//	or not. 
//---------------------------------------------------------------------------
func Login(userName string, password string, res http.ResponseWriter, req *http.Request) error {

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

	// load this from config file
	minutes := 5

	expiration := time.Now().Local().Add(time.Duration(minutes) * time.Second)

	rand.Seed(time.Now().UTC().UnixNano())

	randNum := rand.Intn(100000000)

	sessionID, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(randNum)), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// TODO : make sure cookies are actually expiring
	loginCookie := http.Cookie{Name: "sessionID", Value: string(sessionID[:]), Expires: expiration}

	http.SetCookie(res, &loginCookie)

	Logout(userName)

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

//---------------------------------------------------------------------------
//Inputs: username
//Outputs: This function returns an error if an error occurs.
//Written By: Nathan Huckaba
//Purpose: This function will log a user out.
//---------------------------------------------------------------------------
func Logout(userName string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	_, err = db.Exec("delete from ActiveSessions where UserName=?", userName)

	if err != nil {
		return errors.New("Session failed to delete from the database.")
	}

	return nil
}

// Nathan
// returns T/F if logged in or not, username, and err
//---------------------------------------------------------------------------
//Inputs: http request
//Outputs: This returns a unsername. It also returns an error if an error
//	occurs. 
//Written By: Nathan Huckaba
//Purpose: This function makes sure a user is authorized.
//---------------------------------------------------------------------------
func AuthUser(req *http.Request) (string, error) {

	var userName string

	cookie, err := req.Cookie("sessionID")

	if err != nil {
		return userName, errors.New("Cookie error. Close your browser and try again.")
	}

	id := cookie.Value

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return userName, err
	}

	err = db.QueryRow("SELECT UserName FROM ActiveSessions WHERE SessionID=?", id).Scan(&userName)

	if err != nil {
		return userName, errors.New("You are not logged in.")
	}

	return userName, err

}

/*
// Nathan, Abdullah, Brad
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

		var err error

		if req.FormValue("Action") == "login" {
			err = login(username, password, res, req)
		} else {
			auth, _ := isLoggedIn(username)
			if auth == true {
				res.Write([]byte("is logged in"))
			} else {
				res.Write([]byte("not logged in"))
			}
			return
		}

		if err != nil {
			res.Write([]byte(err.Error()))
		} else {
			res.Write([]byte("Logged you in."))
		}

		// check priv level and redirect them

	})); err != nil {
		ErrorResponse("Server error occurred.")
	}

}
*/
