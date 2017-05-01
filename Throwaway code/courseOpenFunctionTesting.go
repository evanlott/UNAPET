// Eileen Drass
// courseOpen
// This file was used to test the courseOpen function.
package main

import (
	"database/sql"
	"errors"
	"fmt"
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
//Inputs: course name
//Outputs: This function returns true if a course is open. It returns
//	false if a course is closed. It returns an error if an error occurs.
//Written By: Evan Lott and Eileen Drass
//Purpose: This function determines whether a course is open or not.
//---------------------------------------------------------------------------
func courseOpen(courseName string) (bool, error) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME+"?parseTime=true")

	if err != nil {
		fmt.Println("False - No conncection")
		return false, errors.New("No connection")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println("False - Failed to connect to the database.")
		return false, errors.New("Failed to connect to the database.")
	}

	rows, err := db.Query("SELECT StartDate, EndDate FROM CourseDescription WHERE courseName =?", courseName)

	if err != nil {
		fmt.Println("Error retrieving start/end date.")
		return false, errors.New("Error retrieving start/end date.")
	}

	defer rows.Close()

	if rows.Next() == false {
		fmt.Println("No courses matched with query.")
		return false, errors.New("No courses matched with query.")
	}

	var startDate, endDate time.Time
	currentTime := time.Now()

	rows.Scan(&startDate, &endDate)

	fmt.Println("Start Date: " + startDate.Format("01/02/2006 15:04:05"))
	fmt.Println("End Date: " + endDate.Format("01/02/2006 15:04:05"))
	fmt.Println("Current Date: " + currentTime.Format("01/02/2006 15:04:05"))

	if startDate.Format("01/02/2006 15:04:05") <= currentTime.Format("01/02/2006 15:04:05") && endDate.Format("01/02/2006 15:04:05") >= currentTime.Format("01/02/2006 15:04:05") {
		fmt.Println("True")
		return true, nil
	}

	fmt.Println("False")
	return false, nil
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

	err = db.Ping()

	if err != nil {
		return userName, errors.New("Failed to connect to the database.")
	}

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

//------------------------------------------------------------------------------
func main() {

	courseOpen("TerwilligerCS15501SP1617")
}
