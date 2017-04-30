// Eileen Drass
// editStartEndCourseFunctionTesting
// This file was used to test the editStartEndCourse function.
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
//Inputs: course name, start date for course, end date for course
//Outputs: returns errors if the start and end date for the course could not
//	be updated
//Written By: Hannah Hopkins and Nathan Huckaba
//Purpose: This function will be used by the administrator to edit the start
//	and end dates for a course. It will update the course in the
//	CourseDescription table in the database.
//---------------------------------------------------------------------------
func editStartEndCourse(courseName string, startDate string, endDate string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return errors.New("Failed to connect to the database.")
	}

	res, err := db.Exec("update CourseDescription set StartDate=?, EndDate=? where courseName=?", startDate+" 23:59:59", endDate+" 23:59:59", courseName)

	if err != nil {
		fmt.Println("Start/end update failed.")
		return errors.New("Start/end update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("Query didn't match any courses.")
		return errors.New("Query didn't match any courses.")
	}

	return nil

}

//---------------------------------------------------------------------------
func main() {

	//editStartEndCourse("TerwilligerCS15502SP17", "2017-04-25 19:14:59", "2017-05-16 00:00:02")
	editStartEndCourse("TerwilligerCS15506SP17", "2017-04-25 19:14:59", "2017-05-16 00:00:02")

}
