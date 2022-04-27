package repoimpl

import (
	"analytics-api/database"
	"analytics-api/models"
	"analytics-api/repository"

	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.uber.org/zap"
)

type SessionRepoImpl struct {
	influx *database.InfluxDB
	log    *zap.Logger
}

func NewSessionRepo(influx *database.InfluxDB, log *zap.Logger) repository.SessionRepo {
	return &SessionRepoImpl{
		influx: influx,
		log:    log,
	}
}

func (s *SessionRepoImpl) Insert(session models.Session, events models.Events) error {
	for _, event := range events.Events {
		timestampStr := strconv.FormatInt(event.Timestamp, 10)
		typeStr := strconv.FormatInt(event.Type, 10)

		mJson, errMarshal := json.Marshal(event.Data)
		if errMarshal != nil {
			s.log.Error("Error marshal event data ", zap.Error(errMarshal))
			return errMarshal
		}
		jsonStr := string(mJson)

		p := influxdb2.NewPointWithMeasurement(s.influx.Measurement).
			AddTag("clientID", session.Client.ClientID).
			AddTag("sessionID", session.SessionID).
			AddTag("userID", session.User.UserID).
			AddTag("userName", session.User.UserName).
			AddTag("browser", session.Client.Browser).
			AddTag("os", session.Client.OS).
			AddTag("userAgent", session.Client.UserAgent).
			AddTag("version", session.Client.Version).
			AddField("data", jsonStr).
			AddTag("type", typeStr).
			AddTag("timestamp", timestampStr).
			AddTag("updatedAt", session.UpdatedAt).
			SetTime(time.Now())

		writeAPI := s.influx.Client.WriteAPIBlocking(s.influx.Organization, s.influx.Bucket)
		err := writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			s.log.Error("Influxdb fails insert record ", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *SessionRepoImpl) GetSessionByID(sessionID string, session *models.Session) error {
	var event models.Event
	data := event.Data

	queryAPI := s.influx.Client.QueryAPI(s.influx.Organization)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -7d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> filter(fn: (r) => r["sessionID"] == "%s")`, s.influx.Bucket, s.influx.Measurement, sessionID)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			values := result.Record().Values()

			session.SessionID = values["sessionID"].(string)
			session.Client.ClientID = values["clientID"].(string)

			session.Client.OS = values["os"].(string)
			session.Client.UserAgent = values["userAgent"].(string)
			session.Client.Version = values["version"].(string)
			session.Client.Browser = values["browser"].(string)

			session.User.UserID = values["userID"].(string)
			session.User.UserName = values["userName"].(string)

			session.UpdatedAt = values["updatedAt"].(string)

			timestampString := values["timestamp"].(string)
			timestamp, err := stringToInt64(timestampString)
			if err != nil {
				s.log.Error("Error convert type timestamp from string to int64 ", zap.Error(err))
				return err
			}
			event.Timestamp = timestamp

			typeString := values["type"].(string)
			typeEvent, err := stringToInt64(typeString)
			if err != nil {
				s.log.Error("Error convert type event from string to int64 ", zap.Error(err))
				return err
			}
			event.Type = typeEvent

			json.Unmarshal([]byte(result.Record().Value().(string)), &data)
			event.Data = data

			session.Events = append(session.Events, event)
		}
		if result.Err() != nil {
			s.log.Error("Query parsing error ", zap.Error(result.Err()))
			return result.Err()
		}

	} else {
		s.log.Error("Query parsing error ", zap.Error(result.Err()))
		return err
	}

	return nil
}

func (s *SessionRepoImpl) GetAllSessionID() ([]string, error) {
	var listID []string

	queryAPI := s.influx.Client.QueryAPI(s.influx.Organization)

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -7d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> group(columns: ["_measurement"])`, s.influx.Bucket, s.influx.Measurement)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			listID = append(listID, result.Record().ValueByKey("sessionID").(string))
		}

		listID = removeDuplicateValues(listID)

		if result.Err() != nil {
			s.log.Error("Query parsing error ", zap.Error(result.Err()))
			return nil, result.Err()
		}

	} else {
		s.log.Error("Query parsing error ", zap.Error(result.Err()))
		return nil, err
	}

	return listID, nil
}

func (s *SessionRepoImpl) GetAllSession(listID []string, session models.Session) ([]models.Session, error) {
	var sessions []models.Session
	var events models.Events
	var event models.Event

	queryAPI := s.influx.Client.QueryAPI(s.influx.Organization)

	for _, id := range listID {
		query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -7d)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["sessionID"] == "%s")`, s.influx.Bucket, s.influx.Measurement, id)

		result, err := queryAPI.Query(context.Background(), query)
		if err == nil {
			for result.Next() {
				values := result.Record().Values()

				session.SessionID = values["sessionID"].(string)
				session.User.UserName = values["userName"].(string)
				session.UpdatedAt = values["updatedAt"].(string)

				events.Events = append(events.Events, event)

			}

			if result.Err() != nil {
				s.log.Error("Query parsing error ", zap.Error(result.Err()))
				return nil, result.Err()
			}

		} else {
			s.log.Error("Query parsing error ", zap.Error(result.Err()))
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *SessionRepoImpl) GetEventLimitBySessionID(sessionID string, limit, offset int, events *models.Events) error {
	var event models.Event
	queryAPI := s.influx.Client.QueryAPI(s.influx.Organization)

	fmt.Printf("\n")
	s.log.Sugar().Info("limit ", limit, " offset ", offset)

	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -7d)
		|> filter(fn: (r) => r["_measurement"] == "%s")
		|> filter(fn: (r) => r["sessionID"] == "%s")
		|> group()
		|> limit(n:%d, offset: %d)`, s.influx.Bucket, s.influx.Measurement, sessionID, limit, offset)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			values := result.Record().Values()

			timestampString := values["timestamp"].(string)
			timestamp, err := stringToInt64(timestampString)
			if err != nil {
				s.log.Error("Error convert type timestamp from string to int64 ", zap.Error(err))
				return err
			}
			event.Timestamp = timestamp

			typeString := values["type"].(string)
			typeEvent, err := stringToInt64(typeString)
			if err != nil {
				s.log.Error("Error convert type event from string to int64 ", zap.Error(err))
				return err
			}
			event.Type = typeEvent

			json.Unmarshal([]byte(result.Record().Value().(string)), &event.Data)

			if len(events.Events) <= limit {
				s.log.Sugar().Info("event timestamp ", event.Timestamp)
				events.Events = append(events.Events, event)
			}
		}
		if result.Err() != nil {
			s.log.Error("Query parsing error ", zap.Error(result.Err()))
			return result.Err()
		}

	} else {
		s.log.Error("Query parsing error ", zap.Error(result.Err()))
		return err
	}
	s.log.Sugar().Info("Number event ", len(events.Events))

	return nil
}

func (s *SessionRepoImpl) GetTotalColumn(sessionID string) (int64, error) {
	queryAPI := s.influx.Client.QueryAPI(s.influx.Organization)
	var numberColumn int64

	query := fmt.Sprintf(`from(bucket: "%s")
	|> range(start: -7d)
	|> filter(fn: (r) => r["_measurement"] == "%s")
	|> group(columns: ["_measurement"])
	|> filter(fn: (r) => r["sessionID"] == "%s")
	|> group()
    |> count(column: "sessionID")`, s.influx.Bucket, s.influx.Measurement, sessionID)

	result, err := queryAPI.Query(context.Background(), query)
	if err == nil {
		for result.Next() {
			numberColumn = result.Record().Values()["sessionID"].(int64)
		}
		if result.Err() != nil {
			s.log.Error("Query parsing error ", zap.Error(result.Err()))
			return 0, err
		}
	} else {
		s.log.Error("Query parsing error ", zap.Error(result.Err()))
		return 0, err
	}

	return numberColumn, nil
}

func stringToInt64(str string) (int64, error) {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func removeDuplicateValues(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
