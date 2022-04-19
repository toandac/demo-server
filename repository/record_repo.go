package repository

import "demo-server/models"

type RecordRepo interface {
	Insert(record models.Record) error
	Query(id string, record *models.Record) error
	QueryAll(listID []string, record models.Record) ([]models.Record, error)
	QueryAllSessionID() ([]string, error)
}
