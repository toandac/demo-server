package database

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDB struct {
	Client       influxdb2.Client
	URL          string
	Token        string
	Bucket       string
	Measurement  string
	Organization string
}

func (influx *InfluxDB) NewInfluxDB() {
	influx.Client = influxdb2.NewClientWithOptions(influx.URL, influx.Token, influxdb2.DefaultOptions().SetHTTPRequestTimeout(50))
}

func (influx *InfluxDB) Close() {
	influx.Client.Close()
}
