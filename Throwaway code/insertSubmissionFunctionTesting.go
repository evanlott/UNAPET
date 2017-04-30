// Eileen Drass
// isEnabledFunctionTesting
// This file was used to test the isEnabled function.
package main

import (
	"database/sql"
	"errors"
	//"time"
	"fmt"

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
//Inputs: student, course name, assignment name, submission number
//Outputs: This function returns an error if an error occurs.
//Written By: Nathan Huckaba
//Purpose: This function inserts a submission into the Submissions table.
//---------------------------------------------------------------------------
func insertSubmission(student int, courseName string, assignmentName string, subNum int) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println("Failed to connect to the database.")
		return errors.New("Failed to connect to the database.")
	}

	res, err := db.Exec("INSERT INTO `Submissions` (`courseName`, `AssignmentName`, `Student`, `Grade`, `comment`, `Compile`, `Results`, `SubmissionNumber`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", courseName, assignmentName, student, nil, nil, nil, nil, subNum)

	if err != nil {
		fmt.Println("Could not insert submission into database.")
		return errors.New("Could not insert submission into database.")
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("DB insert failure.")
		return errors.New("DB insert failure")
	}

	return nil

}

//------------------------------------------------------------------------------
func main() {

	insertSubmission(11234, "TerwilligerCS15501SP17", "0", 3)
}
