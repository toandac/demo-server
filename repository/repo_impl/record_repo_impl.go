package repoimpl

import (
	"demo-server/common"
	"demo-server/database"
	"demo-server/models"
	"demo-server/repository"
	"encoding/json"
	"fmt"
	"strconv"

	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type RecordRepoImpl struct {
	influx *database.InfluxDB
}

func NewRecordRepo(influx *database.InfluxDB) repository.RecordRepo {
	return &RecordRepoImpl{
		influx: influx,
	}
}

func (r *RecordRepoImpl) Insert(record models.Record, events models.Events) error {
	for _, event := range events.Events {
		timestampStr := strconv.FormatInt(event.Timestamp, 10)
		typeStr := strconv.FormatInt(event.Type, 10)

		mJson, errMarshal := json.Marshal(event.Data)
		if errMarshal != nil {
			log.Println(errMarshal)
		}
		jsonStr := string(mJson)

		p := influxdb2.NewPointWithMeasurement(r.influx.Measurement).
			AddTag("clientID", record.Client.ClientID).
			AddTag("sessionID", record.ID).
			AddTag("userID", record.User.ID).
			AddTag("name", record.User.Name).
			AddTag("browser", record.Client.Browser).
			AddTag("os", record.Client.OS).
			AddTag("userAgent", record.Client.UserAgent).
			AddTag("version", record.Client.Version).
			AddField("data", jsonStr).
			AddTag("type", typeStr).
			AddTag("timestamp", timestampStr).
			AddTag("updatedAt", record.UpdatedAt).
			SetTime(time.Now())

		writeAPI := r.influx.Client.WriteAPIBlocking(r.influx.Organization, r.influx.Bucket)
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			log.Println("Influxdb fails insert record: ", err)
			return err
		}
	}

	return nil
}

func (r *RecordRepoImpl) QueryRecordByID(id string, record *models.Record) error {
	var event models.Event
	data := event.Data

	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -1d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> filter(fn: (r) => r["sessionID"] == "%s")`, r.influx.Bucket, r.influx.Measurement, id)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			values := result.Record().Values()

			record.ID = values["sessionID"].(string)
			record.Client.ClientID = values["clientID"].(string)

			record.Client.OS = values["os"].(string)
			record.Client.UserAgent = values["userAgent"].(string)
			record.Client.Version = values["version"].(string)
			record.Client.Browser = values["browser"].(string)

			record.User.ID = values["userID"].(string)
			record.User.Name = values["name"].(string)

			record.UpdatedAt = values["updatedAt"].(string)

			timestampString := values["timestamp"].(string)
			timestamp, err := common.StringToInt64(timestampString)
			if err != nil {
				log.Println(err)
			}
			event.Timestamp = timestamp

			typeString := values["type"].(string)
			typeEvent, err := common.StringToInt64(typeString)
			if err != nil {
				log.Println(err)
			}
			event.Type = typeEvent

			event.ID = values["sessionID"].(string)

			json.Unmarshal([]byte(result.Record().Value().(string)), &data)
			event.Data = data

			record.Events = append(record.Events, event)
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}

	} else {
		log.Println(err)
	}

	return nil
}

func (r *RecordRepoImpl) QueryAllSessionID() ([]string, error) {
	var listID []string

	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -1d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> group(columns: ["_measurement"])`, r.influx.Bucket, r.influx.Measurement)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			listID = append(listID, result.Record().ValueByKey("sessionID").(string))
		}

		listID = common.RemoveDuplicateValues(listID)

		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}

	} else {
		log.Println("err", err)
	}

	return listID, nil
}

func (r *RecordRepoImpl) QueryAllRecord(listID []string, record models.Record) ([]models.Record, error) {
	var records []models.Record
	var events models.Events
	var event models.Event

	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

	for _, id := range listID {
		query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -1d)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["sessionID"] == "%s")`, r.influx.Bucket, r.influx.Measurement, id)

		result, err := queryAPI.Query(context.Background(), query)
		if err == nil {
			for result.Next() {
				values := result.Record().Values()

				record.ID = values["sessionID"].(string)
				record.User.Name = values["name"].(string)
				record.UpdatedAt = values["updatedAt"].(string)

				events.Events = append(events.Events, event)

			}

			if result.Err() != nil {
				fmt.Printf("query parsing error: %s\n", result.Err().Error())
			}

		} else {
			log.Println("err", err)
		}

		records = append(records, record)
	}

	return records, nil
}

func (r *RecordRepoImpl) QueryEventByID(id string, events *models.Events) error {
	var event models.Event
	data := event.Data

	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -1d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> filter(fn: (r) => r["sessionID"] == "%s")`, r.influx.Bucket, r.influx.Measurement, id)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			values := result.Record().Values()

			timestampString := values["timestamp"].(string)
			timestamp, err := common.StringToInt64(timestampString)
			if err != nil {
				log.Println(err)
			}
			event.Timestamp = timestamp

			typeString := values["type"].(string)
			typeEvent, err := common.StringToInt64(typeString)
			if err != nil {
				log.Println(err)
			}
			event.Type = typeEvent

			event.ID = values["sessionID"].(string)

			json.Unmarshal([]byte(result.Record().Value().(string)), &data)
			event.Data = data

			events.Events = append(events.Events, event)
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}

	} else {
		log.Println(err)
	}

	return nil
}
