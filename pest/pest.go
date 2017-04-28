// Nathan

package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"strconv"
)

/*
type LoggedInUser struct {
	userName
	userID
	sessionKey
}
*/

type Request struct {
	actionType            string
	action                string
	courseName            string
	courseDisplayName     string
	assignmentName        string
	assignmentDisplayName string
	subNum                int
	grade                 int
	userID                int
	userName              string
	courseDescription     string
	startDate             string
	endDate               string
	firstName             string
	MI                    string
	lastName              string
	privLevel             int
	comments              string
	fileName              string
	compilerOptions       string
	numTestCases          int
	runtime               int
	instructor            int
	si1                   int
	si2                   int
	siGradeFlag           bool
	siTestCaseFlag        bool
	fromPage              string
	password              string
}

//---------------------------------------------------------------------------
//Inputs: message
//Outputs: This function does not return anything. 
//Written By: Nathan Huckaba
//Purpose: This function prints out an error response. 
//---------------------------------------------------------------------------
func errorResponse(msg string) {
	fmt.Printf("Status: 500 Bad\r\n")
	fmt.Printf("Content-Type: text/plain\r\n")
	fmt.Printf("\r\n")
	fmt.Printf("%s\r\n", msg)
	fmt.Printf("If problem persists, please contact system admin.\r\n")
}

//---------------------------------------------------------------------------
//Inputs: page
//Outputs: This function does not return anything. 
//Written By: Nathan Huckaba
//Purpose: This function prints out a oage to redirect to. 
//---------------------------------------------------------------------------
func redirectTo(page string) {
	fmt.Printf("HTTP/1.1 303 See other\r\n")
	fmt.Printf("Location: %s \r\n", page)
	fmt.Printf("\r\n")
}

//---------------------------------------------------------------------------
//Inputs: http request
//Outputs: This function returns a Request struct with all info from form 
//	put into variables 
//Written By: Nathan Huckaba
//Purpose: This function is used to process a form. 
//---------------------------------------------------------------------------
func processForm(req *http.Request) Request {

	var form Request

	form.actionType = req.FormValue("actionType")
	form.action = req.FormValue("action")
	form.userID, _ = strconv.Atoi(req.FormValue("userID"))
	form.courseName = req.FormValue("courseName")
	form.courseDisplayName = req.FormValue("courseDisplayName")
	form.assignmentName = req.FormValue("assignmentName")
	form.subNum, _ = strconv.Atoi(req.FormValue("submissionNum"))
	form.grade, _ = strconv.Atoi(req.FormValue("grade"))
	form.userName = req.FormValue("userName")
	form.fileName = req.FormValue("fileName")
	form.courseDescription = req.FormValue("courseDescription")
	form.startDate = req.FormValue("startDate")
	form.endDate = req.FormValue("endDate")
	form.firstName = req.FormValue("firstName")
	form.MI = req.FormValue("MI")
	form.lastName = req.FormValue("lastName")
	form.privLevel, _ = strconv.Atoi(req.FormValue("privLevel"))
	form.comments = req.FormValue("comments")
	form.fileName = req.FormValue("fileName")
	form.assignmentDisplayName = req.FormValue("assignmentDisplayName")
	form.compilerOptions = req.FormValue("compilerOptions")
	form.numTestCases, _ = strconv.Atoi(req.FormValue("numTestCases"))
	form.runtime, _ = strconv.Atoi(req.FormValue("runtime"))
	form.instructor, _ = strconv.Atoi(req.FormValue("instructor"))
	form.si1, _ = strconv.Atoi(req.FormValue("si1"))
	form.si2, _ = strconv.Atoi(req.FormValue("si2"))
	form.fromPage = req.Referer()
	form.password = req.FormValue("password")

	temp, _ := strconv.Atoi(req.FormValue("siGradeFlag"))
	if temp == 1 {
		form.siGradeFlag = true
	} else {
		form.siGradeFlag = false
	}

	temp, _ = strconv.Atoi(req.FormValue("siTestCaseFlag"))
	if temp == 1 {
		form.siTestCaseFlag = true
	} else {
		form.siTestCaseFlag = false
	}

	return form
}

//---------------------------------------------------------------------------
func main() {
	var req *http.Request

	req, err := cgi.Request()
	if err != nil {
		errorResponse("Cannot process CGI request. Malformed HTTP POST or server error.")
		return
	}

	form := processForm(req)

	if form.actionType == "" || form.action == "" {
		errorResponse("Cannot process CGI request. Malformed HTTP POST or server error.")
		return
	}

	var success bool
	var retString string

	// call the function with name == action, pass it the form
	// returns bool, string
	switch form.actionType {
	case "db":
		success, retString = callDbHelper(form.action, form) // Invoke(dbHelper, form.action, form)
	case "upload":
		success, retString = callUpload(form.action, req)
	case "email":
		success, retString = callEmailFunction(form.action, form)
	default:
		errorResponse("Unknown request received.")
	}

	// redirect to the appropriate page if the action succeeded
	// send an error if it did not
	if success {
		redirectTo(retString)
	} else {
		errorResponse(retString)
	}

}
