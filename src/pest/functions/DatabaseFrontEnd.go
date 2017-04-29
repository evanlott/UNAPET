// DatabaseFrontEnd
package functions

import (
	"database/sql"
	"errors"

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

// Hannah
func BuildCourseStruct(CourseName string) (CourseInfo, error) {

	course := CourseInfo{}

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return course, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select * from CourseDescription where CourseName = ?", CourseName)

	if err != nil {
		return course, errors.New("DB error")
	}

	if rows.Next() == false {
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
		return course, errors.New("DB error")
	}

	if rows.Next() == false {
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

	return course, nil

}

// Hannah
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
		return user, errors.New("DB error")
	}

	if rows.Next() == false {
		return user, errors.New("Invalid User.")
	} else {
		rows.Scan(&user.UserID, &user.FirstName, &user.MiddleInitial,
			&user.LastName, &user.PrivLevel, &user.LastLogin, &user.PwdChangeFlag,
			&user.NumLoginAttempts, &user.Enabled)
	}

	return user, nil

}

// Hannah
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
		return assignment, errors.New("Invalid Assignment.")
	} else {
		rows.Scan(&assignment.CourseName, &assignment.AssignmentDisplayName, &assignment.AssignmentName,
			&assignment.StartDate, &assignment.EndDate, &assignment.Runtime, &assignment.CompilerOptions,
			&assignment.NumTestCases)
	}

	return assignment, nil
}

//Hannah
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

// Nathan and Hannah
// Returns a slice of course structs representing all the courses in the DB, active or not
func LoadAdminCards() ([]CourseInfo, error) {

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
		var CourseName string

		if rows.Next() == false {
			break
		}

		rows.Scan(&CourseName)

		courseStruct, err := BuildCourseStruct(CourseName)

		if err != nil {
			return courses, err
		}

		courses = append(courses, courseStruct)
	}

	return courses, nil
}

// Nathan
// Returns a slice of course structs describing all the courses an instructor has, active or not
func LoadInstructorCards(UserID int) ([]CourseInfo, error) {

	var courses []CourseInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return courses, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select CourseName from CourseDescription where Instructor=?", UserID)

	if err != nil {
		return courses, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var CourseName string

		if rows.Next() == false {
			break
		}

		rows.Scan(&CourseName)

		courseStruct, err := BuildCourseStruct(CourseName)

		if err != nil {
			return courses, err
		}

		courses = append(courses, courseStruct)
	}

	return courses, nil
}

// Nathan
// Returns a course struct describing the course a student is enrolled in, active or not
func LoadStudentCourse(UserID int) (CourseInfo, error) {

	var course CourseInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return course, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select CourseName from StudentCourses where Student=?", UserID)

	if err != nil {
		return course, errors.New("Query error.")
	}

	if rows.Next() == false {
		return course, errors.New("Student is not enrolled in a course.")
	}

	var CourseName string

	rows.Scan(&CourseName)

	course, err = BuildCourseStruct(CourseName)

	if err != nil {
		return course, err
	}

	return course, nil
}

// Nathan
// Returns all the assignments for a course as a slice
func LoadAssignments(CourseName string) ([]AssignmentInfo, error) {

	var assignments []AssignmentInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return assignments, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select AssignmentName from Assignments where CourseName=?", CourseName)

	if err != nil {
		return assignments, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var AssignmentName string

		if rows.Next() == false {
			break
		}

		rows.Scan(&AssignmentName)

		assignmentStruct, err := BuildAssignmentStruct(AssignmentName, CourseName)

		if err != nil {
			return assignments, err
		}

		assignments = append(assignments, assignmentStruct)
	}

	return assignments, nil

}

//Hannah
//returns the last submission for a particular student in a certain course for a certain assignment
func LoadLastSubmission(student int, CourseName string, AssignmentName string) (SubmissionInfo, error) {

	var submission SubmissionInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return submission, errors.New("No connection")
	}

	defer db.Close()

	//I know this seems unecessarily long, but it is what I have to do to survive
	rows, err := db.Query("select * from Submissions where student=? and CourseName=? and AssignmentName=? and SubmissionNumber=(Select Max(SubmissionNumber)from Submissions where student=? and CourseName=? and AssignmentName=?)", student, CourseName, AssignmentName, student, CourseName, AssignmentName)

	if err != nil {
		return submission, errors.New("Query error.")
	}

	//should just be able to call buildSubmissionStruct function, but not sure
	//how to do this with specific submission number
	if rows.Next() == false {
		return submission, errors.New("Invalid submission.")
	} else {
		rows.Scan(&submission.CourseName, &submission.AssignmentName,
			&submission.StudentUserID, &submission.Grade, &submission.Comments,
			&submission.Compiled, &submission.Results, &submission.SubmissionNum)
	}

	return submission, nil
}

//Hannah
func LoadStudentsInCourse(CourseName string) ([]UserInfo, error) {

	var users []UserInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return users, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("SELECT Username from Users where UserID IN (SELECT Student from StudentCourses where CourseName = ?)", CourseName)

	if err != nil {
		return users, errors.New("Query error.")
	}

	for i := 0; ; i++ {
		var user string

		if rows.Next() == false {
			break
		}

		rows.Scan(&user)

		userStruct, err := BuildUserStruct(user)

		if err != nil {
			return users, err
		}

		users = append(users, userStruct)
	}

	return users, nil
}

//returns all users in the system for admin use
//Hannah
func LoadAllUsers() ([]UserInfo, error) {

	var users []UserInfo

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return users, errors.New("No connection")
	}

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

		userStruct, err := BuildUserStruct(user)

		if err != nil {
			return users, err
		}

		users = append(users, userStruct)
	}

	return users, nil
}

// returns some kind of Grades for a user, does this need a struct?
// func loadGrades(student int, CourseName string) (??, error) {}

/*
func main() {


		course, err := buildAssignmentStruct("0", "TerwilligerCS15501SP17")

		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("%+v\n", course)


	x, err := loadAllUsers()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%v+\n\n\n\n", x)

	x, err = loadStudentsInCourse("TerwilligerCS15501SP17")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%v+\n\n\n\n\n", x)

	z, err := loadSubmission(10006, "JerkinsCS15502SP17", "1")

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%v+\n", z)

	y, err := loadStudentCourse(10004)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%v+\n", y)

	return

}
*/
