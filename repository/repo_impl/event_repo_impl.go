package repoimpl

import (
	"context"
	"demo-server/database"
	"demo-server/models"
	"demo-server/repository"
	"fmt"
	"log"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type EventRepoImpl struct {
	influx *database.InfluxDB
}

func NewEventRepo(influx *database.InfluxDB) repository.EventRepo {
	return &EventRepoImpl{
		influx: influx,
	}
}

func (e *EventRepoImpl) Insert(record models.Record) error {
	for _, elem := range record.Events {
		p := influxdb2.NewPointWithMeasurement("test5").
			AddTag("clientID", record.Client.ClientID).
			AddTag("sessionID", record.ID).
			AddTag("userID", record.User.ID).
			AddTag("email", record.User.Email).
			AddTag("name", record.User.Name).
			AddTag("browser", record.Client.Browser).
			AddTag("os", record.Client.OS).
			AddTag("userAgent", record.Client.UserAgent).
			AddTag("version", record.Client.Version).
			AddField("data", elem.Data).
			AddField("type", elem.Type).
			AddField("timestamp", elem.Timestamp).
			SetTime(record.UpdatedAt)

		writeAPI := e.influx.Client.WriteAPIBlocking("tanda_organization", "records")
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			log.Println("Influxdb fails insert: ", err)
			return err
		}
	}

	return nil
}

func (e *EventRepoImpl) Query(id string) (models.Record, error) {
	var record models.Record
	var event models.Events

	queryAPI := e.influx.Client.QueryAPI("tanda_organization")

	query := fmt.Sprintf(`from(bucket: "records")
	|> range(start: -1d)
	|> filter(fn: (r) => r["_measurement"] == "test5")
	|> filter(fn: (r) => r["sessionID"] == "%s")`, id)

	// fmt.Println(query)

	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Access data

			values := result.Record().Values()
			record.ID = values["sessionID"].(string)
			// record.ClientID = values["clientID"].(string)

			record.Client.OS = values["os"].(string)
			record.Client.UserAgent = values["userAgent"].(string)
			record.Client.Version = values["version"].(string)

			// record.User.ID = values["userID"].(string)
			record.User.Email = values["email"].(string)
			record.User.Name = values["name"].(string)

			record.Client.Browser = values["browser"].(string)
			record.Client.Browser = values["browser"].(string)

			if strings.Contains(result.Record().Field(), "type") {
				event.Type = result.Record().Value().(int64)
			}

			if strings.Contains(result.Record().Field(), "data") {
				event.Data = result.Record().Value()
			}

			if strings.Contains(result.Record().Field(), "timestamp") {
				event.Timestamp = result.Record().Value().(int64)
			}

			record.Events = append(record.Events, event)

		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}

	} else {
		log.Println(err)
	}

	return record, nil
}
