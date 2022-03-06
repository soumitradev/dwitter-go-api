// Package auth provides functions useful for using authentication in this API.
package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/soumitradev/Dwitter/backend/common"
)

// TODO: Add some kind of map thingie for unverified users, and delete user accounts after an hour
// Also, some kind of API endpoint for verifying users
func SendVerificationEmail(emailID string) {
	from := mail.NewEmail("Dwitter", os.Getenv("SENDGRID_SENDER_EMAIL_ADDR"))
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Recipient", emailID)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	response, err := common.SendgridClient.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
