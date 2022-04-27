package repository

import "analytics-api/models"

type SessionRepo interface {
	Insert(session models.Session, events models.Events) error
	GetTotalColumn(sessionID string) (int64, error)
	GetEventLimitBySessionID(sessionID string, limit, offset int, events *models.Events) error
	GetSessionByID(sessionID string, session *models.Session) error
	GetAllSession(listID []string, session models.Session) ([]models.Session, error)
	GetAllSessionID() ([]string, error)
}
