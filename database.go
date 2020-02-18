package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	//_ "github.com/go-sql-driver/mysql"
)

//DBmetric - stuct for working with metrics in database
type DBmetric struct {
	db *sqlx.DB
}

func newDBmetric() *DBmetric {

	s := fmt.Sprintf("%v:%v@/%v", os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB"))

	database, err := sqlx.Open("postgress", s)

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
