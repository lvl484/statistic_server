package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", ExampleHandler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

	// Double check it's a post request being made
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid_http_method")
		return
	}

	for i := 1; i < 3; i++ {
		// Must call ParseForm() before working with data
		r.ParseForm()

		// Log all data. Form is a map[]
		log.Println(r.Form)

		for key, values := range r.Form { //range over map
			for _, value := range values { //range over [string]
				fmt.Println(key, value)
			}
		}
	}
}
