package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Log the HTTP request
func logHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
}

// mockHandler responds with "ok" as the response body
func mockHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok\n")
}

func headerFunc(w http.ResponseWriter, r *http.Request) {
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			fmt.Printf("%s=%s\n", k, v[0])
		}
	}

	r.ParseForm()
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			fmt.Printf("%s=%s\n", k, v[0])
		}
	}
}

// rootHandler used to process all inbound HTTP requests
func rootHandler(w http.ResponseWriter, r *http.Request) {
	logHandler(w, r)
	headerFunc(w, r)
	mockHandler(w, r)
}

type API struct {
	Port int `json:"port"`
}

var api API

// Start an HTTP server which dispatches to the rootHandler
func main() {
	raw, err := ioutil.ReadFile("./api.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(raw, &api)
	if err != nil {
		log.Fatal(" ", err)
	}

	port := strconv.Itoa(api.Port) //"3000"

	http.HandleFunc("/", rootHandler)

	log.Printf("server is listening on %v\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
