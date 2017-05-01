// Eileen Drass
// getLastAssignmentFunctionTesting
// This file was used to test the getLastAssignment function
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
//Inputs: course name
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

	err = db.Ping()

	if err != nil {
		return name, errors.New("Failed to connect to the database.")
	}

	rows, err := db.Query("select AssignmentName from Assignment where courseName=? order by AssignmentName DESC limit 1", courseName)

	if err != nil {
		fmt.Println("DB error.")
		return name, errors.New("DB error")
	}

	defer rows.Close()

	if rows.Next() == false {
		fmt.Println("No assignment names matched with query.")
		return name, errors.New("No assignment names matched with query.")
	}

	rows.Scan(&name)

	fmt.Println(name)
	return name, nil

}

//---------------------------------------------------------------------------
func main() {
	getLastAssignmentName("TerwilligerCS15501SP17")
}
