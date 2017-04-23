package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

// TODO not tested yet
//---------------------------------------------------------------------------
//Inputs: course name, course display name, course description, instructor
//	start date for course, end date for course, supplemental 
//	instructor 1, supplemental instructor 2, grade flag for 
//	supplemental instructors, test flag for supplemental instructors
//Outputs: returns errors if the course could not be created 
//Written By: Evan Lott 
//Purpose: This function will be used by the administrator to create a 
//	course. It will add the course to the CourseDescription table in
//	the database.  
//---------------------------------------------------------------------------
func createCourse(courseName string, courseDisplayName string, courseDescription string, instructor int, startDate string, endDate string, si1 int, si2 int, siGradeFlag bool, siTestCaseFlag bool) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO CourseDescription VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", courseName, courseDisplayName, courseDescription, instructor, startDate, endDate, si1, si2, siGradeFlag, siTestCaseFlag)

	if err != nil {
		return errors.New("Error inserting course.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any assignments.")
	}

	return nil
}

// TODO : delete course's folder from disk
//---------------------------------------------------------------------------
//Inputs: course name
//Outputs: returns errors if the course fails to delete 
//Written By: Hannah Hopkins
//Purpose: This function will be used by the administrator to delete a 
//	course. It will remove the course from the CourseDescription
//	table in the database.  
//---------------------------------------------------------------------------
func deleteCourse(courseName string) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("delete from CourseDescription where CourseName =?", courseName)

	if err != nil {
		return errors.New("Course failed to delete.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any courses.")
	}

	return nil
}

//---------------------------------------------------------------------------
//Inputs: course name, course description 
//Outputs: returns errors if the course description could not be updated
//Written By: Eileen Drass and Evan Lott 
//Purpose: This function will be used by the instructor to edit the course
//	description for a course. It will update the course in the 
//	in the CourseDescription table in the database.  
//---------------------------------------------------------------------------
func editCourseDescription(courseName string, courseDescription string) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	_, err = db.Exec("UPDATE CourseDescription SET CourseDescription.CourseDescription=? WHERE CourseDescription.CourseName=? ", courseDescription, courseName)

	if err != nil {
		return errors.New("Update failed.")
	}

	return nil

}

//---------------------------------------------------------------------------
//Inputs: course name, start date for course, end date for course
//Outputs: returns errors if the start and end date for the course could not
//	be updated
//Written By: Hannah Hopkins and Nathan Huckaba 
//Purpose: This function will be used by the administrator to edit the start 
//	and end dates for a course. It will update the course in the 
//	CourseDescription table in the database.  
//---------------------------------------------------------------------------
func editStartEndCourse(courseName string, startDate string, endDate string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("update CourseDescription set StartDate=?, EndDate=? where CourseName=?", startDate+" 23:59:59", endDate+" 23:59:59", courseName)

	if err != nil {
		return errors.New("Start/end update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any courses.")
	}

	return nil

}
