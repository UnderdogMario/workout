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

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println(request)
	setupResponse(&writer, request)
	if (*request).Method == "OPTIONS" {
		return
	}

	request.ParseForm()

	var email string
	var password string
	for key, value := range request.Form {
		fmt.Println(key, value)
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

