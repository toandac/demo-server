package repository

import "demo-server/models"

type EventRepo interface {
	Insert(events models.Events) error
}
