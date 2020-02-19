package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//DBmetric - stuct for working with metrics in database
type DBmetric struct {
	db *sql.DB
}

func newDBmetric() *DBmetric {

	s := fmt.Sprintf("host=localhost port=5432 user=%v password=%v dbname=%v sslmode=disable", os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB"))

	//s := fmt.Sprintf("%v:%v@/%v", os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB"))

	database, err := sql.Open("postgres", s)
	defer database.Close()

	err = database.Ping()
	if err != nil {
		log.Println(err)
	}

	if err != nil {
		log.Println(err)
	}

	return &DBmetric{
		db: database,
	}
}

//MetricsCreate record metris from endpoints end record it to db
func (database *DBmetric) MetricsCreate(w http.ResponseWriter, r *http.Request) {
	var metrics Metrics

	vars := mux.Vars(r)
	a := vars["ServiceName"]

	err := json.NewDecoder(r.Body).Decode(&metrics)
	fmt.Println(a, metrics.MetricValue, metrics.MetricName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = database.db.Exec(
		"INSERT INTO metrics(servicename,metricvalue,metricname) VALUES(?,?,?)",
		a, metrics.MetricValue, metrics.MetricName)

	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)

}
