package main

import (
	"net/http"
)

//---------------------------------------------------------------------------
//Inputs: The inputs are the name of the case and the request form.
//Outputs: This function either returns a bool along with a string or the
// 	CGI helper function to be called.
//Written By: Nathan Huckaba
//Purpose: This function contains a switch statement which will determine
//	which CGI helper function will be called.
//---------------------------------------------------------------------------
func callDbHelper(name string, form Request) (bool, string) {

	switch name {
	case "callGradeSubmission":
		return callGradeSubmission(form)
	case "callEvaluate":
		return callEvaluate(form)
	case "callDeleteCourse":
		return callDeleteCourse(form)
	case "callEditCourseDescription":
		return callEditCourseDescription(form)
	case "callEditStartEndCourse":
		return callEditStartEndCourse(form)
	case "callEditUser":
		return callEditUser(form)
	case "callEditStartEndAssignment":
		return callEditStartEndAssignment(form)
	case "callMakeSubmissionComment":
		return callMakeSubmissionComment(form)
	case "callDeleteUser":
		return callDeleteUser(form)
	case "callDeleteAssignment":
		return callDeleteAssignment(form)
	case "callCreateCourse":
		return callCreateCourse(form)
	case "callCreateUser":
		return callCreateUser(form)
	case "callChangePassword":
		return callChangePassword(form)
	}

	return false, "Requested action is not implemented, or you have made an invalid request."

}

//---------------------------------------------------------------------------
//Inputs: The inputs are the name of the case and the http request.
//Outputs:  This function either returns a bool along with a string or the
// 	function to be called.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to handle functions which
//	deal with creating an assignment and uploading source code and CSV files.
//---------------------------------------------------------------------------
func callUpload(name string, req *http.Request) (bool, string) {

	switch name {
	case "callCreateAssignment":
		return callCreateAssignment(req)
	case "sourceCodeUpload":
		return sourceCodeUpload(req)
	case "uploadCSV":
		return uploadCSV(req)
	}

	return false, "Requested action is not implemented, or you have made an invalid request."

}

//---------------------------------------------------------------------------
//Inputs: The inputs are the name of the case and the request form.
//Outputs:  This function either returns a bool along with a string or the
// 	CGI helper function callSendRandomPassword.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the
//	callSendRandomPassword function.
//---------------------------------------------------------------------------
func callEmailFunction(name string, form Request) (bool, string) {

	switch name {
	case "sendRandomPassword":
		return callSendRandomPassword(form)
	}

	return false, "Requested action is not implemented, or you have made an invalid request."

}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call sendRandomPassword
//	function.
//---------------------------------------------------------------------------
func callSendRandomPassword(form Request) (bool, string) {

	userName := form.userName

	err := sendRandomPassword(userName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the gradeSubmission
//	function.
//---------------------------------------------------------------------------
func callGradeSubmission(form Request) (bool, string) {

	err := gradeSubmission(form.userID, form.courseName, form.assignmentName, form.subNum, form.grade)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the evaluate
//	function.
//---------------------------------------------------------------------------
func callEvaluate(form Request) (bool, string) {

	/*
		// also check if this user is who they claim to be
		if (!(isLoggedIn(form.userID)) || !(isEnrolled(form.userID, form.courseName))) {
			return false, "You are not logged in, or you do not have permission to submit to this class."
		}
	*/

	err := evaluate(form.courseName, form.assignmentName, form.userName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the createCourse
//	function.
//---------------------------------------------------------------------------
func callCreateCourse(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || (getPrivLevel(form.userID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you do not have permission to create a course."
		}
	*/

	err := createCourse(form.courseName, form.courseDisplayName, form.courseDescription, form.instructor, form.startDate, form.endDate, form.si1, form.si2, form.siGradeFlag, form.siTestCaseFlag)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the deleteCourse
//	function.
//---------------------------------------------------------------------------
func callDeleteCourse(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || (getPrivLevel(form.userID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you are not an admin user."
		}
	*/

	err := deleteCourse(form.courseName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the
//	editCourseDescription function.
//---------------------------------------------------------------------------
func callEditCourseDescription(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || (getPrivLevel(form.userID) < PRIV_INSTRUCTOR)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editCourseDescription(form.courseName, form.courseDescription)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the
//	editStartEndCourse function.
//---------------------------------------------------------------------------
func callEditStartEndCourse(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || (getPrivLevel(form.userID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editStartEndCourse(form.courseName, form.startDate, form.endDate)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the editUser
//	function.
//---------------------------------------------------------------------------
func callEditUser(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || (getPrivLevel(form.userID) < PRIV_INSTRUCTOR)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editUser(form.userID, form.firstName, form.MI, form.lastName, form.privLevel)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the
//	editStartEndAssignment function.
//---------------------------------------------------------------------------
func callEditStartEndAssignment(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || !(isInstructor(form.userID, form.courseName))) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editStartEndAssignment(form.courseName, form.assignmentName, form.startDate, form.endDate)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the
//	makeSubmissionComment function.
//---------------------------------------------------------------------------
func callMakeSubmissionComment(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || !(isInstructor(form.userID, form.courseName))) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := makeSubmissionComment(form.userID, form.assignmentName, form.subNum, form.comments)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the deleteUser
//	function.
//---------------------------------------------------------------------------
func callDeleteUser(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) || (getPrivLevel(form.userID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := deleteUser(form.userID)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the deleteAssignment
//	function.
//---------------------------------------------------------------------------
func callDeleteAssignment(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.userID)) ||  !(isInstructor(form.userID, form.courseName))) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := deleteAssignment(form.courseName, form.assignmentName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the createUser
//	function.
//---------------------------------------------------------------------------
func callCreateUser(form Request) (bool, string) {

	err := createUser(form.firstName, form.MI, form.lastName, form.userName, form.password, form.privLevel, form.courseName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage

}

//---------------------------------------------------------------------------
//Inputs: This function's input includes the request form.
//Outputs: This function returns true if the request does not fail along with
//	the form's URL. It returns false and an error if the request fails.
//Written By: Nathan Huckaba
//Purpose: This function will be used in order to call the changePassword
//	function.
//---------------------------------------------------------------------------
func callChangePassword(form Request) (bool, string) {

	err := changePassword(form.userName, form.password)

	if err != nil {
		return false, err.Error()
	}

	return true, form.fromPage
}
