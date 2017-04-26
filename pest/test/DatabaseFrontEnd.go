// DatabaseFrontEnd
package main

import (
	"database/sql"
	"errors"
	"fmt"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)

type CourseInfo struct {
	courseName         string
	displayName        string
	startDate          string
	endDate            string
	courseDescription  string
	instructorUserID   int
	instructorUsername string
	si1UserID          int
	si2UserID          int
	siGradeFlag        int
	siTestFlag         int
	semester           string
	year               string
}

type UserInfo struct {
	userID           int
	firstName        string
	middleInitial    string
	lastName         string
	privLevel        int
	lastLogin        string
	pwdChangeFlag    string
	numLoginAttempts int
	enabled          int
}

type AssignmentInfo struct {
	assignmentName        int
	courseName            string
	assignmentDisplayName string
	startDate             string
	endDate               string
	runtime               int
	compilerOptions       string
	numTestCases          int
}

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

// Hannah
func buildCourseStruct(courseName string) (CourseInfo, error) {

	course := CourseInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return course, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from CourseDescription where CourseName = ?", courseName)

	if err != nil {
		return course, errors.New("DB error")
	}

	if rows.Next() == false {
		return course, errors.New("Invalid Course.")
	} else {
		rows.Scan(&course.courseName, &course.displayName, &course.courseDescription,
			&course.instructorUserID, &course.startDate, &course.endDate, &course.si1UserID,
			&course.si2UserID, &course.siGradeFlag, &course.siTestFlag)
	}

	rows, err = db.Query("select Username from Users where UserID = ?", course.instructorUserID)

	if err != nil {
		return course, errors.New("DB error")
	}

	if rows.Next() == false {
		return course, errors.New("Invalid Instructor UserID")
	} else {
		rows.Scan(&course.instructorUsername)
	}
	/*
		rows, err = db.Query("select Username from Users where UserID = ?", course.si1UserID)

		if err != nil {
			return course, errors.New("DB error")
		}

		if rows.Next() == false {
			return course, errors.New("Invalid SI1 UserID")
		} else {
			rows.Scan(&course.si1Username)
		}

		rows, err = db.Query("select Username from Users where UserID = ?", course.si2UserID)

		if err != nil {
			return course, errors.New("DB error")
		}

		if rows.Next() == false {
			return course, errors.New("Invalid SI2 UserID")
		} else {
			rows.Scan(&course.si2Username)
		}
	*/

	course.semester = string(course.courseName[(len(course.courseName) - 4):(len(course.courseName) - 2)])
	course.year = "20" + string(course.courseName[(len(course.courseName)-2):len(course.courseName)])

	if course.semester == "FA" {
		course.semester = "FALL"
	} else if course.semester == "SP" {
		course.semester = "SPRING"
	} else {
		course.semester = "SUMMER"
	}

	return course, nil

}

// Hannah
func buildUserStruct(username string) (UserInfo, error) {
	user := UserInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return user, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select UserID, FirstName, MiddleInitial, LastName, PrivLevel, LastLogin, PwdChangeFlag, NumLoginAttempts, Enabled from Users where Username = ?", username)

	if err != nil {
		return user, errors.New("DB error")
	}

	if rows.Next() == false {
		return user, errors.New("Invalid User.")
	} else {
		rows.Scan(&user.userID, &user.firstName, &user.middleInitial,
			&user.lastName, &user.privLevel, &user.lastLogin, &user.pwdChangeFlag,
			&user.numLoginAttempts, &user.enabled)
	}

	return user, nil

}

// Hannah
func buildAssignmentStruct(assignmentName string, courseName string) (AssignmentInfo, error) {
	assignment := AssignmentInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return assignment, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from Assignments where CourseName = ? and AssignmentName = ?", courseName, assignmentName)

	if err != nil {
		return assignment, errors.New("DB error")
	}

	if rows.Next() == false {
		return assignment, errors.New("Invalid Assignment.")
	} else {
		rows.Scan(&assignment.courseName, &assignment.assignmentDisplayName, &assignment.assignmentName,
			&assignment.startDate, &assignment.endDate, &assignment.runtime, &assignment.compilerOptions,
			&assignment.numTestCases)
	}

	return assignment, nil
}

// Nathan and Hannah
// Returns a slice of course structs representing all the courses in the DB, active or not
func loadAdminCards() ([]CourseInfo, error) {

	var courses []CourseInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return courses, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select CourseName from CourseDescription")

	if err != nil {
		return courses, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var courseName string

		if rows.Next() == false {
			break
		}

		rows.Scan(&courseName)

		courseStruct, err := buildCourseStruct(courseName)

		if err != nil {
			return courses, err
		}

		courses = append(courses, courseStruct)
	}

	return courses, nil
}

// Nathan
// Returns a slice of course structs describing all the courses an instructor has, active or not
func loadInstructorCards(userID int) ([]CourseInfo, error) {

	var courses []CourseInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return courses, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select CourseName from CourseDescription where Instructor=?", userID)

	if err != nil {
		return courses, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var courseName string

		if rows.Next() == false {
			break
		}

		rows.Scan(&courseName)

		courseStruct, err := buildCourseStruct(courseName)

		if err != nil {
			return courses, err
		}

		courses = append(courses, courseStruct)
	}

	return courses, nil
}

// Nathan
// Returns a course struct describing the course a student is enrolled in, active or not
func loadStudentCourse(userID int) (CourseInfo, error) {

	var course CourseInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return course, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select CourseName from StudentCourses where Student=?", userID)

	if err != nil {
		return course, errors.New("Query error.")
	}

	if rows.Next() == false {
		return course, errors.New("Student is not enrolled in a course.")
	}

	var courseName string

	rows.Scan(&courseName)

	course, err = buildCourseStruct(courseName)

	if err != nil {
		return course, err
	}

	return course, nil
}

// Nathan
// Returns all the assignments for a course as a slice
func loadAssignments(courseName string) ([]AssignmentInfo, error) {

	var assignments []AssignmentInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return assignments, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select AssignmentName from Assignments where CourseName=?", courseName)

	if err != nil {
		return assignments, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var assignmentName string

		if rows.Next() == false {
			break
		}

		rows.Scan(&assignmentName)

		assignmentStruct, err := buildAssignmentStruct(assignmentName, courseName)

		if err != nil {
			return assignments, err
		}

		assignments = append(assignments, assignmentStruct)
	}

	return assignments, nil

}

/*



 */

type SubmissionInfo struct {

	// Some stuff

}

// loads -last- submission info for a student in a course
// func loadSubmissions(student int, courseName string) (SubmissionInfo, error) {}

// returns a slice of UserInfo structs representing all students in a course
// func loadStudentsInCourse(courseName string) ([]UserInfo, error) {}

// returns all users in the system for admin use
// func loadAllUsers() ([]UserInfo, error) {}

// returns some kind of grades for a user, does this need a struct?
// func loadGrades(student int, courseName string) (??, error) {}

/*





 */

func main() {

	/*
		course, err := buildAssignmentStruct("0", "TerwilligerCS15501SP17")

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("%+v\n", course)
	*/

	x, err := loadAdminCards()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%v+\n", x)

	return

}
