package pest

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

/*


 */

// not tested yet
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

/*


 */

// TODO : delete course's folder from disk
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

/*


 */

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

/*


 */

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
