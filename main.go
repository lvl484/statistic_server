package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	address := flag.String("address", ":80", "address of server")
	flag.Parse()

	dbHandling := newDBmetric()
	defer dbHandling.db.Close()

	router := mux.NewRouter()
	router.Use(Middleware)

	api := router.PathPrefix("/api/v1").Subrouter()
	//	api.HandleFunc("/{metric}", dbHandling.MetricsCreate).Methods(http.MethodGet)
	api.HandleFunc("/", dbHandling.MetricsCreate).Methods(http.MethodGet)

	if err := http.ListenAndServe(*address, router); err != nil {
		log.Fatal(err.Error())
	}

}
