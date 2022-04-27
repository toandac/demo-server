package router

import (
	"analytics-api/handle"

	"github.com/gin-gonic/gin"
)

type API struct {
	Gin           *gin.Engine
	SessionHandle handle.SessionHandle
}

func (api *API) SetupRouter() {
	api.Gin.GET("/", api.SessionHandle.RenderListSession)
	api.Gin.POST("/session/save", api.SessionHandle.SaveSession)
	api.Gin.GET("/session/:session_id", api.SessionHandle.RenderSessionPlay)
	api.Gin.GET("/session/events/:session_id", api.SessionHandle.GetAllEventLimitByID)
	api.Gin.GET("/api/v1/session/:session_id", api.SessionHandle.GetAllSessionByID)
}
