package main

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"database/sql"
	"encoding/json"
    "errors"
    "io"
    "strings"
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
	if r.Header.Get("Content-Type") != "" {
        value, _ := ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            msg := "Content-Type header is not application/json"
            http.Error(w, msg, http.StatusUnsupportedMediaType)
            return
        }
    }

    r.Body = http.MaxBytesReader(w, r.Body, 1048576)

    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    err := dec.Decode(&metrics)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError

        switch {
        case errors.As(err, &syntaxError):
            msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
            http.Error(w, msg, http.StatusBadRequest)

        case errors.Is(err, io.ErrUnexpectedEOF):
            msg := fmt.Sprintf("Request body contains badly-formed JSON")
            http.Error(w, msg, http.StatusBadRequest)

        case errors.As(err, &unmarshalTypeError):
            msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
            http.Error(w, msg, http.StatusBadRequest)

        case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
            http.Error(w, msg, http.StatusBadRequest)

        case errors.Is(err, io.EOF):
            msg := "Request body must not be empty"
            http.Error(w, msg, http.StatusBadRequest)

        case err.Error() == "http: request body too large":
            msg := "Request body must not be larger than 1MB"
            http.Error(w, msg, http.StatusRequestEntityTooLarge)

        default:
            log.Println(err.Error())
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return
    }

    if dec.More() {
        msg := "Request body must only contain a single JSON object"
        http.Error(w, msg, http.StatusBadRequest)
        return
    }

	_, err = database.db.Exec(
		"INSERT INTO metrics(ServiceName,MetricName,MetricValue) VALUES(?,?,?)",
		metrics.ServiceName, metrics.MetricName, metrics.MetricValue)
	if err != nil {
		log.Println(err)
		return
	}
}
