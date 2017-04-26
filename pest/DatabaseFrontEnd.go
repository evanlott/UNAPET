// DatabaseFrontEnd
package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type CourseInfo struct {
	courseName         string
	displayName        string
	startDate          time.Time
	endDate            time.Time
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

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"
const SHELL_NAME string = "ksh"

func main() {
	fmt.Println("Hello World!")
}

func buildCourseStruct(courseName string) {
	course := CourseInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from CourseDescription where CourseName = ?", courseName)

	if err != nil {
		return ("DB error")
	}

	if rows.Next() == false {
		return ("Invalid Course.")
	} else {
		rows.Scan(&course.courseName, &course.displayName, &course.courseDescription,
			&course.instructorUserID, &course.startDate, &course.endDate, &course.si1UserID,
			&course.si2UserID, &course.siGradeFlag, &course.siTestFlag)
	}

	rows, err := db.Query("select Username from Users where UserID = ?", course.instructorUserID)

	if err != nil {
		return ("DB error")
	}

	if rows.Next() == false {
		return ("Invalid Instructor UserID")
	} else {
		rows.Scan(&course.instructorUsername)
	}

	rows, err := db.Query("select Username from Users where UserID = ?", course.si1UserID)

	if err != nil {
		return ("DB error")
	}

	if rows.Next() == false {
		return ("Invalid SI1 UserID")
	} else {
		rows.Scan(&course.si1Username)
	}

	rows, err := db.Query("select Username from Users where UserID = ?", course.si2UserID)

	if err != nil {
		return ("DB error")
	}

	if rows.Next() == false {
		return ("Invalid SI2 UserID")
	} else {
		rows.Scan(&course.si2Username)
	}

}
