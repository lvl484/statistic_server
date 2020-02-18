package main

import (
	"database/sql"
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
	//var metrics Metrics

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "InvalidMethod")
		return
	}

	r.ParseForm()

	//
	log.Println(r.Form)

	for key, values := range r.Form {
		for _, value := range values {
			fmt.Println(key, value)
		}
	}
	//

	/* fmt.Fprintf(w, r.Form.Get("name"))
	}

		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters  in the query.
		_, err = database.db.Exec(
			"INSERT INTO users(ServiceName,MetricName,MetricValue) VALUES(?,?,?,?)",
	        metrics.ServiceName, metrics.MetricName, metrics.MetricValue)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
	*/
}
