// Eileen Drass
// getLastSubmissionNum
// This function was used to test the getLastSubmissionNum function
package main

import (
	"database/sql"
	"errors"
	"fmt"
	//"time"

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
//Inputs: course name, assignment name, student
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

	err = db.Ping()

	if err != nil {
		fmt.Println("Failed to connect to the database.")
		return lastSubNum, errors.New("Failed to connect to the database.")
	}

	rows, err := db.Query("select SubmissionNumber from Submissions where Student=? and courseName=? and AssignmentName=? order by SubmissionNumber DESC limit 1", student, courseName, assignmentName)

	if err != nil {
		fmt.Println("DB error")
		return lastSubNum, errors.New("DB error")
	}

	defer rows.Close()

	if rows.Next() == false {
		fmt.Println("No submission for user.")
		return lastSubNum, errors.New("No submission for user.")
	}

	rows.Scan(&lastSubNum)

	fmt.Println(lastSubNum)
	return lastSubNum, nil
}

//---------------------------------------
func main() {
	getLastSubmissionNum("TerwilligerCS15501SP17", "0", 10102)
}
