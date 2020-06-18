package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Email string `json: "email"`
	Password string `json: "password"`
}

// default handler to see the connection is running
func DefaultHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Default, Up And Running, Good Connection")
}

// this is setup to avoid CORS policy
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// this is used to handle user login, and will give the user a cookie, if the
// user password and email is correct
func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	var user User
	setupResponse(&writer, request)
	if (*request).Method == "OPTIONS" {
		return
	}

	err := json.NewDecoder(request.Body).Decode(&user)

	// the request body have something wrong
	if err != nil || user.Email == "" || user.Password == ""{
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate the data passed in
	if !ValidateUserInformation(user.Email, user.Password) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// generate session token
	sessionToken := CreateSessionToken(user.Email)

	// see if the session token is successfully generated
	if sessionToken == "" {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

// this is used to handle user register for a new account
func RegisterHandler(writer http.ResponseWriter, request *http.Request)  {
	var user User
	setupResponse(&writer, request)
	if (*request).Method == "OPTIONS" {
		return
	}
	err := json.NewDecoder(request.Body).Decode(&user)

	// the request body have something wrong
	if err != nil || user.Email == "" || user.Password == ""{
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// create new user for the client
	if !CreateNewUser(user.Email, user.Password) {
		writer.WriteHeader(http.StatusBadRequest)
		return
	} else {
		writer.WriteHeader(http.StatusAccepted)
		writer.Write([]byte("Success, Your account is created"))
		return
	}

}