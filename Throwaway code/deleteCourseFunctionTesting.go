// Eileen Drass
// deleteCourseFunctionTesting
// This file was used to test the deleteCourse function.
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
//Inputs: course name
//Outputs: returns errors if the course fails to delete
//Written By: Hannah Hopkins
//Purpose: This function will be used by the administrator to delete a
//	course. It will remove the course from the CourseDescription
//	table in the database.
//---------------------------------------------------------------------------
func deleteCourse(courseName string) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return errors.New("Failed to connect to the database.")
	}

	res, err := db.Exec("delete from CourseDescription where courseName =?", courseName)

	if err != nil {
		fmt.Println("Course failed to delete from the database.")
		return errors.New("Course failed to delete from the database.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("Query didn't match any courses.")
		return errors.New("Query didn't match any courses.")
	}

	// This was changed for testing purposes.
	//courseFolder := "/var/www/data/" + courseName
	courseFolder := "/home/eileen/testing/deleteCourseFunction" + courseName
	err = os.RemoveAll(courseFolder)

	if err != nil {
		fmt.Println("err")
		return err
	}

	_, err = os.Stat(courseFolder)

	// err will be nil if the folder still exists
	if err == nil {
		return errors.New("Course deleted from the database, but course directory not removed from the disk. This folder should be manually removed.")
	}

	return nil
}

//---------------------------------------------------------------------------
func main() {

	//createCourse("TerwilligerCS15502SP17", "Computer Science I", "TEST: This is a tearable class.", 10002, "2017-04-24 19:14:59", "2017-05-13 00:00:00", 10006, 10009, false, false)
	deleteCourse("TerwilligerCS15502SP17")

}
