package main

import (
	"database/sql"
	"errors"

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

// returns true or false if user is enrolled in class or not
func isEnrolled(userID int, courseName string) (bool, error) {
	
	enabled := 1
	queriedCourseName := ""
	
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	
	if err != nil {
		return false, errors.New("No connection")
	}
	
	defer db.Close()
	
	rowsAffected, err := db.Query("SELECT Enabled FROM Users WHERE UserID =?", userID)
	
	if err != nil {
		return false, errors.New("Error retrieving enrollment status.")
	}
	
	if rowsAffected.Next() == false {
		return false, errors.New("Query didn't match any users.")
	}
	
	rowsAffected.Scan(&enabled)
	
	// if they are not enabled, they are not enrolled
	if enabled == 0 {
		return false, nil
	}
	
	rowsAffected, err := db.Query("SELECT CourseName FROM StudentCourses where UserID =?", userID)
	
	// user is not in any course
	if rowsAffected.Next() == false {
		return false, errors.New("Query didn't match any users.")
	}
	
	rowsAffected.scan(&queriedCourseName)
	
	if queriedCourseName != courseName {
		return false, nil
	}
	
	return true, nil
}


// returns T or F if assignment is availible or not... assignment start dateTime < time.NOW() < assignment end dateTime
func assignmentOpen(courseName string, assignmentName string) (bool, error) {
	
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return false, errors.New("No connection")
	}

	defer db.Close()

	rowsAffected, err := db.Query("SELECT StartDate, EndDate FROM Assignments WHERE AssignmentName =?", assignmentName)
	
	if err != nil {
		return false, errors.New("Error retrieving start/end date.")
	}
	
	if rowsAffected.Next() == false {
		return false, errors.New("No assignments matched with query.")
	}
	
	var startDate, endDate time.Time
	currentTime := time.Now()
	
	rowsAffected.Scan(&startDate, endDate)
	
	// TODO : figure out how to parse vars as time.Time
	
	// note: will not work because startDate/endDate have not been converted from SQL DATETIME to golang time.Time
	if startDate <= currentTime && endDate >= currentTime {
		return true, nil
	}
	
	return false, nil
}

/*
// returns T or F if course is open or not
func courseOpen(courseName string) (bool, error) {}

func changePassword(userID int, newPassword string) error {}

func getLastAssignmentname(courseName string) (string, string) {}

func getLastSubmissionName(courseName string, assignmentName, student int) (string, string) {}

func zipAssignment(courseName string, assignmentName) {}

// may or may not need this
func deleteTestCase(courseName string, assignmentName string, testCaseNum int) error {}
*/

// return a users priv level
func getPrivLevel(userID int) (int, error) {

	privLevel := -1
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return privLevel, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("SELECT privelegeLevel FROM Users WHERE UserID =?", userID)

	if err != nil {
		return privLevel, errors.New("Error retrieving privelege level.")
	}

	if rows.Next() == false {
		return privLevel, errors.New("Query didn't match any users.")
	}

	rows.Scan(&privLevel)

	return privLevel, nil
}

// returns T or F if user is instructor for the course or not
func isInstructor(userID int, courseName string) (bool, error) {

	retVal := false

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return retVal, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("SELECT firstName FROM Users WHERE UserID =?", userID)

	if err != nil {
		return retVal, errors.New("Error retrieving instructor name.")
	}

	if rows.Next() == false {
		return retVal, errors.New("Query didn't match any users.")
	}

	// compare firstName to the name in courseName
	var firstName string
	nameSubstr := courseName[:len(firstName)]

	// fill firstName variable
	rows.Scan(&firstName)

	if firstName != nameSubstr {
		retVal = false
	} else {
		retVal = true
	}

	return retVal, nil
}
