package pest

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

/*


 */

func gradeSubmission(userID int, courseName string, assignmentName string, submissionNum int, grade int) error {

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
	}

	rowsAffected, err := editStatement.RowsAffected()

	if rowsAffected != 1 {
		return errors.New("Query didn't match any users.")
	}

	return nil
}
