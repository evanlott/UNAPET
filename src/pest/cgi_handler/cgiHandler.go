package main

import (
	"net/http"
	"net/http/cgi"
	"pest/functions"
)

func main() {
	var req *http.Request

	req, err := cgi.Request()
	if err != nil {
		functions.ErrorResponse("Cannot process CGI request. Malformed HTTP POST or server error.")
		return
	}

	form := functions.ProcessForm(req)

	if form.ActionType == "" || form.Action == "" {
		functions.ErrorResponse("Cannot process CGI request. Malformed HTTP POST or server error.")
		return
	}

	var success bool
	var retString string

	// call the function with name == Action, pass it the form
	// returns bool, string
	switch form.ActionType {
	case "db":
		success, retString = functions.CallDbHelper(form.Action, form) // Invoke(dbHelper, form.Action, form)
	case "upload":
		success, retString = functions.CallUpload(form.Action, req)
	case "email":
		success, retString = functions.CallEmailFunction(form.Action, form)
	default:
		functions.ErrorResponse("Unknown request received.")
	}

	// redirect to the appropriate page if the Action succeeded
	// send an error if it did not
	if success {
		functions.RedirectTo(retString)
	} else {
		functions.ErrorResponse(retString)
	}

}
