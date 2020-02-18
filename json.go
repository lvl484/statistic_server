package main

//Metrics ...
type Metrics struct {
	ServiceName string `json:"ServiceName"`
	MetricName  string `json:"MetricName"`
	MetricValue int    `json:"MetricValue"`
}
