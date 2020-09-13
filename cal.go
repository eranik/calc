package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const addr = "127.0.0.1:8000"

// HealthEndPoints health check endpoints for api
var HealthEndPoints []string

var serverDoneWG sync.WaitGroup
var server *http.Server

func init() {
	HealthEndPoints = []string{
		"/health",
		"/health/",
		"/add/health",
		"/add/health/",
		"/subtract/health",
		"/subtract/health/",
	}
	server = nil
}

func castVars(s1, s2 string) (x int, y int, err error) {
	x, xErr := strconv.Atoi(s1)
	y, yErr := strconv.Atoi(s2)
	err = nil
	if xErr != nil {
		err = xErr
	} else if yErr != nil {
		err = yErr
	}
	return
}

// AddHandler handles /add?x=1&y=1 and /add/?x=1&y=1 to add two int values
// returns x + y as string
func AddHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	x, y, err := castVars(r.FormValue("x"), r.FormValue("y"))

	if err != nil {
		fmt.Fprintf(w, "error parsing values")
	} else {
		res := x + y
		fmt.Fprintf(w, strconv.Itoa(res))
	}
}

// SubtractHandler handles /subtract?x=1&y=1 and /subtract/?x=1&y=1 to subtract two int values
// returns x - y as string
func SubtractHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	x, y, err := castVars(r.FormValue("x"), r.FormValue("y"))

	if err != nil {
		fmt.Fprintf(w, "error parsing values")
	} else {
		res := x - y
		fmt.Fprintf(w, strconv.Itoa(res))
	}
}

// SubtractNotFoundHandler comment
func SubtractNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"help text": "Please provide correct query parameters: e.g. /subtract?x=20&y=10"}`)
}

// HealthCheckHandler used to check endpoints health; simple implementation
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"healhty": true}`)
}

// StartHTTPServer comment
func StartHTTPServer(addr string, wg *sync.WaitGroup) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/add", AddHandler).Queries("x", "{x:[0-9]+}", "y", "{y:[0-9]+}")
	router.HandleFunc("/add/", AddHandler).Queries("x", "{x:[0-9]+}", "y", "{y:[0-9]+}")

	router.HandleFunc("/subtract", SubtractHandler).Queries("x", "{x:[0-9]+}", "y", "{y:[0-9]+}")
	router.HandleFunc("/subtract/", SubtractHandler).Queries("x", "{x:[0-9]+}", "y", "{y:[0-9]+}")

	// Health
	for _, endpoint := range HealthEndPoints {
		router.HandleFunc(endpoint, HealthCheckHandler)
	}

	router.Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    addr,
		// Enforcing timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

// StopHTTPServer comment
func StopHTTPServer() {
	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

func main() {
	serverDoneWG.Add(1)
	server = StartHTTPServer(addr, &serverDoneWG)
	serverDoneWG.Wait()
}
