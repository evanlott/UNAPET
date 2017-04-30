// Test2
package main

/*
// TODO : change the paths in this file to reference a constant instead of hard coding
*/

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type CourseInfo struct {
	CourseName         string
	DisplayName        string
	StartDate          string
	EndDate            string
	CourseDescription  string
	InstructorUserID   int
	InstructorUsername string
	Si1UserID          int
	Si2UserID          int
	SiGradeFlag        int
	SiTestFlag         int
	Semester           string
	Year               string
}

type UserInfo struct {
	UserName         string
	UserID           int
	FirstName        string
	MiddleInitial    string
	LastName         string
	PrivLevel        int
	LastLogin        string
	PwdChangeFlag    string
	NumLoginAttempts int
	Enabled          int
}

type AssignmentInfo struct {
	AssignmentName        int
	CourseName            string
	AssignmentDisplayName string
	StartDate             string
	EndDate               string
	Runtime               int
	CompilerOptions       string
	NumTestCases          int
}

type SubmissionInfo struct {
	CourseName      string
	AssignmentName  string
	StudentUserID   int
	StudentUsername string
	Grade           int
	Comments        string
	Compiled        bool // changed this, will it cause problems? was int
	Results         string
	SubmissionNum   int
}

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"
const SHELL_NAME string = "ksh"

//Todd Gibson Test Driver
func main() {
	//main driver -- Todd Gibson Testing
	//isEnrolled(10113, "JerkinsCS15502SP17")
	//fmt.Println(isEnrolled(10113, "JerkinsCS15502SP17"))

	//Test BuildCourseStruct
	/*
			BuildCourseStruct("JerkinsCS15502SP17")
			BuildCourseStruct("RodenCS15501SP19")
			BuildCourseStruct("UNACourseSP17")

		BuildCourseStruct("JerkinsCS1550217") //should fail for wrong course name -- passed test case
	*/

	//Test BuildUserStruct --  Todd Gibson Testing
	/*
		BuildUserStruct("tggibson12")
		BuildUserStruct("tggibson1234")
	*/

	//Test BuildAssignmentStruct -- Todd Gibson Testing
	/*
		BuildAssignmentStruct("2", "JerkinsCS15502SP17")
		BuildAssignmentStruct("vd", "JerkinsCS15502SP17") -- should fail -- passed test case
	*/

	//Test buildSubmissionStruct
	//buildSubmissionStruct(AssignmentName string, CourseName string)
}

func buildSubmissionStruct(AssignmentName string, CourseName string) (SubmissionInfo, error) {

	submission := SubmissionInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return submission, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from Submissions where CourseName = ? and AssignmentName = ?", CourseName, AssignmentName)

	if err != nil {
		return submission, errors.New("DB error")
	}

	if rows.Next() == false {
		return submission, errors.New("Invalid submission.")
	} else {

		var grade sql.NullInt64
		var comments sql.NullString
		var compiled sql.NullBool
		var results sql.NullString

		rows.Scan(&submission.CourseName, &submission.AssignmentName,
			&submission.StudentUserID, &submission.Grade, &submission.Comments,
			&submission.Compiled, &submission.Results, &submission.SubmissionNum)

		submission.Grade = nullInt(grade)
		submission.Comments = nullString(comments)
		submission.Compiled = nullBool(compiled)
		submission.Results = nullString(results)

	}

	rows, err = db.Query("select Username from Users where UserID = ?", submission.StudentUserID)

	if err != nil {
		return submission, errors.New("DB error")
	}

	if rows.Next() == false {
		return submission, errors.New("Cannot get username.")
	} else {
		rows.Scan(&submission.StudentUsername)
	}

	return submission, nil
}

func nullInt(candidate sql.NullInt64) int {

	if candidate.Valid == false {
		return -1
	}

	return int(candidate.Int64)
}
func nullBool(candidate sql.NullBool) bool {

	if candidate.Valid == false {
		return false
	}

	return candidate.Bool
}

func nullString(candidate sql.NullString) string {

	if candidate.Valid == false {
		return candidate.String
	}

	return ""
}
func BuildUserStruct(username string) (UserInfo, error) {
	user := UserInfo{}

	user.UserName = username

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return user, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select UserID, FirstName, MiddleInitial, LastName, PrivLevel, LastLogin, PwdChangeFlag, NumLoginAttempts, Enabled from Users where Username = ?", username)

	if err != nil {
		fmt.Println("error on this line")
		return user, errors.New("DB error")
	}

	if rows.Next() == false {
		fmt.Println("error on this line -- bad user")
		return user, errors.New("Invalid User.")
	} else {
		rows.Scan(&user.UserID, &user.FirstName, &user.MiddleInitial,
			&user.LastName, &user.PrivLevel, &user.LastLogin, &user.PwdChangeFlag,
			&user.NumLoginAttempts, &user.Enabled)
	}

	//print user information -- Todd Gibson
	/*
		UserName         string
		UserID           int
		FirstName        string
		MiddleInitial    string
		LastName         string
		PrivLevel        int
		LastLogin        string
		PwdChangeFlag    string
		NumLoginAttempts int
		Enabled
	*/
	fmt.Println("username")
	fmt.Println(user.UserName)
	fmt.Println("UserID")
	fmt.Println(user.UserID)
	fmt.Println("FirstName")
	fmt.Println(user.FirstName)
	fmt.Println("MiddleInitial")
	fmt.Println(user.MiddleInitial)
	fmt.Println("LastName")
	fmt.Println(user.LastName)
	fmt.Println("PrivLevel")
	fmt.Println(user.PrivLevel)
	fmt.Println("LastLogin")
	fmt.Println(user.LastLogin)
	fmt.Println("PwdChangeFlag")
	fmt.Println(user.PwdChangeFlag)
	fmt.Println("NumLoginAttempts")
	fmt.Println(user.NumLoginAttempts)
	fmt.Println("Enabled")
	fmt.Println(user.Enabled)

	return user, nil

}

func BuildCourseStruct(CourseName string) (CourseInfo, error) {

	course := CourseInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return course, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from CourseDescription where CourseName = ?", CourseName)

	if err != nil {
		fmt.Println("Database error")
		return course, errors.New("DB error")
	}

	if rows.Next() == false {
		fmt.Println("invalid course name")
		return course, errors.New("Invalid Course.")
	} else {

		var si1 sql.NullInt64
		var si2 sql.NullInt64

		rows.Scan(&course.CourseName, &course.DisplayName, &course.CourseDescription,
			&course.InstructorUserID, &course.StartDate, &course.EndDate, &si1,
			&si2, &course.SiGradeFlag, &course.SiTestFlag)

		course.Si1UserID = nullInt(si1)
		course.Si2UserID = nullInt(si2)
	}

	rows, err = db.Query("select Username from Users where UserID = ?", course.InstructorUserID)

	if err != nil {
		fmt.Println("Database error 2")
		return course, errors.New("DB error")
	}

	if rows.Next() == false {
		fmt.Println("Invalid Instructor UserID")
		return course, errors.New("Invalid Instructor UserID")
	} else {
		rows.Scan(&course.InstructorUsername)
	}
	/*
		rows, err = db.Query("select Username from Users where UserID = ?", course.Si1UserID)
		if err != nil {
			return course, errors.New("DB error")
		}
		if rows.Next() == false {
			return course, errors.New("Invalid SI1 UserID")
		} else {
			rows.Scan(&course.si1Username)
		}
		rows, err = db.Query("select Username from Users where UserID = ?", course.Si2UserID)
		if err != nil {
			return course, errors.New("DB error")
		}
		if rows.Next() == false {
			return course, errors.New("Invalid SI2 UserID")
		} else {
			rows.Scan(&course.si2Username)
		}
	*/

	course.Semester = string(course.CourseName[(len(course.CourseName) - 4):(len(course.CourseName) - 2)])
	course.Year = "20" + string(course.CourseName[(len(course.CourseName)-2):len(course.CourseName)])

	if course.Semester == "FA" {
		course.Semester = "FALL"
	} else if course.Semester == "SP" {
		course.Semester = "SPRING"
	} else {
		course.Semester = "SUMMER"
	}

	//print course name to check database and see if results match
	fmt.Println(course)
	return course, nil

}
func BuildAssignmentStruct(AssignmentName string, CourseName string) (AssignmentInfo, error) {
	assignment := AssignmentInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return assignment, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from Assignments where CourseName = ? and AssignmentName = ?", CourseName, AssignmentName)

	if err != nil {
		return assignment, errors.New("DB error")
	}

	if rows.Next() == false {
		fmt.Println("invalid assignment")
		return assignment, errors.New("Invalid Assignment.")
	} else {
		rows.Scan(&assignment.CourseName, &assignment.AssignmentDisplayName, &assignment.AssignmentName,
			&assignment.StartDate, &assignment.EndDate, &assignment.Runtime, &assignment.CompilerOptions,
			&assignment.NumTestCases)
	}

	fmt.Println(assignment)
	return assignment, nil
}
