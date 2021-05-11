package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/gddo/httputil/header"
)

type TokenType struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Split "Bearer XXXXXXXXXXXX" and return the token part
func SplitAuthToken(headerString string) string {
	tokenArr := strings.Split(headerString, " ")
	tokenString := ""
	if len(tokenArr) == 2 {
		tokenString = tokenArr[1]
	}
	return tokenString
}

// Create an Access Token
func AccessToken(userid string) (string, error) {
	os.Setenv("ACCESS_SECRET", "MYVERYSECRETKEY")

	token_claims := jwt.MapClaims{}
	token_claims["authorized"] = true
	token_claims["username"] = userid
	token_claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	acces_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token_claims)

	token, err := acces_token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Create a Refresh Token
func RefreshToken(userid string) (string, error) {
	os.Setenv("REFRESH_SECRET", "MYOTHERVERYSECRETKEY")

	token_claims := jwt.MapClaims{}
	token_claims["authorized"] = true
	token_claims["username"] = userid
	token_claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	acces_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token_claims)

	token, err := acces_token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Authorize users and return tokens
func GenerateTokens(username string, password string) (TokenType, error) {
	authenticated, authErr := CheckCreds(username, password)
	if authenticated {
		JWT, err := AccessToken(username)
		if err != nil {
			return TokenType{}, errors.New("internal server error while authenticating")
		}

		refTok, err := RefreshToken(username)
		if err != nil {
			return TokenType{}, errors.New("internal server error while authenticating")
		}

		return TokenType{
			AccessToken:  JWT,
			RefreshToken: refTok,
		}, err
	}
	return TokenType{}, authErr
}

// Verify an Access Token
func VerifyToken(tokenString string) (jwt.MapClaims, bool, error) {
	// Validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return jwt.MapClaims{}, false, fmt.Errorf("authentication error: %v", err)
	}

	// Extract metadata from token
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		// Check for username field
		_, ok := claims["username"].(string)
		if !ok {
			return jwt.MapClaims{}, false, errors.New("field username not found in authorization token")
		}
		return claims, true, nil
	} else {
		return jwt.MapClaims{}, false, errors.New("unauthorized")
	}
}

// Handle login requests
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if content type is "application/json"
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// Read a maximum of 1MB from body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Create a JSON decoder and decode the request JSON
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var loginData LoginType
	err := dec.Decode(&loginData)

	// If any error occurred during the decoding, send an appropriate response
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		// Return errors based on what error JSON parser returned
		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			http.Error(w, msg, http.StatusBadRequest)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			http.Error(w, msg, http.StatusBadRequest)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)

		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Decode it and check for an external JSON error
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// After checking for any errors, log the user in, and generate tokens
	tokenData, err := GenerateTokens(loginData.Username, loginData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	// Send the refresh token in a HTTPOnly cookie
	c := http.Cookie{
		Name:     "jid",
		Value:    tokenData.RefreshToken,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &c)

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	// Send the access token in JSON
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken: tokenData.AccessToken,
	})

}
