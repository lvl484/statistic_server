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

const (
	sourceFormat     = "host=localhost port=5432 user=%v password=%v dbname=%v sslmode=disable"
	serviceName      = "ServiceName"
	postgresqlQuerry = "SELECT servicename,metricvalue, metricname,status,time FROM metrics WHERE servicename=$1"
)

func newDBmetric() *DBmetric {

	s := fmt.Sprintf(sourceFormat, os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB"))

	database, err := sql.Open("postgres", s)

	if err != nil {
		database.Close()
		log.Fatal(err)
	}

	err = database.Ping()

	if err != nil {
		log.Println(err)
	}

	return &DBmetric{
		db: database,
	}
}

//MetricsCreate record metris from endpoints to db
func (database *DBmetric) MetricsCreate(w http.ResponseWriter, r *http.Request) {

	var metrics Metrics
	var Status string

	err := json.NewDecoder(r.Body).Decode(&metrics)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	Status = "Successful"

	vars := mux.Vars(r)
	metrics.ServiceName = vars[serviceName]

	time := time.Now()

	sqlStatement := `
	INSERT INTO metrics(servicename,metricvalue,metricname,status,time)
	VALUES($1,$2,$3,$4,$5) `

	_, err = database.db.Exec(sqlStatement, metrics.ServiceName, metrics.MetricValue, metrics.MetricName, Status, time)

	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)

}

//GetMetricsForService gets all metrics for the collection unit
func (database *DBmetric) GetMetricsForService(w http.ResponseWriter, r *http.Request) {

	var metrics Metrics

	vars := mux.Vars(r)
	metrics.ServiceName = vars[serviceName]

	rows, err := database.db.Query(postgresqlQuerry, metrics.ServiceName)

	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	ms := make([]*Metrics, 0)

	for rows.Next() {
		m := new(Metrics)
		err := rows.Scan(&m.ServiceName, &m.MetricValue, &m.MetricName, &m.Status, &m.Time)
		if err != nil {
			log.Fatal(err)
		}
		ms = append(ms, m)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, m := range ms {
		fmt.Fprintf(w, "%s %g %s %s %s \n", m.ServiceName, m.MetricValue, m.MetricName, m.Status, m.Time)
	}

	err = json.NewEncoder(w).Encode(metrics)
	if err != nil {
		log.Println(err)
	}
}
