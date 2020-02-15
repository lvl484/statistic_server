package main

//Metrics - my form for saving metrics
type Metrics struct {
	ServiceName int `json:"ServiceName"`
	MetricName  int `json:"MetricName"`
	MetricValue int `json:"MetricValue"`
}
