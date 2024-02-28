package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)

	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
