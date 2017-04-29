// Eileen Drass
// editStartEndAssignment.go
// This file was used to test the editStartEndAssignment function.

//package functions
package main

import (
	"database/sql"
	"errors"
	"fmt"
	//"os"

	_ "github.com/go-sql-driver/mysql"
)

const SHELL_NAME string = "ksh"
const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

//---------------------------------------------------------------------------
//Inputs: course name, assignment name, start date for the assignment,
//	end date for an assignment
//Outputs: returns errors if the start and end dates could not be updated
//Written By: Hannah Hopkins and Nathan Huckaba
//Purpose: This function will be used by the instructors to edit the
//	start and end date for an assignment. It will update the
//	Assignments table in the database.
//---------------------------------------------------------------------------
func editStartEndAssignment(courseName string, assignmentName string, startDate string, endDate string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return errors.New("Failed to connect to the database.")
	}

	res, err := db.Exec("update Assignments set StartDate=?, EndDate=? where courseName=? and AssignmentName=?", startDate+" 23:59:59", endDate+" 23:59:59", courseName, assignmentName)

	if err != nil {
		fmt.Println("Start/end update failed.")
		return errors.New("Start/end update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("Query didn't match any assignments.")
		return errors.New("Query didn't match any assignments.")
	}

	return nil
}

//---------------------------------------------------------------------------
func main() {

	editStartEndAssignment("TerwilligerCS15501SP17", "6", "2017-04-30 12:00:00", "2017-05-30 12:00:00")
}
