package repoimpl

import (
	"demo-server/database"
	"demo-server/models"
	"demo-server/repository"

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

func (r *RecordRepoImpl) Insert(record models.Record) error {
	p := influxdb2.NewPointWithMeasurement(r.influx.Measurement).
		AddTag("clientID", record.Client.ClientID).
		AddTag("sessionID", record.ID).
		AddTag("userID", record.User.ID).
		AddTag("email", record.User.Email).
		AddTag("name", record.User.Name).
		AddTag("browser", record.Client.Browser).
		AddTag("os", record.Client.OS).
		AddTag("userAgent", record.Client.UserAgent).
		AddTag("version", record.Client.Version).
		AddTag("updatedAt", record.UpdatedAt).
		SetTime(time.Now())

	writeAPI := r.influx.Client.WriteAPIBlocking(r.influx.Organization, r.influx.Bucket)
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Println("Influxdb fails insert: ", err)
		return err
	}

	return nil
}

// func (r *RecordRepoImpl) QueryRecordByID(id string, record *models.Record) error {
// 	var event models.Events
// 	data := event.Data

// 	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

// 	query := fmt.Sprintf(`from(bucket: "%s")
// 	|> range(start: -1d)
// 	|> filter(fn: (r) => r["_measurement"] == "%s")
// 	|> filter(fn: (r) => r["sessionID"] == "%s")`, r.influx.Bucket, r.influx.Measurement, id)

// 	result, err := queryAPI.Query(context.Background(), query)
// 	if err == nil {
// 		for result.Next() {
// 			values := result.Record().Values()

// 			record.ID = values["sessionID"].(string)
// 			// record.ClientID = values["clientID"].(string)

// 			record.Client.OS = values["os"].(string)
// 			record.Client.UserAgent = values["userAgent"].(string)
// 			record.Client.Version = values["version"].(string)
// 			record.Client.Browser = values["browser"].(string)

// 			// record.User.ID = values["userID"].(string)
// 			record.User.Email = values["email"].(string)
// 			record.User.Name = values["name"].(string)

// 			record.UpdatedAt = values["updatedAt"].(string)

// 			timestampString := values["timestamp"].(string)
// 			timestamp, err := common.StringToInt64(timestampString)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			event.Timestamp = timestamp

// 			typeString := values["type"].(string)
// 			typeEvent, err := common.StringToInt64(typeString)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			event.Type = typeEvent

// 			json.Unmarshal([]byte(result.Record().Value().(string)), &data)
// 			event.Data = data

// 			record.Events = append(record.Events, event)
// 		}
// 		if result.Err() != nil {
// 			fmt.Printf("query parsing error: %s\n", result.Err().Error())
// 		}

// 	} else {
// 		log.Println(err)
// 	}

// 	return nil
// }

// func (r *RecordRepoImpl) QueryAllSessionID() ([]string, error) {
// 	var listID []string

// 	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

// 	query := fmt.Sprintf(`from(bucket: "%s")
// 	|> range(start: -1d)
// 	|> filter(fn: (r) => r["_measurement"] == "%s")
// 	|> group(columns: ["_measurement"])`, r.influx.Bucket, r.influx.Measurement)

// 	result, err := queryAPI.Query(context.Background(), query)
// 	if err == nil {
// 		for result.Next() {
// 			listID = append(listID, result.Record().ValueByKey("sessionID").(string))
// 		}

// 		listID = common.RemoveDuplicateValues(listID)

// 		if result.Err() != nil {
// 			fmt.Printf("query parsing error: %s\n", result.Err().Error())
// 		}

// 	} else {
// 		log.Println("err", err)
// 	}

// 	return listID, nil
// }

// func (r *RecordRepoImpl) QueryAllRecord(listID []string, record models.Record) ([]models.Record, error) {
// 	var records []models.Record
// 	var event models.Events

// 	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

// 	for _, id := range listID {
// 		query := fmt.Sprintf(`from(bucket: "%s")
// 		|> range(start: -1d)
// 		|> filter(fn: (r) => r["_measurement"] == "%s")
// 		|> filter(fn: (r) => r["sessionID"] == "%s")`, r.influx.Bucket, r.influx.Measurement, id)

// 		result, err := queryAPI.Query(context.Background(), query)
// 		if err == nil {
// 			for result.Next() {
// 				values := result.Record().Values()

// 				record.ID = values["sessionID"].(string)
// 				record.User.Name = values["name"].(string)
// 				record.UpdatedAt = values["updatedAt"].(string)

// 				record.Events = append(record.Events, event)

// 			}

// 			if result.Err() != nil {
// 				fmt.Printf("query parsing error: %s\n", result.Err().Error())
// 			}

// 		} else {
// 			log.Println("err", err)
// 		}

// 		records = append(records, record)
// 	}

// 	return records, nil
// }

// func (r *RecordRepoImpl) QueryEventDataByID(id string, record *models.Record) error {
// 	var event models.Events
// 	data := event.Data

// 	queryAPI := r.influx.Client.QueryAPI(r.influx.Organization)

// 	query := fmt.Sprintf(`from(bucket: "%s")
// 	|> range(start: -1d)
// 	|> filter(fn: (r) => r["_measurement"] == "%s")
// 	|> filter(fn: (r) => r["sessionID"] == "%s")`, r.influx.Bucket, r.influx.Measurement, id)

// 	result, err := queryAPI.Query(context.Background(), query)
// 	if err == nil {
// 		for result.Next() {
// 			json.Unmarshal([]byte(result.Record().Value().(string)), &data)
// 			event.Data = data

// 			record.Events = append(record.Events, event)
// 		}
// 		if result.Err() != nil {
// 			fmt.Printf("query parsing error: %s\n", result.Err().Error())
// 		}

// 	} else {
// 		log.Println(err)
// 	}

// 	return nil
// }
