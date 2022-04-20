package handle

import (
	"demo-server/models"
	"demo-server/repository"
	"encoding/json"
	"log"
	"net/http"
)

type EventHandle struct {
	EventRepo repository.EventRepo
	URL       string
}

func (e *EventHandle) SaveEvents(w http.ResponseWriter, r *http.Request) {
	var reqEvents models.Events
	if err := json.NewDecoder(r.Body).Decode(&reqEvents); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var events models.Events

	events.Events = append(events.Events, reqEvents.Events...)

	err := e.EventRepo.Insert(events)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}
