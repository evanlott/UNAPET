// Nathan

package functions

import (
	"net/http"
)

func CallDbHelper(name string, form Request) (bool, string) {

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

	return false, "Requested Action is not implemented, or you have made an invalid request."

}

func CallUpload(name string, req *http.Request) (bool, string) {

	switch name {
	case "callCreateAssignment":
		return callCreateAssignment(req)
	case "sourceCodeUpload":
		return sourceCodeUpload(req)
	case "uploadCSV":
		return uploadCSV(req)
	}

	return false, "Requested Action is not implemented, or you have made an invalid request."

}

func CallEmailFunction(name string, form Request) (bool, string) {

	switch name {
	case "sendRandomPassword":
		return callSendRandomPassword(form)
	}

	return false, "Requested Action is not implemented, or you have made an invalid request."

}

func callSendRandomPassword(form Request) (bool, string) {

	userName := form.UserName

	err := sendRandomPassword(userName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callGradeSubmission(form Request) (bool, string) {

	err := gradeSubmission(form.UserID, form.CourseName, form.AssignmentName, form.SubNum, form.Grade)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callEvaluate(form Request) (bool, string) {

	/*
		// also check if this user is who they claim to be

		if (!(isLoggedIn(form.UserID)) || !(isEnrolled(form.UserID, form.CourseName))) {
			return false, "You are not logged in, or you do not have permission to submit to this class."
		}
	*/

	err := evaluate(form.CourseName, form.AssignmentName, form.UserName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callCreateCourse(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || (getPrivLevel(form.UserID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you do not have permission to create a course."
		}
	*/

	err := createCourse(form.CourseName, form.CourseDisplayName, form.CourseDescription, form.Instructor, form.StartDate, form.EndDate, form.Si1, form.Si2, form.SiGradeFlag, form.SiTestCaseFlag)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

/*
func callImportCSV(form Request) (bool, string) {


		if (!(isLoggedIn(form.UserID)) || (getPrivLevel(form.UserID) < PRIV_INSTRUCTOR)) {
			return false, "You are not logged in, or you do not have permission to upload a csv."
		}


	err := importCSV(form.FileName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}
*/

func callDeleteCourse(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || (getPrivLevel(form.UserID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you are not an admin user."
		}
	*/

	err := deleteCourse(form.CourseName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callEditCourseDescription(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || (getPrivLevel(form.UserID) < PRIV_INSTRUCTOR)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editCourseDescription(form.CourseName, form.CourseDescription)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callEditStartEndCourse(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || (getPrivLevel(form.UserID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editStartEndCourse(form.CourseName, form.StartDate, form.EndDate)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callEditUser(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || (getPrivLevel(form.UserID) < PRIV_INSTRUCTOR)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editUser(form.UserID, form.FirstName, form.MI, form.LastName, form.PrivLevel)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callEditStartEndAssignment(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || !(isInstructor(form.UserID, form.CourseName))) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := editStartEndAssignment(form.CourseName, form.AssignmentName, form.StartDate, form.EndDate)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callMakeSubmissionComment(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || !(isInstructor(form.UserID, form.CourseName))) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := makeSubmissionComment(form.UserID, form.AssignmentName, form.SubNum, form.Comments)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callDeleteUser(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) || (getPrivLevel(form.UserID) < PRIV_ADMIN)) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := deleteUser(form.UserID)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callDeleteAssignment(form Request) (bool, string) {

	/*
		if (!(isLoggedIn(form.UserID)) ||  !(isInstructor(form.UserID, form.CourseName))) {
			return false, "You are not logged in, or you do not have permission for this operation."
		}
	*/

	err := deleteAssignment(form.CourseName, form.AssignmentName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}

func callCreateUser(form Request) (bool, string) {

	err := createUser(form.FirstName, form.MI, form.LastName, form.UserName, form.PrivLevel, form.CourseName)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage

}

func callChangePassword(form Request) (bool, string) {

	err := changePassword(form.UserName, form.Password)

	if err != nil {
		return false, err.Error()
	}

	return true, form.FromPage
}
