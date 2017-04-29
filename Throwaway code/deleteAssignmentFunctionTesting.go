// Eileen Drass
// deleteAssingmentFunctionTesting
// This file was used to test the deleteAssignment function.

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const SHELL_NAME string = "ksh"
const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

//---------------------------------------------------------------------------
//Inputs: course name, assignment name
//Outputs: returns errors if the assignment failed to delete
//Written By: Hannah Hopkins
//Purpose: This function will be used by the instructors to delete an
//	assignment for their class. It will remove an assignment from
//	the Assignments table in the database.
//---------------------------------------------------------------------------
func deleteAssignment(courseName string, assignmentName string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return errors.New("Failed to connect to the database.")
	}

	res, err := db.Exec("delete from Assignments where (AssignmentName =? and courseName =?)", assignmentName, courseName)

	if err != nil {
		fmt.Println("Assignment failed to delete from the database.")
		return errors.New("Assignment failed to delete from the database.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("Query didn't match any assignments or courses.")
		return errors.New("Query didn't match any assignments or courses.")
	}

	// Changed the directory for testing purposes.
	assignmentFolder := "/home/eileen/testing/" + courseName + "/" + assignmentName

	os.Remove(assignmentFolder)

	_, err = os.Stat(assignmentFolder)

	// err will be nil if the folder still exists
	if err == nil {
		fmt.Println("Assignment deleted from the database, but assignment directory not removed from the disk. This folder should be manually removed before another assignment is created.")
		return errors.New("Assignment deleted from the database, but assignment directory not removed from the disk. This folder should be manually removed before another assignment is created.")
	}

	return nil
}

//---------------------------------------------------------------------------
func main() {

	deleteAssignment("TerwilligerCS15501SP17", "3")

}
