package repository

import "demo-server/models"

type EventRepo interface {
	Insert(record models.Record) error
	Query(id string) (models.Record, error)
}
