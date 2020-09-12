package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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

// AddHandler handles /add/{x}/{y} and /add/{x}/{y}/ routes for adding two int values of x and y
// returns x + y as string result
func AddHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	x, y, err := castVars(vars["x"], vars["y"])

	if err != nil {
		fmt.Fprintf(w, "error parsing values")
	} else {
		res := x + y
		fmt.Fprintf(w, strconv.Itoa(res))
	}
}

// SubtractHandler handles /subtract/{x}/{y} and /subtract/{x}/{y}/ routes for subtract two int values of x and y
// returns x - y as string result
func SubtractHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	x, y, err := castVars(vars["x"], vars["y"])

	if err != nil {
		fmt.Fprintf(w, "error parsing values")
	} else {
		res := x - y
		fmt.Fprintf(w, strconv.Itoa(res))
	}
}

func setupRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/add/{x}/{y}", AddHandler)
	router.HandleFunc("/add/{x}/{y}/", AddHandler)
	router.HandleFunc("/subtract/{x}/{y}", SubtractHandler)
	router.HandleFunc("/subtract/{x}/{y}/", SubtractHandler)
	router.Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Enforcing timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func main() {
	setupRouter()
}
