// Test
// April 29, 2017
// Testing each method and testing that every error message appears accordingly
// Todd Gibson
//------------------------------------------------------------------------------
package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

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

const DB_USER_NAME string = "dbadmin"
const DB_PASSWORD string = "EX0evNtl"
const DB_NAME string = "pest"
const SHELL_NAME string = "ksh"

func main() {
	//sendRandomPassword("tgibson12")

	//createUser("Todd", "G", "Gibson", "todd", 1, "JerkinsCS15502SP17")
	//user := UserInfo{}
	//user.UserName = tggibson

	//editUser(10114, "Todd", "J", "Gibson", 1)
	//deleteUser(10116)

	//changePassword("tggibson123", "todd")
	changePassword("nhuckaba555", "todd")
	//$2a$10$Ip9.oQug99IV8meAHC47k.x/0CYN0XoEQEArd.WZPek/KxLQOd1E6  -- old password
	//$2a$10$oWgso1F9e.XswsocaatfUOuGnK5nY0XUgFAp/jPpFlBkX3d3sxBg6  -- changed password
}

//-------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------
func changePassword(userName string, newPassword string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		return errors.New("No connection to DB")
	}

	defer db.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		return errors.New("Error encrypting password")
	}

	res, err := db.Exec("update Users set Password=? where UserName=?", hashedPassword, userName)

	if err != nil {
		return errors.New("DB error")
	}

	rowsChanged, err := res.RowsAffected()

	if (err != nil) || (rowsChanged != 1) {
		fmt.Println("user password not changed, could not find user.")
		return errors.New("Could not change password. User not found.")
	}

	return nil
}

//-------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------
func deleteUser(userID int) error {
	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("delete from Users where UserID=? and not exists(select 1 from StudentCourses where Student=? limit 1)", userID, userID)

	if err != nil {
		fmt.Println("student enrolled in class")
		return errors.New("User is currently enrolled in a class. Please remove the student from the class before deleting the user.")
	}
	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("Query didn't match any users. deleteUser")
		return errors.New("Query didn't match any users.")
	}

	return nil
}
func editUser(userID int, firstName string, MI string, lastName string, privLevel int) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)
	if err != nil {
		return errors.New("No connection")
	}

	res, err := db.Exec("update Eileen set FirstName=?, MiddleInitial=?, LastName=?, PrivLevel=? where UserID=?", firstName, MI, lastName, privLevel, userID)

	if err != nil {
		fmt.Println("ERROR!!! failed to update user!")
		return errors.New("User update failed.")
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected != 1 {
		fmt.Println("ERROR!!! query didn't match any users!")
		return errors.New("Query didn't match any users.")
	}

	return nil
}

//---------------------------------------------------------------------------
//Inputs: user's first name, user's middle initial, user's last name,
//	username, password, user's priv level
//Outputs: returns errors if the user cannot be created
//Written By: Hannah Hopkins and Brad Lanford
//Purpose: This function will be used by the instructor or administrator to
//	create a user. It will insert the user in the Users table in the
//	database.
//---------------------------------------------------------------------------
func createUser(firstName string, MI string, lastName string, username string, privLevel int, courseName string) error {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		fmt.Println("error 1")
		return errors.New("No connection")
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		return errors.New("failed to connect to database")
	}

	password := "todd2"                                                                      // TGG - added standard password for testing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //TGG

	if err != nil {
		fmt.Println("error 2")
		return errors.New("Error")
	}

	//added password																						//added ? rather than 'default'								//added hashedPassword
	_, err = db.Exec("INSERT INTO Users(FirstName, MiddleInitial, LastName, Username, Password, PrivLevel) VALUES(?, ?, ?, ?, ?, ?)", firstName, MI, lastName, username, hashedPassword, privLevel)

	if err != nil {
		fmt.Println("error 3, user creation failed")
		return errors.New("User creation failed.")
	}

	//sendRandomPassword(username)
	user, _ := BuildUserStruct(username)

	_, err = db.Exec("INSERT INTO StudentCourses (Student, CourseName) VALUES (?, ?)", user.UserID, courseName)

	if err != nil {
		fmt.Println("error 4 User unable to be added to the student courses." + username + courseName + "F")
		return errors.New("User unable to be added to the student courses." + username + courseName + "F")
	}

	/*_, err = db.Exec("INSERT INTO GradeReport" + courseName + "(Student) VALUES(select UserID from users where Username=" + username + ")")
	if err != nil {
		return errors.New("User unable to be added to GradeReport table.")
	}*/

	return nil
}

// Hannah
func BuildUserStruct(username string) (UserInfo, error) {
	user := UserInfo{}

	//user.UserName = username//

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@unix(/var/run/mysql/mysql.sock)/"+DB_NAME)

	if err != nil {
		fmt.Println("error 5")
		return user, errors.New("No connection")
	}

	defer db.Close()

	rows, err := db.Query("select UserID, FirstName, MiddleInitial, LastName, PrivLevel, LastLogin, PwdChangeFlag, NumLoginAttempts, Enabled from Users where Username = ?", username)

	if err != nil {
		fmt.Println("error 6")
		return user, errors.New("DB error")
	}

	if rows.Next() == false {
		return user, errors.New("Invalid User.")
		fmt.Println("error 6")
	} else {
		rows.Scan(&user.UserID, &user.FirstName, &user.MiddleInitial,
			&user.LastName, &user.PrivLevel, &user.LastLogin, &user.PwdChangeFlag,
			&user.NumLoginAttempts, &user.Enabled)
	}

	return user, nil

}
