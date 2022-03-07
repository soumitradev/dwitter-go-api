// Package auth provides functions useful for using authentication in this API.
package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/soumitradev/Dwitter/backend/common"
	"github.com/soumitradev/Dwitter/backend/prisma/db"
)

// TODO: Make the email look better
func SendVerificationEmail(emailID string, link string) (*rest.Response, error) {
	from := mail.NewEmail("Dwitter", os.Getenv("SENDGRID_SENDER_EMAIL_ADDR"))
	subject := "Dwitter account verification"
	to := mail.NewEmail("Recipient", emailID)
	plainTextContent := "Your verification link for Dwitter is: " + link + ".\nClick the link to verify your account.\nUnverified accounts are deleted after 1 hour."
	htmlContent := "Your verification link for Dwitter is: " + "<a href=\"" + link + "\">" + link + "</a>\r\n" + "Click the link to verify your account.\r\n" + "<strong>Unverified accounts are deleted after 1 hour.</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	response, err := common.SendgridClient.Send(message)
	return response, err
}

func DeleteUserAfterExpire(minutes int, token string) {
	time.Sleep(time.Minute * time.Duration(minutes))
	if username, found := common.AccountCreatedButNotVerified[token]; found {
		_, err := common.InternalDeleteUser(username)
		if err != nil {
			fmt.Printf("Error deleting user: %v", err)
		}
	}
}

// TODO: Make the thing actually send the emoji
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	if username, found := common.AccountCreatedButNotVerified[token]; found {
		delete(common.AccountCreatedButNotVerified, token)

		_, err := common.Client.User.FindUnique(
			db.User.Username.Equals(username),
		).Update(
			db.User.Verified.Set(true),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			res := fmt.Sprintf("user not found: %v", err)
			w.Write([]byte(res))
		}
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			res := fmt.Sprintf("internal server error: %v", err)
			w.Write([]byte(res))
		}

		verifiedHTML := "Account verified!\nYou may close this tab and sign in now."
		// Set the response headers
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(verifiedHTML))
	} else {
		unknownHTML := "Unrecognized verification link\nIs your account already verified? 🤔"

		// Set the response headers
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(unknownHTML))
	}
}
