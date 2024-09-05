package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	http.HandleFunc("/form", formHandler)

	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting server at port 8080\n")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func formHandler(response http.ResponseWriter, request *http.Request) {
	url := request.URL.Path
	method := request.Method
	formParseErr := request.ParseForm()

	if url != "/form" || method != "POST" {
		http.Error(response, "404 Not Found", http.StatusNotFound)
		return
	}

	if formParseErr != nil {
		fmt.Fprintf(response, "ParseForm() err: %v", formParseErr)
		return
	}

	fmt.Fprintln(response, "POST requset successful")

	name := request.FormValue("name")
	fmt.Fprintf(response, "your name is: %s\n", name)
}

func helloHandler(response http.ResponseWriter, request *http.Request) {
	url := request.URL.Path
	method := request.Method
	if url != "/hello" || method != "GET" {
		http.Error(response, "404 Not Found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(response, "Hello from GO Server")
}
