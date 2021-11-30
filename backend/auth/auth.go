// Package auth provides functions useful for using authentication in this API.
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/soumitradev/Dwitter/backend/common"
	"github.com/soumitradev/Dwitter/backend/prisma/db"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/gddo/httputil/header"
)

// A tokenType stores an access and refresh token
type tokenType struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// A loginResponse stores the response to authentication
type loginResponse struct {
	AccessToken string `json:"accessToken"`
}

// A loginType stores login info
type loginType struct {
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
func splitCookie(cookieString string) string {
	arr := strings.Split(cookieString, "=")
	val := ""
	if len(arr) == 2 {
		val = arr[1]
	}
	return val
}

// Generate an Access Token
func generateAccessToken(username string) (string, error) {
	// Check if user exists
	_, err := common.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return "", errors.New("user doesn't exist")
	}

	// Save data in claims and generate token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["authorized"] = true
	tokenClaims["username"] = username
	tokenClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	token, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Generate a Refresh Token
func generateRefreshToken(username string) (string, error) {
	// Check if user exists
	userDB, err := common.Client.User.FindUnique(
		db.User.Username.Equals(username),
	).Exec(common.BaseCtx)
	if err == db.ErrNotFound {
		return "", errors.New("user doesn't exist")
	}

	// Save data in claims and generate token
	tokenClaims := jwt.MapClaims{}
	tokenClaims["authorized"] = true
	tokenClaims["username"] = username
	tokenClaims["token_version"] = userDB.TokenVersion
	tokenClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	token, err := accessToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Authorize user and return tokens
func generateTokens(username string, password string) (tokenType, error) {
	authenticated, authErr := common.CheckCreds(username, password)
	if authenticated {
		JWT, err := generateAccessToken(username)
		if err != nil {
			return tokenType{}, err
		}

		refTok, err := generateRefreshToken(username)
		if err != nil {
			return tokenType{}, err
		}

		return tokenType{
			AccessToken:  JWT,
			RefreshToken: refTok,
		}, err
	}
	return tokenType{}, authErr
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
		_, err = common.Client.User.FindUnique(
			db.User.Username.Equals(claims["username"].(string)),
		).Exec(common.BaseCtx)
		if err == db.ErrNotFound {
			return jwt.MapClaims{}, false, errors.New("user doesn't exist")
		}
		return claims, true, nil
	} else {
		return jwt.MapClaims{}, false, errors.New("unauthorized")
	}
}

// Verify a Refresh Token
func verifyRefreshToken(tokenString string) (jwt.MapClaims, bool, error) {
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

		userDB, err := common.Client.User.FindUnique(
			db.User.Username.Equals(username),
		).Exec(common.BaseCtx)

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

// Handles login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if content type is "application/json"
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: "Content-Type header is not application/json",
			})
			return
		}
	}

	// Read a maximum of 1MB from body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Create a JSON decoder and decode the request JSON
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var loginData loginType
	err := decoder.Decode(&loginData)

	// If any error occurred during the decoding, send an appropriate response
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		// Return errors based on what error JSON parser returned
		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: err.Error(),
			})
		}
		return
	}

	// Decode it and check for an external JSON error
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: msg,
		})
		return
	}

	// After checking for any errors, log the user in, and generate tokens
	tokenData, err := generateTokens(loginData.Username, loginData.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: err.Error(),
		})
		return
	}

	// Send the refresh token in a HTTPOnly cookie
	c := http.Cookie{
		Name:     "jid",
		Value:    tokenData.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/refresh_token",
	}
	http.SetCookie(w, &c)

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	// Send the access token in JSON
	json.NewEncoder(w).Encode(loginResponse{
		AccessToken: tokenData.AccessToken,
	})

}

// Handle refresh-token requests
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// Check if content type is "application/json"
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: "Content-Type header is not application/json",
			})
			return
		}
	}
	// Read a maximum of 1MB from body
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Create a JSON decoder and decode the request JSON
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var loginData loginType
	err := decoder.Decode(&loginData)

	// If any error occurred during the decoding, send an appropriate response
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		// Return errors based on what error JSON parser returned
		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: msg,
			})

		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(common.HTTPError{
				Error: err.Error(),
			})
		}
		return
	}

	// Decode it and check for an external JSON error
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: msg,
		})
		return
	}

	cookieString, err := r.Cookie("jid")
	if err != nil {
		msg := "Refresh Token not present"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: msg,
		})
		return
	}
	token := splitCookie(cookieString.String())

	claims, verified, err := verifyRefreshToken(token)
	if (err != nil) || (!verified) {
		msg := fmt.Sprintf("Unauthorized: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: msg,
		})
		return
	}

	userID, ok := claims["username"].(string)
	if !ok {
		msg := "Invalid refresh token"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: msg,
		})
		return
	}

	refTok, err := generateRefreshToken(userID)
	if err != nil {
		msg := "Invalid refresh token"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: msg,
		})
		return
	}

	// Send the refresh token in a HTTPOnly cookie
	c := http.Cookie{
		Name:     "jid",
		Value:    refTok,
		HttpOnly: true,
		Secure:   true,
		Path:     "/api/refresh_token",
	}
	http.SetCookie(w, &c)

	accessTok, err := generateAccessToken(userID)
	if err != nil {
		msg := "Invalid refresh token"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(common.HTTPError{
			Error: msg,
		})
		return
	}
	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	// Send the access token in JSON
	json.NewEncoder(w).Encode(loginResponse{
		AccessToken: accessTok,
	})
}

// Check header of request and authenticate
func Authenticate(authHeader string) (string, error) {
	tokenString := SplitAuthToken(authHeader)
	data, isAuth, err := VerifyAccessToken(tokenString)
	if (err != nil) || !isAuth {
		return "", errors.New("Unauthorized")
	}

	username := data["username"].(string)
	return username, nil
}
