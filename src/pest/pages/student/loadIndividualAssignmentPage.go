// Nathan

package main

import (
	"html/template"
	"net/http"
	"net/http/cgi"
	"pest/functions"
	"pest/pages"

	_ "github.com/go-sql-driver/mysql"
)

const THIS_PAGE string = "Student/"

type PageContents struct {
	User       functions.UserInfo
	Assignment functions.AssignmentInfo
	Submission functions.SubmissionInfo
	SourceCode string
}

func main() {

	if err := cgi.Serve(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var Page PageContents

		// load this!
		Page.SourceCode = "Source!!!"

		username, err := functions.AuthUser(req)

		userInfo, err := functions.BuildUserStruct(username)
		functions.HandleError(err)

		Page.User = userInfo

		assignmentName := req.URL.Query().Get("assignmentName")
		courseName := req.URL.Query().Get("courseName")

		pages.EnsureNotNull(assignmentName)
		pages.EnsureNotNull(courseName)

		//assignmentName := "0"
		//courseName := "TerwilligerCS15501SP17"

		Page.Assignment, err = functions.BuildAssignmentStruct(assignmentName, courseName)
		functions.HandleError(err)

		Page.Submission, err = functions.LoadLastSubmission(userInfo.UserID, courseName, assignmentName)
		functions.HandleError(err)

		//tpl, err := template.ParseGlob(pages.WWW_ROOT + THIS_PAGE)
		//functions.HandleError(err)

		header := res.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		tpl, err := template.ParseFiles(pages.WWW_ROOT + THIS_PAGE + "individualAssignmentPage.gohtml")
		functions.HandleError(err)
		//template.ParseFiles()

		err = tpl.Execute(res, Page)
		functions.HandleError(err)

	})); err != nil {
		functions.ErrorResponse("CGI request failed.")
	}
}
