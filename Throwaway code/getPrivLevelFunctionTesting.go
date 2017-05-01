// Eileen Drass
// getPrivLevelFunctionTesting
// This function was used to test the getPrivLevel function
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
//Outputs: This function returns a user's privilege level. It returns an
//	error if an error occurs.
//Written By: Evan Lott
//Purpose: This function gets a  user's privilege level.
//---------------------------------------------------------------------------
func getPrivLevel(userID int) (int, error) {

	privLevel := -1
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return privLevel, errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return privLevel, errors.New("Failed to connect to the database.")
	}

	rows, err := db.Query("SELECT PrivLevel FROM Users WHERE UserID =?", userID)

	if err != nil {
		fmt.Println("Error retrieving privelege level.")
		return privLevel, errors.New("Error retrieving privelege level.")
	}

	defer rows.Close()

	if rows.Next() == false {
		fmt.Println("Query didn't match any users.")
		return privLevel, errors.New("Query didn't match any users.")
	}

	rows.Scan(&privLevel)

	fmt.Println(privLevel)
	return privLevel, nil
}

//--------------------------------------------------------------------------
func main() {

	getPrivLevel(10001) // 15
	getPrivLevel(10067) // 10
	getPrivLevel(10072) // 5
	getPrivLevel(10004) // 1
	getPrivLevel(-1)    // Query didn't match any users

}
