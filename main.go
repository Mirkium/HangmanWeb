package main

import (
	"net/http"
)

func main() {

	FirstConnect = true

	Menu()
	Hangman()

	fileServer := http.FileServer(http.Dir("./assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	http.ListenAndServe("localhost:8080", nil)

}
