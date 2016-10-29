package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"net/http"
	"math/rand"
	"time"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello", index).Methods("GET")
	router.HandleFunc("/health", health).Methods("GET")
	router.HandleFunc("/unhealthy", unhealthy).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /hello request")
	log.Println(r.UserAgent())
	w.WriteHeader(http.StatusOK)
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	forward := r.Header.Get("X-Forwarded-For")
	for _, addr := range addrs {
	    if ipv4 := addr.To4(); ipv4 != nil {
					fmt.Fprintf(w, "{\n\"Host\": \"%v\",\n\"IPv4\": \"%v\"\n,\"RemoteAddr\": \"%v\"\n,\"X-Forwarded-For\": \"%v\"\n}\n", host, ipv4, r.RemoteAddr, )
	    }
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /health request")
	log.Println(r.UserAgent())

	w.WriteHeader(http.StatusOK)
}

func unhealthy(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /unhealthy request")
	log.Println(r.UserAgent())
	t :=  time.Now()
	i := int64(t.Nanosecond())
	rand.Seed(i)
	answers := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
	}

	w.WriteHeader(answers[rand.Intn(len(answers))])
}
