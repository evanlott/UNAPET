// Eileen Drass
// editCourseDescriptionFunction
// This file was used to test the editCourseDescription function.
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
//Inputs: course name, course description
//Outputs: returns errors if the course description could not be updated
//Written By: Eileen Drass and Evan Lott
//Purpose: This function will be used by the instructor to edit the course
//	description for a course. It will update the course in the
//	in the CourseDescription table in the database.
//---------------------------------------------------------------------------
func editCourseDescription(courseName string, courseDescription string) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return errors.New("Failed to connect to the database.")
	}

	res, err := db.Exec("UPDATE CourseDescription SET CourseDescription.CourseDescription=? WHERE CourseDescription.courseName=? ", courseDescription, courseName)

	if err != nil {
		fmt.Println("Update failed.")
		return errors.New("Update failed.")
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

	//createCourse("TerwilligerCS15502SP17", "Computer Science I", "TEST: This is a tearable class.", 10002, "2017-04-24 19:14:59", "2017-05-13 00:00:00", 10006, 10009, false, false)
	//deleteCourse("TerwilligerCS15502SP17")
	//editCourseDescription("TerwilligerCS15506SP17", "TEST: This is a not-so tearable class.")
	editCourseDescription("TerwilligerCS15506SP17", "TESTS: This is a not-so tearable class.")

}
