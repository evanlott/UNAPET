package main

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: test this
//---------------------------------------------------------------------------
//Inputs: course name, assignment display name, assignment name, runtime,
//	number of test cases, compiler options, start date of assignment
//	end date of assignment
//Outputs: returns errors if it failed to create an assignment
//Written By: Evan Lott
//Purpose: This function will be used by the instructors to create an
//	assignment for their class. It will add an assignment to the
//	Assignments table in the database.
//---------------------------------------------------------------------------
func createAssignment(courseName string, assignmentDisplayName string, assignmentName string, runtime int, numTestCases int, compilerOptions string, startDate string, endDate string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO `Assignments` (`CourseName`, `AssignmentDisplayName`, `AssignmentName`, `StartDate`, `EndDate`, `MaxRuntime`, `CompilerOptions`, `NumTestCases`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", courseName, assignmentDisplayName, assignmentName, startDate+" 23:59:59", endDate+" 23:59:59", runtime, compilerOptions, numTestCases)

	if err != nil {
		return errors.New("Create assignment failed. Please fill out alll fields." + courseName + assignmentDisplayName + assignmentName + compilerOptions + startDate + endDate)
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Could not create assignment.")
	}

	// panic crashing program
	/*
		// TODO : verify query on server, figure out how to pull test cases from UI and upload to server
		res, err = db.Exec("ALTER TABLE GradeReport ADD " + assignmentName + " tinyint")

		if err != nil {
			panic("Error adding assignment to GradeReport")
		}

		// need rows affected check

	*/

	return nil

}

// TODO : delete assignment's folder from disk
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

	res, err := db.Exec("delete from Assignments where (AssignmentName =? and CourseName =?)", assignmentName, courseName)

	if err != nil {
		return errors.New("Assignment failed to delete from the database.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any assignments or courses.")
	}

	assignmentFolder := "/var/www/data/" + courseName + "/" + assignmentName

	os.Remove(assignmentFolder)

	_, err = os.Stat(assignmentFolder)

	// err will be nil if the folder still exists
	if err == nil {
		return errors.New("Assignment deleted from the database, but assignment directory not removed from the disk. This folder should be manually removed before another assignment is created.")
	}

	return nil
}

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

	db.Close()

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
