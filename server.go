package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PostStruct struct {
	Buffer string
}

func serve(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePost(&w, r)
	case "GET":
		handleGet(&w, r)
	}
}

func handleGet(w *http.ResponseWriter, r *http.Request) {
	var buffer, dot, output string

	getJsonData(&buffer)

	output = "dot"

	dot, err := processBytes([]byte(buffer), &output)
	if err != nil {
		dot = err.Error()
	}

	fmt.Fprintf(*w, "GET request\nRequest struct = %v\n\nJSON data:\n%s\nDot:\n%s\n", r, buffer, dot)
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body '%s': %s", body, err)
		return
	}
}
