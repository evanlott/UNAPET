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

type SubmissionInfo struct {
	courseName		string
	assignmentName		string
	studentUserID		int
	studentUsername		string
	grade			int
	comments		string
	compiled		int
	results			string
	submissionNum		int
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

//Hannah
func buildSubmissionStruct(assignmentName string, courseName string) (SubmissionInfo, error) {
	
	submission := SubmissionInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return submission, errors.New("No connection")
	}

	defer db.Close()
	
	rows, err := db.Query("select * from Submissions where CourseName = ? and AssignmentName = ?", courseName, assignmentName)

	if err != nil {
		return submission, errors.New("DB error")
	}

	if rows.Next() == false {
		return submission, errors.New("Invalid submission.")
	} else {
		rows.Scan(&submission.courseName, &submission.assignmentName, 
		&submission.studentUserID, &submission.grade, &submission.comments, 
		&submission.compiled, &submission.results, &submission.submissionNum)
	}
	
	rows, err := db.Query("select Username from Users where UserID = ?", submission.studentUserID)

	if err != nil {
		return submission, errors.New("DB error")
	}
	
	if rows.Next() == false {
		return submission, errors.New("Cannot get username.")
	} else {
		rows.Scan(&submission.studentUsername)
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

//Hannah
//returns the last submission for a particular student in a certain course for a certain assignment
func loadSubmission(student int, courseName string, assignmentName string) (SubmissionInfo, error){
	
	var submission SubmissionInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return submission, errors.New("No connection")
	}

	defer db.Close()

	//I know this seems unecessarily long, but it is what I have to do to survive 
	rows, err := db.Query("select * from Submissions where student=? and courseName=? and assignmentName=? and SubmissionNumber=(Select Max(SubmissionNumber)from Submissions where student=? and courseName=? and assignmentName=?)", student, courseName, assignmentName, student, courseName, assignmentName)

	if err != nil {
		return submission, errors.New("Query error.")
	}

	if rows.Next() == false {
		return submission, errors.New("Student is not enrolled in a course.")
	}

	//should just be able to call buildSubmissionStruct function, but not sure
	//how to do this with specific submission number 
	if rows.Next() == false {
		return submission, errors.New("Invalid submission.")
	} else {
		rows.Scan(&submission.courseName, &submission.assignmentName, 
		&submission.studentUserID, &submission.grade, &submission.comments, 
		&submission.compiled, &submission.results, &submission.submissionNum)
	}

	return submission, nil
}

//Hannah
func loadStudentsInCourse(courseName string) ([]UserInfo, error) {

	var users []UserInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return users, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("SELECT Username from Users where UserID IN (SELECT Student from StudentCourses where CourseName = ?)", courseName)

	if err != nil {
		return users, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var user string

		if rows.Next() == false {
			break
		}

		rows.Scan(&user)

		userStruct, err := buildUserStruct(user) 

		if err != nil {
			return users, err
		}

		users = append(users, userStruct)
	}

	return users, nil
}

//returns all users in the system for admin use 
//Hannah
func loadAllUsers() ([]UserInfo, error) {} {

	var users []UserInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return users, errors.New("No connection")
	}

<<<<<<< HEAD


 */
=======
	defer db.Close()

	rows, err := db.Query("SELECT Username from Users")

	if err != nil {
		return users, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var user string

		if rows.Next() == false {
			break
		}

		rows.Scan(&user)

		userStruct, err := buildUserStruct(user)

		if err != nil {
			return users, err
		}

		users = append(users, userStruct)
	}

	return users, nil
}

// returns some kind of grades for a user, does this need a struct?
// func loadGrades(student int, courseName string) (??, error) {}
>>>>>>> origin/master

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
