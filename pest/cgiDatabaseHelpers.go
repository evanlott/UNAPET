package pest

import "reflect"

// calls a function of an interface by name and passes it args
func Invoke(function interface{}, name string, form Request) (bool, string) {

	input := make([]reflect.Value, 1)

	input[0] = reflect.ValueOf(form)

	method := reflect.ValueOf(function).MethodByName(name)

	// function exists, call it
	if method.IsValid() {
		result := method.Call(input)

		// function existed return what it returns
		return result[0].Interface().(bool), result[1].Interface().(string)
	}

	// function doesn't exist, return false and err msg
	return false, "Requested action is not implemented, or you have made an invalid request."
}

type dbHelpers struct{}

// Will have to check user access priv. before executing these!
// So need a checkPriv function that goes to the database and returns their priv level
// Will also have to check that they are logged in as the user they claim to be in the request!

func (_ dbHelpers) callEvaluate(form Request) (bool, string) {

	err := evaluate(form.courseName, form.assignmentName, form.userName)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callGradeAssignment(form Request) (bool, string) {

	err := gradeAssignment(form.userID, form.courseName, form.assignmentName, form.subNum, form.grade)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callImportCSV(form Request) (bool, string) {

	err := importCSV(form.fileName)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callDeleteCourse(form Request) (bool, string) {
	err := deleteCourse(form.courseName)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callEditCourseDescription(form Request) (bool, string) {
	err := editCourseDescription(form.courseName, form.courseDescription)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callEditStartEndCourse(form Request) (bool, string) {

	err := editStartEndCourse(form.courseName, form.startDate, form.endDate)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callEditUser(form Request) (bool, string) {

	err := editUser(form.userID, form.firstName, form.MI, form.lastName, form.privLevel)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callEditStartEndAssignment(form Request) (bool, string) {

	err := editStartEndAssignment(form.courseName, form.assignmentName, form.startDate, form.endDate)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callEditSubmissionComments(form Request) (bool, string) {

	err := editSubmissionComments(form.userID, form.assignmentName, form.comments)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callDeleteUser(form Request) (bool, string) {

	err := deleteUser(form.userID)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}

func (_ dbHelpers) callDeleteAssignment(form Request) (bool, string) {

	err := deleteAssignment(form.courseName, form.assignmentName)

	if err != nil {
		return false, err.Error()
	}

	return true, "?.html"
}
