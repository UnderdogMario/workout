package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// default handler to see the connection is running
func DefaultHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Default, Up And Running, Good Connection")
}

// this is setup to avoid CORS policy
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Expose-Headers", "Access-Token, X-Auth")
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

	success, userInfo := ValidateUserInformation(user.Email, user.Password)
	// validate the data passed in
	if !success{
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

	writer.Header().Set("X-Auth", sessionToken)
	json.NewEncoder(writer).Encode(userInfo)
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

// this is used to handle user update their information
func UserProfileUpdateHandler(writer http.ResponseWriter, request *http.Request)  {
	setupResponse(&writer, request)
	if (*request).Method == "OPTIONS" {
		return
	}
	sid := request.Header.Get("X-Auth")
	if !ValidateSessionID(sid) {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	var userInfo UserInfo
	json.NewDecoder(request.Body).Decode(&userInfo)
	UserProfileUpdate(userInfo)
}