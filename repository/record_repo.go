package repository

import "demo-server/models"

type RecordRepo interface {
	Insert(record models.Record) error
	// QueryRecordByID(id string, record *models.Record) error
	// QueryEventDataByID(id string, record *models.Record) error
	// QueryAllRecord(listID []string, record models.Record) ([]models.Record, error)
	// QueryAllSessionID() ([]string, error)
}
