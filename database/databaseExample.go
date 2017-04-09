package main // test.go

// http://www.wadewegner.com/2014/12/easy-go-programming-setup-for-windows/

// username: unapet
// pass: cs455

//package main

import (
	"database/sql"
	"fmt"
	//	"strconv"

	//_ "github.com/go-sql-driver/mysql"
)

const DB_USER_NAME string = "root"
const DB_PASSWORD string = "cs455"
const DB_NAME string = "UNAPET"

func handleError(err error) {
	if err != nil {
		panic("")
	}
}

func panicFunc() {
	panic("")
}

func queryDB() {

	db, err := sql.Open("mysql", DB_USER_NAME+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/"+DB_NAME) //root:cs455@/UNAPET
	defer db.Close()
	handleError(err)

	rows, err := db.Query("select Username, Email from Users")
	defer rows.Close()
	handleError(err)

	rows.Next()
	var name, email string
	rows.Scan(&name, &email)
	fmt.Println(name + " " + email)

	rows.Next()
	rows.Scan(&name, &email)
	fmt.Println(name + " " + email)

	panicFunc()
	fmt.Println("DB Program")
}

func recoverFromError(r error) bool {
	panicCalled := recover()

	if panicCalled != nil {
		fmt.Println("Someone called panic")
		return true
	} else {
		return false
	}
}

func main() {

	// recover from error template
	defer func() {
		r := recover()
		if r == nil {
			return
		} else {
			fmt.Println("Recovery code")
		}
	}()

	queryDB()
}
