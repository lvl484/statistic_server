package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	address := flag.String("address", ":1433", "address of server")
	flag.Parse()

	dbHandling := newDBmetric()
	defer dbHandling.db.Close()

	router := mux.NewRouter()
	router.Use(Middleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/{ServiceName}", dbHandling.MetricsCreate).Methods(http.MethodPost)
	if err := http.ListenAndServe(*address, router); err != nil {
		log.Fatal(err.Error())
	}

}
