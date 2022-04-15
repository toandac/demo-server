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
	api.Chi.Post("/record", api.RecordHandle.SaveRecord)
	api.Chi.Get("/record.js", api.RecordHandle.RenderRecordScript)
	api.Chi.Get("/records/{id}", api.RecordHandle.RenderRecordPlayer)
	api.Chi.Get("/", api.RecordHandle.RendersRecordsList)
	api.Chi.Get("/api/v1/records/{id}", api.RecordHandle.GetAllRecordByID)
}
