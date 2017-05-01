// Eileen Drass
// assignmentOpen
// This file was used to test the assignmentOpen function.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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
//Inputs: course name and assignment name
//Outputs: This function returns true if an assignment is available. It
//	returns false if an assignment is not available. It returns an
// 	error if an error occurs.
//Written By: Evan Lott and Eileen Drass
//Purpose: This function determines whether an assignment is available or
// or not.
//---------------------------------------------------------------------------
func assignmentOpen(courseName string, assignmentName string) (bool, error) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME+"?parseTime=true")

	if err != nil {
		return false, errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println("Failed to connect to the database.")
		return false, errors.New("Failed to connect to the database.")
	}

	rows, err := db.Query("SELECT StartDate, EndDate FROM Assignments WHERE AssignmentName =?", assignmentName)

	if err != nil {
		fmt.Println("False - Error retrieving start/end date.")
		return false, errors.New("Error retrieving start/end date.")
	}

	defer rows.Close()

	if rows.Next() == false {
		fmt.Println("False - No assignments matched with query.")
		return false, errors.New("No assignments matched with query.")
	}

	var startDate, endDate time.Time
	currentTime := time.Now()

	rows.Scan(&startDate, &endDate)

	fmt.Println("Start Date: " + startDate.Format("01/02/2006 15:04:05"))
	fmt.Println("End Date: " + endDate.Format("01/02/2006 15:04:05"))
	fmt.Println("Current Date: " + currentTime.Format("01/02/2006 15:04:05"))
	fmt.Println("False")

	if startDate.Format("01/02/2006 15:04:05") <= currentTime.Format("01/02/2006 15:04:05") && endDate.Format("01/02/2006 15:04:05") >= currentTime.Format("01/02/2006 15:04:05") {
		fmt.Println("True")
		return true, nil
	}

	return false, nil
}

//------------------------------------------------------------------------------
func main() {

	assignmentOpen("TerwilligerCS15501SP17", "0")
}
