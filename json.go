package main

//Metrics ...
type Metrics struct {
	MetricName  string  `json:"MetricName"`
	MetricValue float64 `json:"MetricValue"`
	ServiceName string  `json:"ServiceName"`
}
