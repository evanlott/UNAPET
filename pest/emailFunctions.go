package main

import (
	"errors"
	"math/rand"
	"net/smtp"
	"time"
)

// Nathan
func sendRandomPassword(userName string) error {

	// verify user exists first..
	// select * from Users where UserName=form.userName

	// edit this to match her specifications
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$?"

	newPasswordLen := 12

	rand.Seed(time.Now().UTC().UnixNano())

	byteSlice := make([]byte, newPasswordLen)
	for i := range byteSlice {
		byteSlice[i] = chars[rand.Intn(len(chars))]
	}

	newPassword := string(byteSlice)

	// constants from the config file
	emailUserName := "univnorthalabamapet@gmail.com"
	emailPassword := "W33kahell0"
	emailServerAddr := "smtp.gmail.com"
	emailServerPort := "587"

	login := smtp.PlainAuth("", emailUserName, emailPassword, emailServerAddr)

	toEmail := userName + "@una.edu"
	fromEmail := "univnorthalabamapet@gmail.com"

	content := []byte("Your new password is " + newPassword + ".\n\n" +
		"If you did not request a new password please contact system admin.\n\n" +
		"This is an automatically generated email.\n" +
		"Please do not reply to this email.\n\n" +
		"-UNAPET")

	err := smtp.SendMail(emailServerAddr+":"+emailServerPort, login, fromEmail, []string{toEmail}, content)

	if err != nil {
		return errors.New("Error. Email failed to send. Your password was not changed. Please try again.")
	}

	// change user's password to the new one
	// call function that doesn't exist yet
	err = changePassword(userName, newPassword)

	if err != nil {
		return errors.New("Database error. Old password could not be changed. You should ignore the new password sent to you by email and try again.")
	}

	// set this user's change password flag

	return nil

}
