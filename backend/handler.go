package backend

import (
	"fmt"
	"net/http"
)

type User struct {
	Email string `json: "email"`
	Password string `json: "password"`
}

func DefaultHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Default, Up And Running, Good Connection")
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	var email string
	var password string
	for key, value := range request.Form {
		if key == "email" {
			email = value[0]
		}
		if key == "password" {
			password = value[0]
		}
	}
	if email == "" || password == "" {
		fmt.Fprint(writer, "Failure, You pass bad data")
	} else {
		result := ValidateUserInformation(email, password)
		if result {
			fmt.Fprint(writer, "You successful login")
		} else {
			fmt.Fprint(writer, "Your information is incorrect")
		}
	}
}

