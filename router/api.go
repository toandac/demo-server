package router

import (
	"demo-server/handle"

	"github.com/go-chi/chi"
)

type API struct {
	Chi          *chi.Mux
	RecordHandle handle.RecordHandle
}

func (api *API) SetupRouter() {
	api.Chi.Get("/", api.RecordHandle.RendersRecordsList)
	api.Chi.Post("/sessions", api.RecordHandle.SaveRecord)
	api.Chi.Get("/sessions/{id}", api.RecordHandle.RenderRecordPlayer)
	api.Chi.Get("/api/v1/sessions/{id}", api.RecordHandle.GetAllRecordByID)
	api.Chi.Get("/record.js", api.RecordHandle.RenderRecordScript)
}
