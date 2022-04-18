package handle

import (
	"demo-server/models"
	"demo-server/repository"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/mssola/user_agent"
)

type RecordHandle struct {
	// RecordRepo repository.RecordRepo
	EventRepo repository.EventRepo
	URL       string
}

func (rc *RecordHandle) SaveRecord(w http.ResponseWriter, r *http.Request) {
	var req models.Record
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rec models.Record
	// if err := rc.RecordRepo.Get(req.ID, &rec); err != nil {
	// 	if !errors.Is(err, bow.ErrNotFound) {
	// 		log.Println(err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	ua := user_agent.New(r.UserAgent())

	rec.ID = req.ID
	rec.Events = append(rec.Events, req.Events...)
	rec.User = req.User
	rec.UpdatedAt = time.Now()

	browserName, browserVersion := ua.Browser()
	rec.Client = models.Client{
		UserAgent: r.UserAgent(),
		OS:        ua.OS(),
		Browser:   browserName,
		Version:   browserVersion,
	}

	// if err := rc.RecordRepo.Put(rec); err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	err := rc.EventRepo.Insert(rec)
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

	// var record models.Record

	// if err := rc.RecordRepo.Get(id, &record); err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }

	record, err := rc.EventRepo.Query(id)
	if err != nil {
		log.Println(err)
	}

	tmplPlayerHTML := template.Must(template.ParseFiles("templates/session_by_id.html"))

	err = tmplPlayerHTML.Execute(w, struct {
		ID     string
		Record models.Record
	}{
		ID:     id,
		Record: record,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// func (rc *RecordHandle) RendersRecordsList(w http.ResponseWriter, r *http.Request) {
// 	var record models.Record

// 	records, err := rc.RecordRepo.Iter(record)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	tmplListHTML := template.Must(template.ParseFiles("templates/session_list.html"))

// 	err = tmplListHTML.Execute(w, struct {
// 		Records []models.Record
// 		URL     string
// 	}{
// 		Records: records,
// 		URL:     rc.URL,
// 	})
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

func (rc *RecordHandle) GetAllRecordByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// var record models.Record

	// if err := rc.RecordRepo.Get(id, &record); err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusNotFound)
	// 	return
	// }

	record, err := rc.EventRepo.Query(id)
	if err != nil {
		log.Println(err)
	}

	if err := json.NewEncoder(w).Encode(&record); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
