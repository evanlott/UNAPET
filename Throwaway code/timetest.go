// timetest
package main

import (
	"fmt"
	"time"
)

func main() {
	start := "January 28, 2015"
	end := "May 30, 2017"

	// The form must be January 2,2006.
	form := "January 2, 2006"

	// Parse the string according to the form.
	t1, _ := time.Parse(form, start)
	t2, _ := time.Parse(form, end)

	fmt.Println(t1.Format("01/02/2006 15:04:05"))
	fmt.Println(t2.Format("01/02/2006 15:04:05"))

	currentTime := time.Now()
	fmt.Println(currentTime.Format("01/02/2006 15:04:05"))

	if t1.Format("01/02/2006 15:04:05") < currentTime.Format("01/02/2006 15:04:05") && t2.Format("01/02/2006 15:04:05") > currentTime.Format("01/02/2006 15:04:05") {
		fmt.Println("data comparison worked.")
	}
}
