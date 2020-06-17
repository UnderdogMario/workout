package backend

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Mario()  {
	router := mux.NewRouter()
	const port string = ":8000"
	router.HandleFunc("/", DefaultHandler)
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(port, router))
}