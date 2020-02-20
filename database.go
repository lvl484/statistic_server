package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	if err != nil {
		log.Println(err)
	}

	err = database.Ping()
	if err != nil {
		log.Println(err)
	}

	return &DBmetric{
		db: database,
	}
}

//CurrentTime ...
func CurrentTime() time.Time {
	return time.Now()
}

//MetricsCreate record metris from endpoints end record it to db
func (database *DBmetric) MetricsCreate(w http.ResponseWriter, r *http.Request) {

	var metrics Metrics

	err := json.NewDecoder(r.Body).Decode(&metrics)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	metrics.ServiceName = vars["ServiceName"]
	ct := CurrentTime()

	fmt.Println(ct)

	fmt.Println(metrics.ServiceName, metrics.MetricValue, metrics.MetricName, metrics.Status, metrics.Time)

	//	_, err = database.db.Exec("INSERT INTO metrics(servicename,metricvalue,metricname) VALUES(?,?,?)", vars["ServiceName"], metrics.MetricValue, metrics.MetricName)

	sqlStatement := `
	INSERT INTO metrics(servicename,metricvalue,metricname,status,time)
	VALUES($1,$2,$3,$4,$5) `

	_, err = database.db.Exec(sqlStatement, metrics.ServiceName, metrics.MetricValue, metrics.MetricName, metrics.Status, metrics.Time)

	if err != nil {
		fmt.Println("bad")
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)

}

//GetMetricsForService gets all metrics for the collection unit
func (database *DBmetric) GetMetricsForService(w http.ResponseWriter, r *http.Request) {

	var metrics Metrics

	vars := mux.Vars(r)
	metrics.ServiceName = vars["ServiceName"]

	rows := database.db.QueryRow("SELECT servicename,metricvalue, metricname FROM metrics WHERE servicename=?", metrics.ServiceName)

	err := rows.Scan(&metrics.ServiceName, &metrics.MetricValue, &metrics.MetricName)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(metrics.ServiceName, metrics.MetricValue, metrics.MetricName)
}

func (database *DBmetric) GetAllSuccessfullyHandled(w http.ResponseWriter, r *http.Request) {

	var metrics Metrics

	vars := mux.Vars(r)
	metrics.ServiceName = vars["ServiceName"]

	rows := database.db.QueryRow("SELECT servicename,metricvalue, metricname FROM metrics WHERE servicename=?", metrics.ServiceName)

	err := rows.Scan(&metrics.ServiceName, &metrics.MetricValue, &metrics.MetricName)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(metrics.ServiceName, metrics.MetricValue, metrics.MetricName)
}
