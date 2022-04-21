package repository

import "demo-server/models"

type RecordRepo interface {
	Insert(record models.Record, events models.Events) error
	QueryRecordByID(id string, record *models.Record) error
	QueryEventByID(id string, events *models.Events) error
	QueryAllRecord(listID []string, record models.Record) ([]models.Record, error)
	QueryAllSessionID() ([]string, error)
}
