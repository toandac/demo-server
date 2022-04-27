package handle

import (
	"encoding/json"
	"net/http"
	"time"

	"analytics-api/models"
	"analytics-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"go.uber.org/zap"
)

type SessionHandle struct {
	SessionRepo repository.SessionRepo
	Log         *zap.Logger
}

func (s *SessionHandle) SaveSession(c *gin.Context) {
	var request models.Request
	var session models.Session
	var events models.Events

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session.SessionID = request.SessionID
	events.Events = append(events.Events, request.Events...)

	session.User.UserName = request.User.UserName
	session.User.UserID = request.UserID

	ua := user_agent.New(c.Request.UserAgent())
	browserName, browserVersion := ua.Browser()
	session.Client = models.Client{
		ClientID:  request.ClientID,
		UserAgent: c.Request.UserAgent(),
		OS:        ua.OS(),
		Browser:   browserName,
		Version:   browserVersion,
	}

	session.UpdatedAt = time.Now().Format("02/01/2006, 15:04:05")

	err := s.SessionRepo.Insert(session, events)
	if err != nil {
		s.Log.Error("Error insert record ", zap.Error(err))
	}

	c.String(200, "ok")
}

func (s *SessionHandle) RenderSessionPlay(c *gin.Context) {
	sessionID := c.Param("session_id")

	c.HTML(200, "session_by_id.html", gin.H{
		"SessionID": sessionID,
	})
}

func (s *SessionHandle) RenderListSession(c *gin.Context) {
	var session models.Session

	listID, err := s.SessionRepo.GetAllSessionID()
	if err != nil {
		return
	}

	sessions, err := s.SessionRepo.GetAllSession(listID, session)
	if err != nil {
		return
	}

	c.HTML(200, "session_list.html", gin.H{
		"Sessions": sessions,
		"URL":      "http://localhost:3000",
	})
}

func (s *SessionHandle) GetAllSessionByID(c *gin.Context) {
	sessionID := c.Param("session_id")
	var session models.Session

	if err := s.SessionRepo.GetSessionByID(sessionID, &session); err != nil {
		return
	}

	c.JSON(200, session)
}

func (s *SessionHandle) GetAllEventLimitByID(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(200)

	sessionID := c.Param("session_id")
	limit := 2
	offset := 0

	totalRecord, err := s.SessionRepo.GetTotalColumn(sessionID)
	if err != nil {
		return
	}
	s.Log.Sugar().Info("Total number rows ", totalRecord)

	msgChan := make(chan []models.Event)
	defer func() {
		close(msgChan)
		msgChan = nil
		s.Log.Info("Client connection is closed")
	}()

	go func() {
		if msgChan != nil {
			for offset <= int(totalRecord) {
				var events models.Events
				if err := s.SessionRepo.GetEventLimitBySessionID(sessionID, limit, offset, &events); err != nil {
					s.Log.Error("Error ", zap.Error(err))
					return
				}
				offset = offset + limit

				msgChan <- events.Events
			}
		}
	}()

	for {
		select {
		case message := <-msgChan:
			enc := json.NewEncoder(c.Writer)
			if err := enc.Encode(message); err != nil {
				return
			}
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
