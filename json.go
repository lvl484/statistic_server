package main

//Metrics ...
type Metrics struct {
	ServiceName string  `db:"ServiceName" json:"ServiceName"`
	MetricValue float64 `db:"MetricValue" json:"MetricValue"`
	MetricName  string  `db:"MetricName" json:"MetricName"`
}
