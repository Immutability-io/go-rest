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

package main

import (
    "log"
    "net/http"
)

func Run(addr string, sslAddr string, ssl map[string]string) chan error {

    errs := make(chan error)

		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/auth", auth).Methods("GET")
		router.HandleFunc("/login", login).Methods("POST")
		router.HandleFunc("/hello", hello).Methods("GET")
		router.HandleFunc("/health", health).Methods("GET")
		router.HandleFunc("/unhealthy", unhealthy).Methods("GET")

    // Starting HTTP server
    go func() {
        log.Printf("Staring HTTP service on %s ...", addr)

        if err := http.ListenAndServe(addr, router); err != nil {
            errs <- err
        }

    }()

    // Starting HTTPS server
    go func() {
        log.Printf("Staring HTTPS service on %s ...", addr)
        if err := http.ListenAndServeTLS(sslAddr, ssl["cert"], ssl["key"], router); err != nil {
            errs <- err
        }
    }()

    return errs
}

func main() {

    errs := Run(":8080", ":443", map[string]string{
        "cert": "/etc/ssl/service.crt",
        "key":  "/etc/ssl/service.key",
    })

    // This will run forever until channel receives error
    select {
    case err := <-errs:
        log.Printf("Could not start serving service due to (error: %s)", err)
    }

}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /hello request")
	log.Println(r.UserAgent())
	w.WriteHeader(http.StatusOK)
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	forward  := r.Header.Get("X-Forwarded-For")
	for _, addr := range addrs {
	    if ipv4 := addr.To4(); ipv4 != nil {
					fmt.Fprintf(w, "{\n\"Version\": \"v0.0.8\",\n\"Host\": \"%v\",\n\"IPv4\": \"%v\",\n\"RemoteAddr\": \"%v\",\n\"X-Forwarded-For\": \"%v\"\n}\n", host, ipv4, r.RemoteAddr, forward)
	    }
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /health request")
	log.Println(r.UserAgent())

	w.WriteHeader(http.StatusOK)
}

func auth(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /auth request")
	log.Println(r.UserAgent())
	var cookie,err = r.Cookie("IMMUTABILITY_SSO")
	if err == nil {
		if cookie.Value == "supersecret" {
			w.WriteHeader(http.StatusOK)
		}	else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}	else {
		w.WriteHeader(http.StatusUnauthorized)
	}


}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /login request")
	log.Println(r.UserAgent())
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie    :=    http.Cookie{Name: "IMMUTABILITY_SSO", Value:"supersecret", Domain: ".immutability.io", Expires:expiration}
	http.SetCookie(w, &cookie)

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
