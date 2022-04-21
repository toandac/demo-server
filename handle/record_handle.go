package handle

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"demo-server/models"
	"demo-server/repository"

	"github.com/go-chi/chi"
	"github.com/mssola/user_agent"
)

type RecordHandle struct {
	RecordRepo repository.RecordRepo
	URL        string
}

func (rc *RecordHandle) SaveRecord(w http.ResponseWriter, r *http.Request) {
	var req models.Request
	var events models.Events
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var record models.Record

	ua := user_agent.New(r.UserAgent())

	record.ID = req.SessionID
	events.Events = append(events.Events, req.Events...)

	record.User.Name = req.User.Name
	record.User.ID = req.UserID

	browserName, browserVersion := ua.Browser()
	record.Client = models.Client{
		ClientID:  req.ClientID,
		UserAgent: r.UserAgent(),
		OS:        ua.OS(),
		Browser:   browserName,
		Version:   browserVersion,
	}

	record.UpdatedAt = time.Now().Format("02/01/2006, 15:04")

	err := rc.RecordRepo.Insert(record, events)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *RecordHandle) RenderRecordScript(w http.ResponseWriter, r *http.Request) {
	tmplRecorder := template.Must(template.ParseFiles("templates/record.js"))

	err := tmplRecorder.Execute(w, struct{ URL string }{URL: rc.URL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *RecordHandle) RenderRecordPlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var events models.Events

	if err := rc.RecordRepo.QueryEventByID(id, &events); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tmplPlayerHTML := template.Must(template.ParseFiles("templates/session_by_id.html"))

	err := tmplPlayerHTML.Execute(w, struct {
		ID     string
		Events models.Events
	}{
		ID:     id,
		Events: events,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *RecordHandle) RendersRecordsList(w http.ResponseWriter, r *http.Request) {
	var record models.Record

	listID, err := rc.RecordRepo.QueryAllSessionID()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	records, err := rc.RecordRepo.QueryAllRecord(listID, record)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tmplListHTML := template.Must(template.ParseFiles("templates/session_list.html"))

	err = tmplListHTML.Execute(w, struct {
		Records []models.Record
		URL     string
	}{
		Records: records,
		URL:     rc.URL,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *RecordHandle) GetAllRecordByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var record models.Record

	if err := rc.RecordRepo.QueryRecordByID(id, &record); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&record); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rc *RecordHandle) GetAllEventByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var events models.Events

	if err := rc.RecordRepo.QueryEventByID(id, &events); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&events); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
