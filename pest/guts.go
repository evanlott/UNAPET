// General TODO's:
//		Maybe pass result struct around as a pointer, more speed, less memory
//		Add some kind of erorr logging, for when true errors occur, i.e. directory not existing that's suppossed to exist. Also need to notify admin
//		UI needs to verify that compiler options are valid before accepting them

package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

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

const SHELL_NAME string = "ksh"

func printResults(results Submission) {
	fmt.Println()
	fmt.Println("------------------------------------------")
	fmt.Println("Course: " + results.courseName)
	fmt.Println("Assignment: " + results.assignmentName)
	fmt.Println("User: " + results.userName)
	fmt.Println("Submission number: " + strconv.Itoa(results.submissionNum))
	fmt.Println("Source file name (no .cpp): " + results.sourceName)
	fmt.Println("Num test cases: " + strconv.Itoa(results.numTestCases))
	if results.compiled {
		fmt.Println("Compiled: true")
	} else {
		fmt.Println("Compiled: false")
	}
	fmt.Println("Runtime limit (ms): " + strconv.Itoa(results.maxRuntime))
	fmt.Println("Compiler options: " + results.compilerOptions)
	fmt.Println("Results:")
	fmt.Print(results.results)
	fmt.Println("------------------------------------------")
	fmt.Println()
}

//---------------------------------------------------------------------------
//Inputs: course name, assignment name, username
//Outputs: Errors will be returned if necessary
//Written By: Tyler Delano, Eileen Drass, Hannah Hopkins, Nathan Huckaba
//	Evan Lott
//Purpose: This function will be called when the student uploads their source
//	code. It will call the other functions that will compile, execute,
//	compare test case output to student's output, and store the
//	results in the database.
//---------------------------------------------------------------------------
func evaluate(courseName string, assignmentName string, userName string) error {

	// build a results struct
	results := Submission{}

	// open database connection
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	// get userID from database
	rows, err := db.Query("select UserID from Users where Username=?", userName)

	if err != nil {
		return errors.New("DB error")
	}

	defer rows.Close()

	if rows.Next() == false {
		return errors.New("Invalid user.")
	} else {
		rows.Scan(&results.userID)
	}

	// get number of test cases, compiler options, and maxRuntime from DB
	rows, err = db.Query("select NumTestCases, MaxRuntime, CompilerOptions from Assignments where CourseName=? and AssignmentName=?", courseName, assignmentName)
	if err != nil {
		return errors.New("DB error")
	}

	if rows.Next() == false {
		return errors.New("Invalid assignment.")
	} else {
		rows.Scan(&results.numTestCases, &results.maxRuntime, &results.compilerOptions)
		rows.Close()
	}

	results.submissionNum, err = getLastSubmissionNum(courseName, assignmentName, results.userID)

	if err != nil {
		return err
	}

	results.courseName = courseName
	results.assignmentName = assignmentName
	results.userName = userName
	results.sourceName = userName + strconv.Itoa(results.submissionNum)
	results.compiled = true
	results.results = ""

	// path to directory containing all submissions for this assignment
	assignmentRoot := "data/" + results.courseName + "/" + results.assignmentName + "/"

	// make sure assignmentRoot folder exists, if not, panic
	_, err = os.Stat(assignmentRoot)

	if err != nil {
		return errors.New("Assignment folder does not exist or permission error: " + assignmentRoot)
	}

	// call to change working directory back to current, defer to ensure execution despite errors
	currentDir, err := os.Getwd()

	defer os.Chdir(currentDir)

	// change working directory to the correct assignment folder
	os.Chdir(assignmentRoot)

	// call compile
	return compile(results)
}

//---------------------------------------------------------------------------
//Inputs: the Submission struct
//Outputs: If the program does not compile, then an error message is returned
//Written By: Tyler Delano, Eileen Drass, Hannah Hopkins, Nathan Huckaba
//	Evan Lott
//Purpose: This function will be used to compile the source code that the
//	student uploads. If it does not compile, it will return an error.
//	If it does compile, the Submission struct will be updated.
//---------------------------------------------------------------------------
func compile(results Submission) error {

	// working directory is still the assignmentRoot

	sourceName := results.sourceName + ".cpp"

	outputName := results.sourceName + "Out"

	// make sure source file exists, if not, panic
	_, err := os.Stat(sourceName)

	if err != nil {
		return errors.New("Source file does not exist or permission eror: " + sourceName)
	}

	// verified that source exists, try to compile it
	compileCmd := exec.Command(SHELL_NAME, "-c", "g++ "+results.compilerOptions+" "+sourceName+" -o "+outputName)

	compileResults, err := compileCmd.CombinedOutput()

	// useless to check err of compile command, will be set even if simply errors in source
	// instead, check if output file was produced
	_, err = os.Stat(outputName)

	if err != nil {
		// did not compile.. or permission error -_-
		results.compiled = false

		// compiler output is in compileResults
		results.results = string(compileResults)
	}

	if results.compiled {
		// call execute
		return execute(results)
	} else {
		// store results with the failed output
		return storeResults(results)
	}
}

//---------------------------------------------------------------------------
//Inputs: the Submission struct
//Outputs: It returns a function call to the store results function so that
//	the results will be stored in the database.
//Written By: Tyler Delano, Eileen Drass, Hannah Hopkins, Nathan Huckaba
//	Evan Lott
//Purpose: This function runs the program and calls the compare output
//	function to determine if the output from the student's program
//	is equivalent to the desired output from each test case. It also
//	determines if the runtime constraint was met. Last, it calls
//	the store results function to store the results in the database.
//---------------------------------------------------------------------------
func execute(results Submission) error {
	// here, program is compiled and sitting on disk
	// working directory is still the assignmentRoot

	// TODO : fix reading test cases concurrently, or is this even a problem?

	programName := results.sourceName + "Out"

	// defer deletion of the executable to ensure it happens
	defer os.Remove(programName)

	testCaseNum := results.numTestCases

	// execute program on each test case
	for i := 0; i < testCaseNum; i++ {

		programOutputFile := programName + strconv.Itoa(i) + ".txt"

		// delete the output .txt file after function finishes
		defer os.Remove(programOutputFile)

		// make sure the test case exists, test cases may have been deleted, account for this
		_, err := os.Stat("test" + strconv.Itoa(i) + ".txt")
		if err != nil {
			testCaseNum++
			continue
		}

		// redirect stdin and stdout:
		//./programName < [input file].txt 1> [output file].txt
		fullCommand := "./" + programName + " < " + "test" + strconv.Itoa(i) + ".txt" + " 1> " + programOutputFile

		var capturedStdError bytes.Buffer

		// this blocks until program exits, well one of these below does, not sure which one, doesn't matter
		execCmd := exec.Command(SHELL_NAME, "-c", fullCommand)

		execCmd.Stderr = &capturedStdError

		execCmd.Start()

		done := make(chan error, 1)

		go func() {
			done <- execCmd.Wait()
		}()

		overRun := false

		// run time constraint implementation
		select {
		case <-time.After(time.Duration(results.maxRuntime) * time.Millisecond):
			{
				fmt.Println("About to kill")
				err := execCmd.Process.Kill()
				fmt.Println("Killed")
				if err != nil {
					return errors.New("Not able to kill process.")
				}
				overRun = true
			}
		case err := <-done:
			{

				// TODO : student programs must return 0 otherwise our program will think they crashed
				// check if program crashed, if not compare test case results
				if err != nil {
					results.results += "Test case " + strconv.Itoa(i) + ": Program crashed. Error: "

					// check if OS printed anything to standard error
					if capturedStdError.String() != "" {
						results.results += capturedStdError.String() + "\n"
					} else {
						results.results += "Unknown error \n"
					}
				} else {

					// call compare once for each test case
					passed := compareOutput(results, i)

					results.results += "Test case " + strconv.Itoa(i) + ": "

					if passed == true {
						results.results += "passed \n"
					} else {
						results.results += "failed \n"
					}
				}
			}
		}

		if overRun {
			results.results += "Test case " + strconv.Itoa(i) + ": Runtime limit reached. \n"
		}

	}

	return storeResults(results)
}

//---------------------------------------------------------------------------
//Inputs: the Submission struct, the test case number that is being compared
//Outputs: returns true if student's output is equivalent to the desired
//	output for that test case and returns false if the output is
//	not equivalent
//Written By: Tyler Delano, Eileen Drass, Hannah Hopkins, Nathan Huckaba
//	Evan Lott
//Purpose: This function compares the output of the student's program to the
//	the desired output from the test case and determines if they are
//	equivalent.
//---------------------------------------------------------------------------
func compareOutput(results Submission, testCaseNum int) bool {

	programName := results.sourceName + "Out"
	programOutputFile := programName + strconv.Itoa(testCaseNum) + ".txt"
	desiredOutputFile := "test" + strconv.Itoa(testCaseNum) + "DesiredOutput.txt"

	fullCommand := "diff " + programOutputFile + " " + desiredOutputFile

	// compare stuff
	compareCmd := exec.Command(SHELL_NAME, "-c", fullCommand)
	_, err := compareCmd.CombinedOutput()

	rvalue := true

	if err != nil {
		rvalue = false
	}

	return rvalue
}

//---------------------------------------------------------------------------
//Inputs: the Submission struct
//Outputs: returns an error if there is no connection, if it failed to
//	prepare, or if the update failed
//Written By: Tyler Delano, Eileen Drass, Hannah Hopkins, Nathan Huckaba
//	Evan Lott
//Purpose: This function will store the data in the Submission struct into
//	the database.
//---------------------------------------------------------------------------
func storeResults(results Submission) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	printResults(results)

	updateStatement, err := db.Prepare("update Submissions set Compile=(?), Results=(?) where Student=(?) and SubmissionNumber=(?) and AssignmentName=(?)")

	if err != nil {
		return errors.New("Failed to prepare.")
	}

	res, err := updateStatement.Exec(results.compiled, results.results, results.userID, results.submissionNum, results.assignmentName)

	if err != nil {
		return errors.New("Update failed.")
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Could not store results into database. Please try again.")
	}

	return nil
}
