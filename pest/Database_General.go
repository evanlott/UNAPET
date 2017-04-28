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

//---------------------------------------------------------------------------
//Inputs: candidate
//Outputs: returns the nullable value
//Written By: Nathan Huckaba
//Purpose: This function handles nullable values so they are handled safely 
//	in the database. 
//---------------------------------------------------------------------------
func mayBeNull(candidate int) sql.NullInt64 {

	if candidate == 0 {
		return sql.NullInt64{}
	}

	retVal := sql.NullInt64{
		Int64: int64(candidate),
		Valid: true,
	}

	return retVal
}


//---------------------------------------------------------------------------
//Inputs: userID and courseName
//Outputs: This function returns true if a user is enrolled in a class. It 
//	returns false if a user is not enrolled in a class. It returns an
// 	error if an error occurs. 
//Written By: Evan Lott
//Purpose: This function determines whether or not a user is enrolled in
//	class. 
//---------------------------------------------------------------------------
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
//---------------------------------------------------------------------------
//Inputs: courseName and assignmentName
//Outputs: This function returns true if an assignment is available. It 
//	returns false if an assignment is not available. It returns an
// 	error if an error occurs. 
//Written By: Evan Lott and Eileen Drass
//Purpose: This function determines whether an assignment is available or 
// or not. 
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: courseName 
//Outputs: This function returns true if a course is open. It returns
//	false if a course is closed. It returns an error if an error occurs.
//Written By: Evan Lott and Eileen Drass
//Purpose: This function determines whether a course is open or not. 
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: student, courseName, assignmentName, submission number
//Outputs: This function returns an error if an error occurs. 
//Written By: Nathan Huckaba
//Purpose: This function inserts a submission into the Submissions table. 
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: userID
//Outputs: This function returns an error if an error occurs. 
//Written By: Nathan Huckaba
//Purpose: This function selects a user from the Users table. 
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: courseName
//Outputs: This function returns an error if an error occurs. 
//Written By: Nathan Huckaba
//Purpose: This function selects the last assignment. 
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: courseName, assignmentName, student
//Outputs: This function returns the last submission number. It returns an
//	error if an error occurs. 
//Written By: Nathan Huckaba
//Purpose: This function selects a student's most recent submission. 
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: userID
//Outputs: This function returns a user's privilege level. It returns an
//	error if an error occurs. 
//Written By: Evan Lott
//Purpose: This function gets a  user's privilege level. 
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: userID and courseName
//Outputs: This function returns true if a user is an instructor. It returns
//	false if a user is not an instructor. It returns an error if an error
//	occurs. 
//Written By: Evan Lott
//Purpose: This function determines whether a user is an instructor or not. 
//---------------------------------------------------------------------------
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
