package pest

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"

func importCSV(name string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	mysql.RegisterLocalFile(name)

	// TODO : Solve password @dummy issue, also CSV quotation issue, trailing comma issue
	_, err = db.Exec("LOAD DATA LOCAL INFILE '" + name + "' INTO TABLE Users FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE 1 LINES (@dummy, FirstName, MiddleInitial, LastName, UserName, Password, @dummy, @dummy, @dummy, @dummy, @dummy)")

	if err != nil {
		return errors.New("Import failed.")
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

func gradeAssignment(userID int, courseName string, assignmentName string, submissionNum int, grade int) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("update Submissions set grade=? where Student=? and CourseName=? and AssignmentName=? and SubmissionNumber=?", grade, userID, courseName, assignmentName, submissionNum)

	if err != nil {
		return errors.New("Grade update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any submissions.")
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

/*


 */

func editUser(userID int, firstName string, MI string, lastName string, privLevel int) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("update Users set FirstName=?, MiddleInitial=?, LastName=?, PrivLevel=? where UserID=?", firstName, MI, lastName, privLevel, userID)

	if err != nil {
		return errors.New("User update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any users.")
	}

	return nil
}

/*
func createAssignment(courseName string, assignmentName string, runtime int, numTestCases int, compilerOptions string, startDate string, endDate string) {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		panic("No connection")
	}

	res, err := db.Exec("INSERT INTO `Assignments` (`CourseName`, `AssignmentName`, `StartDate`, `EndDate`, `MaxRuntime`, `CompilerOptions`, `NumTestCases`) VALUES (?, ?, '2017-04-11 00:00:00', '2017-04-29 00:00:00', '2000', '-Wall', '0');, startDate+" 23:59:59", ?, ?, ?")

	if err != nil {
		panic("Start/end update failed.")
	}

}
*/

func editStartEndAssignment(courseName string, assignmentName string, startDate string, endDate string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("update Assignments set StartDate=?, EndDate=? where CourseName=? and AssignmentName=?", startDate+" 23:59:59", endDate+" 23:59:59", courseName, assignmentName)

	if err != nil {
		return errors.New("Start/end update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any assignments.")
	}

	return nil
}

/*


 */

// add coursename to this, and plug in vars instead of hardcode
func editSubmissionComments(studentId int, assignmentName string, comments string) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}
	t := time.Now()
	currentTime := t.Format("2006-01-02 15:04:05")

	rows, err := db.Query("SELECT Submissions.comment FROM Submissions WHERE Submissions.student=10034 AND AssignmentName = \"Assignment 0\"")

	if err != nil {
		return errors.New("DB error")
	}

	var currentComments string

	if rows.Next() == false {
		return errors.New("Invalid comments.")
	} else {
		rows.Scan(&currentComments)
	}

	currentComments += currentTime + " - " + comments + "\n"

	editStatement, err := db.Exec("UPDATE Submissions SET Submissions.comment =? WHERE Submissions.student =? AND Submissions.AssignmentName =?", currentComments, studentId, assignmentName)

	if err != nil {
		return errors.New("Update failed.")
	} else {
		fmt.Println("Updated submission comments\n")
	}

	rowsAffected, err := editStatement.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any users.")
	}

	return nil
}

func deleteUser(userID int) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("delete from Users where UserID=? and not exists(select 1 from StudentCourses where Student=? limit 1)", userID, userID)

	if err != nil {
		return errors.New("User is currently enrolled in a class. Please remove the student from the class before deleting the user.")
	}
	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any users.")
	}

	return nil
}

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

func deleteAssignment(courseName string, assignmentName string) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("delete from Assignments where (AssignmentName =? and CourseName =?)", assignmentName, courseName)

	if err != nil {
		return errors.New("Assignment failed to delete.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any assignments or courses.")
	}

	return nil
}
