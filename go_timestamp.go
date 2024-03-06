package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TimeApi struct {
	Unix int    `json:"unix"`
	Utc  string `json:"utc"`
}

type Error struct {
	Error string `json:"error"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// if the request only  has numbers after the slash
	fmt.Fprintln(w, "Welcome to the timestamp microservice")
}

func handleTimeApi(w http.ResponseWriter, r *http.Request) {
	// if the request only  has numbers after the slash
	inputTime := strings.TrimPrefix(r.URL.Path, "/api/")
	t, err := validateTime(inputTime)
	if err != nil {
		errorResponse := map[string]string{"error": "Invalid Date"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	timeOutput := TimeApi{
		Unix: int(t.Unix()),
		Utc:  t.UTC().String(),
	}
	json.NewEncoder(w).Encode(timeOutput)
}

func main() {
	port := ":8080"
	fmt.Println("Starting server on port: ", port)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/api/", handleTimeApi)

	http.ListenAndServe(port, nil)

}

// UTILS

func validateTime(s string) (time.Time, error) {
	unixTime, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		// Convert Unix timestamp to time.Time
		return time.Unix(unixTime, 0), nil
	}

	pattern := `^\d{4}-\d{2}-\d{2}$`
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return time.Time{}, err
	}
	if !matched {
		return time.Time{}, fmt.Errorf("Invalid input. Use YYYY-MM-DD or Unix timestamp")
	}

	// Parse the input string as a date
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
