package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Listening on port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	err := routeErrorHandling(w, r, []string{"GET"}, "/hello")
	if err != nil {
		return
	}

	fmt.Printf("Hello there!\n")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	err := routeErrorHandling(w, r, []string{"POST"}, "/form")
	if err != nil {
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "POST request success")
	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Fprintf(w, "Form name=%s address=%s", name, address)
}

func routeErrorHandling(w http.ResponseWriter, r *http.Request, allowedMethods []string, expectedPath string) error {
	if r.URL.Path != expectedPath {
		http.Error(w, "404 not found.", http.StatusNotFound)
	}

	if !slices.Contains(allowedMethods, r.Method) {

		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return errors.New("method not allowed")
	}

	return nil
}
