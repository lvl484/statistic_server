package main

//Metrics ...
type Metrics struct {
	ServiceName string  `json:"ServiceName"`
	MetricValue float64 `json:"MetricValue"`
	MetricName  string  `json:"MetricName"`
}
