// Eileen Drass
// databaseAssignmentsTesting.go
// This file was used to test the createAssignment function.

//package functions
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

	err = db.Ping()

	if err != nil {
		fmt.Println("Failed to connect to the database.")
		return errors.New("Failed to connect to the database.")
	}

	if compilerOptions == "" {
		compilerOptions = "NULL"
		fmt.Println(compilerOptions)
	}

	res, err := db.Exec("INSERT INTO `Assignments` (`courseName`, `AssignmentDisplayName`, `AssignmentName`, `StartDate`, `EndDate`, `MaxRuntime`, `CompilerOptions`, `NumTestCases`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", courseName, assignmentDisplayName, assignmentName, startDate+" 23:59:59", endDate+" 23:59:59", runtime, compilerOptions, numTestCases)

	if err != nil {
		fmt.Println("Create assignment failed. Please do not create a duplicate assignment, and please fill out all fields." + courseName + assignmentDisplayName + assignmentName + compilerOptions + startDate + endDate)
		return errors.New("Create assignment failed. Please fill out all fields." + courseName + assignmentDisplayName + assignmentName + compilerOptions + startDate + endDate)
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("Could not create assignment.")
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

//---------------------------------------------------------------------------
func main() {

	createAssignment("TerwilligerCS15501SP17", "TEST: Assignment 5", "5", 10000, 1, "-Wall", "2017-04-25 15:00:00", "2017-05-25 15:00:00")
	//createAssignment("TerwilligerCS15501SP17", "TEST: Assignment 3", "3", 10000, 1, "-Wall", "2017-04-25 15:00:00", "2017-05-25 15:00:00"")
	//createAssignment("TerwilligerCS15501SP17", "TEST: Assignment 4", "4", 10000, 1, "", "2017-04-25 15:00:00", "2017-05-25 15:00:00")
}
