package main

import (
	"dwitter_go_graphql/prisma/db"
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

// Split "xyz=AAAAAAA" and return the AAAAAAA part
func SplitCookie(cookieString string) string {
	arr := strings.Split(cookieString, "=")
	val := ""
	if len(arr) == 2 {
		val = arr[1]
	}
	return val
}

// Create an Access Token
func AccessToken(userID string) (string, error) {
	_, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Exec(ctx)
	if err == db.ErrNotFound {
		return "", errors.New("user doesn't exist")
	}

	token_claims := jwt.MapClaims{}
	token_claims["authorized"] = true
	token_claims["username"] = userID
	token_claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	acces_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token_claims)

	token, err := acces_token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Create a Refresh Token
func RefreshToken(userID string) (string, error) {
	userDB, err := client.User.FindUnique(
		db.User.Username.Equals(userID),
	).Exec(ctx)
	if err == db.ErrNotFound {
		return "", errors.New("user doesn't exist")
	}

	token_claims := jwt.MapClaims{}
	token_claims["authorized"] = true
	token_claims["username"] = userID
	token_claims["token_version"] = userDB.TokenVersion
	token_claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, token_claims)

	token, err := access_token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
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
			return TokenType{}, err
		}

		refTok, err := RefreshToken(username)
		if err != nil {
			return TokenType{}, err
		}

		return TokenType{
			AccessToken:  JWT,
			RefreshToken: refTok,
		}, err
	}
	return TokenType{}, authErr
}

// Verify an Access Token
func VerifyAccessToken(tokenString string) (jwt.MapClaims, bool, error) {
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
			return jwt.MapClaims{}, false, errors.New("field username not found in access token")
		}
		_, err = client.User.FindUnique(
			db.User.Username.Equals(claims["username"].(string)),
		).Exec(ctx)
		if err == db.ErrNotFound {
			return jwt.MapClaims{}, false, errors.New("user doesn't exist")
		}
		return claims, true, nil
	} else {
		return jwt.MapClaims{}, false, errors.New("unauthorized")
	}
}

// Verify a Refresh Token
func VerifyRefreshToken(tokenString string) (jwt.MapClaims, bool, error) {
	// Validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return jwt.MapClaims{}, false, fmt.Errorf("authentication error: %v", err)
	}

	// Extract metadata from token
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		// Check for username field
		username, ok := claims["username"].(string)
		if !ok {
			return jwt.MapClaims{}, false, errors.New("field username not found in refresh token")
		}
		// Check for token_version field
		tokenV, ok := claims["token_version"].(float64)
		if !ok {
			return jwt.MapClaims{}, false, errors.New("field token_version not found in refresh token")
		}

		userDB, err := client.User.FindUnique(
			db.User.Username.Equals(username),
		).Exec(ctx)

		if err == db.ErrNotFound {
			return jwt.MapClaims{}, false, errors.New("user doesn't exist")
		}
		fmt.Printf("DB: %v, token: %v", userDB.TokenVersion, int(tokenV))
		if userDB.TokenVersion != int(tokenV) {
			return jwt.MapClaims{}, false, errors.New("invalid token version")
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
		return
	}

	// Send the refresh token in a HTTPOnly cookie
	c := http.Cookie{
		Name:     "jid",
		Value:    tokenData.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/refresh_token",
	}
	http.SetCookie(w, &c)

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	// Send the access token in JSON
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken: tokenData.AccessToken,
	})

}

// Handle login requests
func refreshHandler(w http.ResponseWriter, r *http.Request) {
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

	cookieString, err := r.Cookie("jid")
	if err != nil {
		msg := "Refresh Token not present"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	token := SplitCookie(cookieString.String())

	claims, verified, err := VerifyRefreshToken(token)
	if (err != nil) || (!verified) {
		msg := fmt.Sprintf("Unauthorized: %v", err)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	userID, ok := claims["username"].(string)
	if !ok {
		msg := "Invalid refresh token"
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	refTok, err := RefreshToken(userID)
	if err != nil {
		msg := "Invalid refresh token"
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	// Send the refresh token in a HTTPOnly cookie
	c := http.Cookie{
		Name:     "jid",
		Value:    refTok,
		HttpOnly: true,
		Secure:   true,
		Path:     "/refresh_token",
	}
	http.SetCookie(w, &c)

	accessTok, err := AccessToken(userID)
	if err != nil {
		msg := "Invalid refresh token"
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	// Send the access token in JSON
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken: accessTok,
	})
}
