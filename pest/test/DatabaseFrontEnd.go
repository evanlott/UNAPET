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
	si1Username        string
	si2UserID          int
	si2Username        string
	siGradeFlag        int
	siTestFlag         int
}

type UserInfo struct {
	userID			int
	firstName		string
	middleInitial		string
	lastName		string
	privLevel		int
	lastLogin		string
	pwdChangeFlag		string
	numLoginAttempts	int 
	enabled			int
}

type AssignmentInfo struct {
	assignmentName		int
	courseName		string
	assignmentDisplayName	string
	startDate		string
	endDate			string
	runtime			int
	compilerOptions		string
	numTestCases		int
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

	return course, nil

}

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

func buildAssignmentStruct(assignmentName string, courseName string) (AssignmentInfo, error) {
	assignment := AssignmentInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return user, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from Assignments where CourseName = ? and AssignmentName = ?", courseName, assignmentName)

	if err != nil {
		return user, errors.New("DB error")
	}

	if rows.Next() == false {
		return user, errors.New("Invalid Assignment.")
	} else {
		rows.Scan(&assignment.courseName, &assignment.assignmentDisplayName, &assignment.assignmentName,
			&assignment.startDate, &assignment.endDate, &assignment.runtime, &assignment.compilerOptions, 
			&assignment.numTestCases)
	}

	return assignment, nil
}

func main() {
	course, err := buildCourseStruct("JerkinsCS15502SP17")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%+v\n", course)

	return

}
