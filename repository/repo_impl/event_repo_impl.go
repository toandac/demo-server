package repoimpl

import (
	"context"
	"demo-server/database"
	"demo-server/models"
	"demo-server/repository"
	"encoding/json"
	"log"
	"strconv"
	"time"

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

func (e *EventRepoImpl) Insert(events models.Events) error {
	for _, event := range events.Events {
		timestampStr := strconv.FormatInt(event.Timestamp, 10)
		typeStr := strconv.FormatInt(event.Type, 10)

		mJson, errMarshal := json.Marshal(event.Data)
		if errMarshal != nil {
			log.Println(errMarshal)
		}
		jsonStr := string(mJson)

		p := influxdb2.NewPointWithMeasurement(e.influx.Measurement).
			AddTag("sessionID", event.SessionID).
			AddField("data", jsonStr).
			AddTag("type", typeStr).
			AddTag("timestamp", timestampStr).
			SetTime(time.Now())

		writeAPI := e.influx.Client.WriteAPIBlocking(e.influx.Organization, e.influx.Bucket)
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			log.Println("Influxdb fails insert: ", err)
			return err
		}
	}
	return nil
}
