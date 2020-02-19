package main

//Metrics ...
type Metrics struct {
	ServiceName string  `db:"servicename" json:"ServiceName"`
	MetricValue float64 `db:"metricvalue" json:"MetricValue"`
	MetricName  string  `db:"metricname" json:"MetricName"`
}
