package pest

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
	action            string
	courseName        string
	assignmentName    string
	subNum            int
	grade             int
	userID            int
	userName          string
	courseDescription string
	startDate         string
	endDate           string
	firstName         string
	MI                string
	lastName          string
	privLevel         int
	comments          string
	fileName          string
}

func errorResponse(msg string) {
	fmt.Printf("Status: 500 Bad\r\n")
	fmt.Printf("Content-Type: text/plain\r\n")
	fmt.Printf("\r\n")
	fmt.Printf("%s\r\n", msg)
}

func redirectTo(page string) {
	fmt.Printf("HTTP/1.1 303 See other\r\n")
	fmt.Printf("Location: /%s \r\n", page)
	fmt.Printf("\r\n")
}

// returns a Request struct with all info from form put into variables
func processForm(req *http.Request) Request {
	var form Request

	form.action = req.FormValue("action")
	form.userID, _ = strconv.Atoi(req.FormValue("userID"))
	form.courseName = req.FormValue("courseName")
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

	return form
}

func main() {
	var req *http.Request

	req, err := cgi.Request()
	if err != nil {
		errorResponse("Cannot process CGI request. Malformed HTTP POST or server error.")
		return
	}

	form := processForm(req)

	if form.action == "" {
		errorResponse("Cannot process CGI request. Malformed HTTP POST or server error.")
		return
	}

	// call the function with name == action, pass it the form
	// returns bool, string
	success, page := Invoke(dbHelpers{}, form.action, form)

	// redirect to the appropriate page if the action succeeded
	// send an error if it did not
	if success {
		redirectTo(page)
	} else {
		errorResponse(page)
	}

}
