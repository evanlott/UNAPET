package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

/*


 */

// TODO
// test this
func createAssignment(courseName string, assignmentDisplayName string, assignmentName string, runtime int, numTestCases int, compilerOptions string, startDate string, endDate string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("INSERT INTO `Assignments` (`CourseName`, `AssignmentDisplayName`, `AssignmentName`, `StartDate`, `EndDate`, `MaxRuntime`, `CompilerOptions`, `NumTestCases`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", courseName, assignmentDisplayName, assignmentName, startDate+" 23:59:59", endDate+" 23:59:59", runtime, compilerOptions, numTestCases)

	if err != nil {
		return errors.New("Create assignment failed. Please fill out all fields.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Could not create assignment.")
	}

	return nil

}

/*


 */

// TODO : delete assignment's folder from disk
func deleteAssignment(courseName string, assignmentName string) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("delete from Assignments where (AssignmentName =? and CourseName =?)", assignmentName, courseName)

	if err != nil {
		return errors.New("Assignment failed to delete.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any assignments or courses.")
	}

	return nil
}

/*


 */

func editStartEndAssignment(courseName string, assignmentName string, startDate string, endDate string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("update Assignments set StartDate=?, EndDate=? where CourseName=? and AssignmentName=?", startDate+" 23:59:59", endDate+" 23:59:59", courseName, assignmentName)

	if err != nil {
		return errors.New("Start/end update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any assignments.")
	}

	return nil
}
