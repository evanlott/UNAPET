// Eileen Drass
// isInstructorFunctionTesting
// This function was used to test the isInstructor function
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
//Inputs: userID and course name
//Outputs: This function returns true if a user is an instructor. It returns
//	false if a user is not an instructor. It returns an error if an error
//	occurs.
//Written By: Evan Lott
//Purpose: This function determines whether a user is an instructor of a
//	a course or not.
//---------------------------------------------------------------------------
func isInstructor(userID int, courseName string) (bool, error) {

	retVal := false

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return retVal, errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return retVal, errors.New("Failed to connect to the database.")
	}

	rows, err := db.Query("SELECT LastName FROM Users WHERE UserID =?", userID)

	if err != nil {
		return retVal, errors.New("Error retrieving instructor name.")
	}

	defer rows.Close()

	if rows.Next() == false {
		return retVal, errors.New("Query didn't match any users.")
	}

	// compare lastName to the name in courseName
	var lastName string

	rows.Scan(&lastName)

	nameSubstr := courseName[:len(lastName)]

	if lastName != nameSubstr {
		retVal = false
	} else {
		retVal = true
	}

	return retVal, nil
}

//------------------------------------------------------------------------------
func main() {
	//isInstructor(10002, "TerwilligerCS15501SP17") // true
	//isInstructor(10004, "TerwilligerCS15501SP17") // false
	//isInstructor(-1, "TerwilligerCS15501SP17") // "Query didn't match any users."
	//isInstructor(10003, "TerwilligerCS15501SP17") // false
	isInstructor(10002, "TerwilligerCS1550116SP17") //
}
