package main

import (
	"html/template"
	"net/http"
	"net/http/cgi"
	"pest/functions"
	"pest/pages"

	_ "github.com/go-sql-driver/mysql"
)

const THIS_PAGE string = "Admin/"

type PageContents struct {
	User    functions.UserInfo
	Courses []functions.CourseInfo
}

func main() {

	if err := cgi.Serve(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		username, err := functions.AuthUser(req)
		functions.HandleError(err)

		var Page PageContents

		Page.User, err = functions.BuildUserStruct(username)
		functions.HandleError(err)

		Page.Courses, err = functions.LoadAdminCards()
		functions.HandleError(err)

		tpl, err := template.ParseGlob(pages.WWW_ROOT + THIS_PAGE + "*.gohtml")
		functions.HandleError(err)

		header := res.Header()
		header.Set("Content-Type", "text/html; charset=utf-8")

		err = tpl.ExecuteTemplate(res, "createUser.gohtml", Page)
		functions.HandleError(err)

	})); err != nil {
		functions.ErrorResponse("CGI request failed.")
	}
}
