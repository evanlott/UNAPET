// General TODO's:
//		Maybe pass result struct around as a pointer, more speed, less memory
//		Add some kind of erorr logging, for when true errors occur, i.e. directory not existing that's suppossed to exist. Also need to notify admin
//		UI needs to verify that compiler options are valid before accepting them

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

/*


 */

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

/*


 */

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"
const SHELL_NAME string = "ksh"

/*


 */

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

/*


 */

// Maybe, instead of using courseName and assignmentName, make UI pass us a submission ID number..? field would need to be added to the DB
func evaluate(courseName string, assignmentName string, userName string) {

	// build a results struct
	results := Submission{}

	// open database connection
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	defer db.Close()

	// get userID from database
	rows, err := db.Query("select UserID from Users where Username=?", userName)
	if err != nil {
		panic("DB error")
	}

	if rows.Next() == false {
		panic("Invalid user.")
	} else {
		rows.Scan(&results.userID)
	}

	// get number of test cases, compiler options, and maxRuntime from DB
	rows, err = db.Query("select NumTestCases, MaxRuntime, CompilerOptions from Assignments where CourseName=? and AssignmentName=?", courseName, assignmentName)
	if err != nil {
		panic("DB error")
	}

	if rows.Next() == false {
		panic("Invalid assignment.")
	} else {
		rows.Scan(&results.numTestCases, &results.maxRuntime, &results.compilerOptions)
	}

	rows, err = db.Query("select SubmissionNumber from Submissions where Student=? and AssignmentName=? order by SubmissionNumber DESC limit 1", results.userID, assignmentName)
	if err != nil {
		panic("DB error")
	}

	if rows.Next() == false {
		panic("No submission for user.")
	} else {
		rows.Scan(&results.submissionNum)
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
		panic("Assignment folder does not exist or permission error: " + assignmentRoot)
	}

	// call to change working directory back to current, defer to ensure execution despite errors
	currentDir, err := os.Getwd()

	defer os.Chdir(currentDir)

	// change working directory to the correct assignment folder
	os.Chdir(assignmentRoot)

	// call compile
	compile(results)
}

/*


 */

// compile stuff
func compile(results Submission) {

	// working directory is still the assignmentRoot

	sourceName := results.sourceName + ".cpp"

	outputName := results.sourceName + "Out"

	// make sure source file exists, if not, panic
	_, err := os.Stat(sourceName)

	if err != nil {
		panic("Source file does not exist or permission eror: " + sourceName)
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
		execute(results)
	} else {
		// store results with the failed output
		storeResults(results)
	}
}

/*


 */

func execute(results Submission) {
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
					panic("Not able to kill process.")
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

	storeResults(results)
}

/*


 */

func compareOutput(results Submission, testCaseNum int) bool {

	programName := results.sourceName + "Out"
	programOutputFile := programName + strconv.Itoa(testCaseNum) + ".txt"
	desiredOutputFile := "test" + strconv.Itoa(testCaseNum) + "DesiredOutput.txt"

	fullCommand := "diff " + programOutputFile + " " + desiredOutputFile

	// compare stuff
	compareCmd := exec.Command(SHELL_NAME, "-c", fullCommand)
	compareResults, err := compareCmd.CombinedOutput()

	if compareResults != nil {
	}

	rvalue := true

	if err != nil {
		rvalue = false
	}

	return rvalue
}

/*


 */

func storeResults(results Submission) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	printResults(results)

	updateStatement, err := db.Prepare("update Submissions set Compile=(?), Results=(?) where Student=(?) and SubmissionNumber=(?)")

	if err != nil {
		panic("Failed to prepare.")
	}

	_, err = updateStatement.Exec(results.compiled, results.results, 10004, results.submissionNum)

	if err != nil {
		panic("Update failed.")
	} else {
		fmt.Println("Stored the results. \n")
	}
}

func importCSV(name string) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	mysql.RegisterLocalFile(name)

	// TODO : Solve password @dummy issue, also CSV quotation issue, trailing comma issue
	_, err = db.Exec("LOAD DATA LOCAL INFILE '" + name + "' INTO TABLE Users FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE 1 LINES (@dummy, FirstName, MiddleInitial, LastName, UserName, Password, @dummy, @dummy, @dummy, @dummy, @dummy)")

	if err != nil {
		panic("Import failed.")
	}
}

/*


 */

func editCourseDescription(courseName string, courseDescription string)
{
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}
	
	editStatement, err := db.Prepare("UPDATE CourseDescription SET CourseDescription.CourseDescription = " + courseName + " WHERE CourseDescription.CourseName = " + courseDescription + ")
	
	if err != nil {
		panic("Failed to prepare")
	}
	
	_, err = editStatement.Exec();
	
	if err != nil {
		panic("Update failed.")
	} else {
		fmt.Println("Updated course description table\n")
	}
	
}

/*


 */
					 
 func gradeAssignment(userID int, courseName string, assignmentName string, submissionNum int, grade int) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	res, err := db.Exec("update Submissions set grade=? where Student=? and CourseName=? and AssignmentName=? and SubmissionNumber=?", grade, userID, courseName, assignmentName, submissionNum)

	if err != nil {
		panic("Grade update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		panic("Query didn't match any submissions.")
	}

}

/*


 */

func editStartEndCourse(courseName string, startDate string, endDate string) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	res, err := db.Exec("update CourseDescription set StartDate=?, EndDate=? where CourseName=?", startDate+" 23:59:59", endDate+" 23:59:59", courseName)

	if err != nil {
		panic("Start/end update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		panic("Query didn't match any courses.")
	}

}

/*


 */

func editUser(userID int, firstName string, MI string, lastName string, privLevel int) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	res, err := db.Exec("update Users set FirstName=?, MiddleInitial=?, LastName=?, PrivLevel=? where UserID=?", firstName, MI, lastName, privLevel, userID)

	if err != nil {
		panic("User update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		panic("Query didn't match any users.")
	}

}

func editStartEndAssignment(courseName string, assignmentName string, startDate string, endDate string) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	res, err := db.Exec("update Assignments set StartDate=?, EndDate=? where CourseName=? and AssignmentName=?", startDate+" 23:59:59", endDate+" 23:59:59", courseName, assignmentName)

	if err != nil {
		panic("Start/end update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		panic("Query didn't match any assignments.")
	}

}

/*


 */

func main() {
	importCSV("Users.csv")
	//evaluate("TerwilligerCS15501SP17", "Assignment 1", "jdoe")

}

/*






















 */

/*
----------------------------------------------------------
  Throw away code
----------------------------------------------------------


	3/19

main()
{
	test := Results{}

	test.sourceName = "test"

	compile(test)
}

compile()
{
	//	if err != nil {
	//		panic(err)
	//	}



	if results.compiled == false {
		// compilation error occurred

		// TODO : parse out file path from compile errors
		fmt.Println(string(compileResults))

		// store results with the failed output
	}

	if results.compiled {
		// call execute
		fmt.Println("Compiled")
	}
}





from execute:
			/*execResults, err := execCmd.CombinedOutput()
			if string(execResults) == "" {
			}
			done <- err*/

/*

	fullCommand := programName + " < " + "test" + strconv.Itoa(i) + ".txt" + " 1> " + programOutputFile

	var stdError bytes.Buffer

	cmdToRun := programName
	args := []string{""}
	procAttr := new(os.ProcAttr)

	inFile, err := os.Open("test" + strconv.Itoa(i) + ".txt")
	outFile, err := os.Create(programOutputFile)
	if err != nil {
		panic("BAD")
	}
	procAttr.Files = []*os.File{inFile, outFile, os.Stderr}

	execResults, err := os.StartProcess(cmdToRun, args, procAttr)

	inFile.Close()
	outFile.Close()

*/
