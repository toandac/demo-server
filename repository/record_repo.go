package repository

import "demo-server/models"

type RecordRepo interface {
	Get(id string, record *models.Record) error
	Put(record models.Record) error
	Iter(record models.Record) ([]models.Record, error)
}
