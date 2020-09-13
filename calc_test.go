package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// struct being used to define tests for handlers
type handlerTest struct {
	title          string
	routeAndParams string
	hasErr         bool
	expectedRes    int
	handler        func(http.ResponseWriter, *http.Request)
}

// Helper func to be used when testing handlers
func singleHandlerTestHelper(t *testing.T, test handlerTest) {
	req, err := http.NewRequest("GET", test.routeAndParams, nil)
	if test.hasErr == false && err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(test.handler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if test.hasErr == false && rr.Body.String() != strconv.Itoa(test.expectedRes) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), test.expectedRes)
	}
}

func TestStartHTTPServer(t *testing.T) {
	serverDoneWG.Add(1)
	addr := "127.0.0.1:8000"
	server = StartHTTPServer(addr, &serverDoneWG)

	// Wait for other goroutine to start the http server
	// No need for adding a channel, since it is a simple test and has nothing with internals of StartHTTPServer func.
	for server == nil {
		time.Sleep(time.Microsecond) // OK, does not make CPU busy
	}

	for _, route := range HealthEndPoints {
		req, err := http.NewRequest("GET", route, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HealthCheckHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `{"healhty": true}`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}

	// Following two lines: Not really neccessary, since runtime will shutdown the goroutines
	// StopHTTPServer()
	// serverDoneWG.Wait()
}

func TestAddHandler(t *testing.T) {
	tests := []handlerTest{
		{
			title:          "Add OK",
			routeAndParams: "/add?x=1&y=2",
			hasErr:         false,
			expectedRes:    3,
			handler:        AddHandler,
		},
		{
			title:          "Add with Slash",
			routeAndParams: "/add/?x=1&y=2",
			hasErr:         false,
			expectedRes:    3,
			handler:        AddHandler,
		},
		{
			title:          "Add expecting Error",
			routeAndParams: "/add?x=-1.99&y=2",
			hasErr:         true,
			expectedRes:    0,
			handler:        AddHandler,
		},
		{
			title:          "Add OK: No overflow/underflow checking",
			routeAndParams: fmt.Sprintf("/add?x=%v&y=%v", int(^uint(0)>>1), int(^uint(0)>>1)),
			hasErr:         false,
			expectedRes:    -2, // assumed 64-bit system
			handler:        AddHandler,
		},
		// and so on...
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			singleHandlerTestHelper(t, test)
		})
	}
}
func TestSubtractHandler(t *testing.T) {
	tests := []handlerTest{
		{
			title:          "Subtract OK",
			routeAndParams: "/subtract?x=1&y=2",
			hasErr:         false,
			expectedRes:    -1,
			handler:        SubtractHandler,
		},
		{
			title:          "Subtract with Slash",
			routeAndParams: "/subtract/?x=1&y=2",
			hasErr:         false,
			expectedRes:    -1,
			handler:        SubtractHandler,
		},
		{
			title:          "Subtract expecting Error",
			routeAndParams: "/subtract?x=-1.99&y=2",
			hasErr:         true,
			expectedRes:    0,
			handler:        SubtractHandler,
		},
		{
			title:          "Subtract OK: No overflow/underflow checking",
			routeAndParams: fmt.Sprintf("/subtract?x=%v&y=%v", -int(^uint(0)>>1)-1, int(^uint(0)>>1)),
			hasErr:         false,
			expectedRes:    1, // assumed 64-bit system
			handler:        SubtractHandler,
		},
		// and so on...
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			singleHandlerTestHelper(t, test)
		})
	}
}

func TestCastVarsOK(t *testing.T) {
	wantX, wantY := -10, 10
	if gotX, gotY, gotErr := castVars("-10", "10"); gotX != wantX && gotY != wantY && gotErr != nil {
		t.Errorf("castVars() = %v, %v, %v, want %v, %v, %v", gotX, gotY, gotErr, wantX, wantY, nil)
	}
}
func TestCastVarsErrorOnSpace(t *testing.T) {
	if gotX, gotY, gotErr := castVars("-10 ", " 10 "); gotErr == nil {
		t.Errorf(`castVars("-10 ", " 10 ")`+"= %v, %v, %v, want error", gotX, gotY, gotErr)
	}
}
func TestCastVarsErrorOnFloat(t *testing.T) {
	if gotX, gotY, gotErr := castVars("-10.2", "99.2"); gotErr == nil {
		t.Errorf(`castVars("-10.2", "99.2")`+"= %v, %v, %v, want error", gotX, gotY, gotErr)
	}
}
