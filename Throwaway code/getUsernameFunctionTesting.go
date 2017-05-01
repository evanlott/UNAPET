// Eileen Drass
// getUsernameFunctionTesting
// This file was used to test the getUserName function
package main

import (
	"database/sql"
	"errors"
	"fmt"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

const PRIV_DISABLED = 0
const PRIV_STUDENT = 1
const PRIV_SI = 5
const PRIV_INSTRUCTOR = 10
const PRIV_ADMIN = 15

//---------------------------------------------------------------------------
//Inputs: userID
//Outputs: This function returns an error if an error occurs.
//Written By: Nathan Huckaba
//Purpose: This function selects a user from the Users table.
//---------------------------------------------------------------------------
func getUserName(userID int) (string, error) {

	var userName string

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		fmt.Println(userName + " No connection.")
		return userName, errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return userName, errors.New("Failed to connect to the database.")
	}

	rows, err := db.Query("select UserName from Users where UserID=?", userID)

	if err != nil {
		fmt.Println(userName + " DB error")
		return userName, errors.New("DB error")
	}

	defer rows.Close()

	if rows.Next() == false {
		fmt.Println(userName + " No user with this ID found.")
		return userName, errors.New("No user with this ID found.")
	}

	rows.Scan(&userName)

	fmt.Println(userName)
	return userName, nil

}

//---------------------------------------------------------------------------
func main() {
	getUserName(10204)
}
