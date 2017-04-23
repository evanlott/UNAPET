package main

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// TODO func createUser( -some stuff- ) {}

//---------------------------------------------------------------------------
//Inputs: user ID number, user's first name, user's middle initial, 
//	user's last name, user's priv level 
//Outputs: returns errors if the user could not be updated 
//Written By: Hannah Hopkins and Nathan Huckaba 
//Purpose: This function will be used by the instructor or administrator to
//	edit a user's information. It will update the user in the 
//	Users table in the database.  
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: user ID number 
//Outputs: returns errors if the user could not be deleted
//Written By: Hannah Hopkins
//Purpose: This function will be used by the instructor or administrator to
//	delete a user. It will remove the user from the Users table in 
//	the database if the user is not currently associated with a 
//	course. If the user is in a course, the user will not be able to
//	be removed and an error will be generated.  
//---------------------------------------------------------------------------
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

//---------------------------------------------------------------------------
//Inputs: the CSV file name 
//Outputs: returns errors if the csv could not be imported 
//Written By: Tyler Delano, Eileen Drass, Hannah Hopkins, Nathan Huckaba
//	Evan Lott 
//Purpose: This function will be used by the administrator or instructor to
//	import a CSV file of students. 
//---------------------------------------------------------------------------
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
