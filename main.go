package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", index).Methods("GET")
	router.HandleFunc("/health", index).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /health request")
	log.Println("Responsing to /health request")
	log.Println("Responsing to /health request")
	log.Println(r.UserAgent())
	rand.Seed(422)
	answers := []int{
		http.StatusContinue,
		http.StatusSwitchingProtocols,
	}

	w.WriteHeader(answers[rand.Intn(len(answers))])
	fmt.Fprintln(w, "Hello:", name)
}
