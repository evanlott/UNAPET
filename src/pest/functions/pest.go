// Nathan

package functions

import (
	"fmt"
	"net/http"
	"os"
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
	ActionType            string
	Action                string
	CourseName            string
	CourseDisplayName     string
	AssignmentName        string
	AssignmentDisplayName string
	SubNum                int
	Grade                 int
	UserID                int
	UserName              string
	CourseDescription     string
	StartDate             string
	EndDate               string
	FirstName             string
	MI                    string
	LastName              string
	PrivLevel             int
	Comments              string
	FileName              string
	CompilerOptions       string
	NumTestCases          int
	Runtime               int
	Instructor            int
	Si1                   int
	Si2                   int
	SiGradeFlag           bool
	SiTestCaseFlag        bool
	FromPage              string
	Password              string
}

//---------------------------------------------------------------------------
//Inputs: message
//Outputs: This function does not return anything. 
//Written By: Nathan Huckaba
//Purpose: This function prints out an error response. 
//---------------------------------------------------------------------------
func ErrorResponse(msg string) {
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
func HandleError(err error) {

	if err != nil {
		ErrorResponse(err.Error())
		os.Exit(0)
	}
}

/*
func HandleErrorHTTP(err error, res http.ResponseWriter) {

	if err != nil {
		http.Error(res, "Internal server error.", http.StatusInternalServerError)
		os.Exit(0)
	}
}
*/

//---------------------------------------------------------------------------
//Inputs: page
//Outputs: This function does not return anything. 
//Written By: Nathan Huckaba
//Purpose: This function prints out a oage to redirect to. 
//---------------------------------------------------------------------------
func RedirectTo(page string) {
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
func ProcessForm(req *http.Request) Request {

	var form Request

	form.ActionType = req.FormValue("actionType")
	form.Action = req.FormValue("action")
	form.UserID, _ = strconv.Atoi(req.FormValue("userID"))
	form.CourseName = req.FormValue("courseName")
	form.CourseDisplayName = req.FormValue("courseDisplayName")
	form.AssignmentName = req.FormValue("assignmentName")
	form.SubNum, _ = strconv.Atoi(req.FormValue("submissionNum"))
	form.Grade, _ = strconv.Atoi(req.FormValue("grade"))
	form.UserName = req.FormValue("userName")
	form.FileName = req.FormValue("fileName")
	form.CourseDescription = req.FormValue("courseDescription")
	form.StartDate = req.FormValue("startDate")
	form.EndDate = req.FormValue("endDate")
	form.FirstName = req.FormValue("firstName")
	form.MI = req.FormValue("MI")
	form.LastName = req.FormValue("lastName")
	form.PrivLevel, _ = strconv.Atoi(req.FormValue("privLevel"))
	form.Comments = req.FormValue("comments")
	form.FileName = req.FormValue("fileName")
	form.AssignmentDisplayName = req.FormValue("assignmentDisplayName")
	form.CompilerOptions = req.FormValue("compilerOptions")
	form.NumTestCases, _ = strconv.Atoi(req.FormValue("numTestCases"))
	form.Runtime, _ = strconv.Atoi(req.FormValue("runtime"))
	form.Instructor, _ = strconv.Atoi(req.FormValue("instructor"))
	form.Si1, _ = strconv.Atoi(req.FormValue("si1"))
	form.Si2, _ = strconv.Atoi(req.FormValue("si2"))
	form.FromPage = req.Referer()
	form.Password = req.FormValue("password")

	temp, _ := strconv.Atoi(req.FormValue("siGradeFlag"))
	if temp == 1 {
		form.SiGradeFlag = true
	} else {
		form.SiGradeFlag = false
	}

	temp, _ = strconv.Atoi(req.FormValue("siTestCaseFlag"))
	if temp == 1 {
		form.SiTestCaseFlag = true
	} else {
		form.SiTestCaseFlag = false
	}

	return form
}
