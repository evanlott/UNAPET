package main

import (
	"database/sql"
	"errors"
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

// returns true or false if user is enrolled in class or not
// Evan
func isEnrolled(userID int, courseName string) (bool, error) {

	enabled := 1
	queriedCourseName := ""

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return false, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("SELECT Enabled FROM Users WHERE UserID =?", userID)

	if err != nil {
		return false, errors.New("Error retrieving enrollment status.")
	}

	defer rows.Close()

	if rows.Next() == false {
		return false, errors.New("Query didn't match any users.")
	}

	rows.Scan(&enabled)

	// if they are not enabled, they are not enrolled
	if enabled == 0 {
		return false, nil
	}

	rows, err = db.Query("SELECT CourseName FROM StudentCourses where UserID =?", userID)

	// user is not in any course
	if rows.Next() == false {
		return false, errors.New("Query didn't match any users.")
	}

	rows.Scan(&queriedCourseName)

	if queriedCourseName != courseName {
		return false, nil
	}

	return true, nil
}

// returns T or F if assignment is availible or not... assignment start dateTime < time.NOW() < assignment end dateTime
// Evan, Eileen
func assignmentOpen(courseName string, assignmentName string) (bool, error) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME+"?parseTime=true")

	if err != nil {
		return false, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("SELECT StartDate, EndDate FROM Assignments WHERE AssignmentName =?", assignmentName)

	if err != nil {
		return false, errors.New("Error retrieving start/end date.")
	}

	defer rows.Close()

	if rows.Next() == false {
		return false, errors.New("No assignments matched with query.")
	}

	var startDate, endDate time.Time
	currentTime := time.Now()

	rows.Scan(&startDate, endDate)

	if startDate.Format("01/02/2006 15:04:05") <= currentTime.Format("01/02/2006 15:04:05") && endDate.Format("01/02/2006 15:04:05") >= currentTime.Format("01/02/2006 15:04:05") {
		return true, nil
	}

	return false, nil
}

// returns T or F if course is open or not
// Evan, Eileen
func courseOpen(courseName string) (bool, error) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME+"?parseTime=true")

	if err != nil {
		return false, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("SELECT StartDate, EndDate FROM CourseDescription WHERE CourseName =?", courseName)

	if err != nil {
		return false, errors.New("Error retrieving start/end date.")
	}

	defer rows.Close()

	if rows.Next() == false {
		return false, errors.New("No courses matched with query.")
	}

	var startDate, endDate time.Time
	currentTime := time.Now()

	rows.Scan(&startDate, endDate)

	if startDate.Format("01/02/2006 15:04:05") <= currentTime.Format("01/02/2006 15:04:05") && endDate.Format("01/02/2006 15:04:05") >= currentTime.Format("01/02/2006 15:04:05") {
		return true, nil
	}

	return false, nil
}

/*

func zipAssignment(courseName string, assignmentName) {}

// may or may not need this
func deleteTestCase(courseName string, assignmentName string, testCaseNum int) error {}
*/

// Nathan
func insertSubmission(student int, courseName string, assignmentName string, subNum int) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO `Submissions` (`CourseName`, `AssignmentName`, `Student`, `Grade`, `comment`, `Compile`, `Results`, `SubmissionNumber`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", courseName, assignmentName, student, nil, nil, nil, nil, subNum)

	if err != nil {
		return errors.New("Could not insert submission into database.")
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("DB insert failure")
	}

	return nil

}

// Nathan
func getUserName(userID int) (string, error) {

	var userName string

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return userName, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select UserName from Users where UserID=?", userID)

	if err != nil {
		return userName, errors.New("DB error")
	}

	defer rows.Close()

	if rows.Next() == false {
		return userName, errors.New("No user with this ID found.")
	}

	rows.Scan(&userName)

	return userName, nil

}

// Nathan
func getLastAssignmentName(courseName string) (string, error) {

	name := "-1"

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return name, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select AssignmentName from Assignments where CourseName=? order by AssignmentName DESC limit 1", courseName)

	if err != nil {
		return name, errors.New("DB error")
	}

	defer rows.Close()

	if rows.Next() == false {
		return name, errors.New("No submission for user.")
	}

	rows.Scan(&name)

	return name, nil

}

// Nathan
func getLastSubmissionNum(courseName string, assignmentName string, student int) (int, error) {

	lastSubNum := -1

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return lastSubNum, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select SubmissionNumber from Submissions where Student=? and CourseName=? and AssignmentName=? order by SubmissionNumber DESC limit 1", student, courseName, assignmentName)

	if err != nil {
		return lastSubNum, errors.New("DB error")
	}

	defer rows.Close()

	if rows.Next() == false {
		return lastSubNum, errors.New("No submission for user.")
	}

	rows.Scan(&lastSubNum)

	return lastSubNum, nil
}

// return a users priv level
// Evan
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

	defer rows.Close()

	if rows.Next() == false {
		return privLevel, errors.New("Query didn't match any users.")
	}

	rows.Scan(&privLevel)

	return privLevel, nil
}

// returns T or F if user is instructor for the course or not
// Evan
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

	defer rows.Close()

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
