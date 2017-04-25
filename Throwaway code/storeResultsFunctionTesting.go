// storeResultsFunctionTesting
package main

import (
	"database/sql"
	"errors"
	"fmt"
	//"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Submission struct {
	courseName      string
	assignmentName  string
	userName        string
	userID          int
	sourceName      string // without .cpp, call this submission name..?
	numTestCases    int
	compiled        bool
	results         string
	maxRuntime      int
	compilerOptions string
	submissionNum   int
}

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"
const SHELL_NAME string = "ksh"

//------------------------------------------------------------------------------
func storeResults(results Submission) error {

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

	//printResults(results)

	updateStatement, err := db.Prepare("update Submissions set Compile=(?), Results=(?) where Student=(?) and SubmissionNumber=(?) and AssignmentName=(?)")

	if err != nil {
		//fmt.Println("Failed to prepare.")
		return errors.New("Failed to prepare.")
	}

	res, err := updateStatement.Exec(results.compiled, results.results, results.userID, results.submissionNum, results.assignmentName)
	//res, err := updateStatement.Exec(results.compiled, results.results, results.userID, results.submissionNum, results.assignmentName, results.numTestCases)

	if err != nil {
		//fmt.Println("Update failed.")
		return errors.New("Update failed.")
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected != 1 {
		//fmt.Println("Could not store results into database. Please try again.")
		return errors.New("Could not store results into database. Please try again.")
	}

	return nil
}

//------------------------------------------------------------------------------
func main() {

	results := Submission{}

	results.submissionNum = 0
	results.courseName = "TerwilligerCS15501SP17"
	results.assignmentName = "0"
	results.compiled = true
	results.results = "TEST: Some results..."
	results.userID = 11111

	//fmt.Println("Test...")
	storeResults(results)

}
