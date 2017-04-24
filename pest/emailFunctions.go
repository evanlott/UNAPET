package main

import (
	"math/rand"
	"net/smtp"
)

// Nathan
func sendRandomPassword(form Request) (bool, string) {

	// verify user exists first..
	// select * from Users where UserName=form.userName

	// edit this to match her specifications
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$?"

	newPasswordLen := 12

	byteSlice := make([]byte, newPasswordLen)
	for i := range byteSlice {
		byteSlice[i] = chars[rand.Intn(len(chars))]
	}

	newPassword := string(byteSlice)

	emailUserName := "??"
	emailPassword := "??"
	emailServerAddr := "smtp.??"
	emailServerPort := "587"

	login := smtp.PlainAuth("", emailUserName, emailPassword, emailServerAddr)

	toEmail := form.userName + "@una.edu"
	fromEmail := "DoNotReply@??.com"

	content := []byte("Your new password is " + newPassword + ".\n\n" +
		"If you did not request a new password please contact system admin.\n\n" +
		"This is an automatically generated email.\n" +
		"Please do not reply to this email.\n\n" +
		"-UNAPET")

	err := smtp.SendMail(emailServerAddr+":"+emailServerPort, login, fromEmail, []string{toEmail}, content)

	if err != nil {
		return false, "Error. Email failed to send. Your password was not changed. Please try again."
	}

	// change user's password to the new one
	// call function that doesn't exist yet
	// err, _ := changePassword(userName, newPassword)
	// if err != nil {
	//		return false, "Database error. Old password could not be changed. You should ignore the new password sent to you by email and try again."
	// }

	return true, form.fromPage

}
