package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//DBmetric - stuct for working with metrics in database
type DBmetric struct {
	db *sql.DB
}

func newDBmetric() *DBmetric {

	s := fmt.Sprintf("%v:%v@/%v", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS"), os.Getenv("MYSQL_DB"))

	database, err := sql.Open("mysql", s)

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

	err := json.NewDecoder(r.Body).Decode(&metrics)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = database.db.Exec(
		"INSERT INTO metrics(ServiceName,MetricName,MetricValue) VALUES(?,?,?)",
		metrics.ServiceName, metrics.MetricName, metrics.MetricValue)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)

}
