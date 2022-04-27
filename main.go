package main

import (
	"analytics-api/database"
	"analytics-api/handle"
	"analytics-api/middlewares"
	repoimpl "analytics-api/repository/repo_impl"
	"analytics-api/router"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		return
	}
}

func main() {
	log, _ := handle.NewLog()

	influx := &database.InfluxDB{
		URL:          os.Getenv("INFLUX_URL"),
		Token:        os.Getenv("INFLUX_TOKEN"),
		Bucket:       os.Getenv("INFLUX_BUCKET"),
		Measurement:  os.Getenv("INFLUX_MEASUREMENT"),
		Organization: os.Getenv("INFLUX_ORGANIZATION"),

		Logger: log,
	}
	influx.NewInfluxDB()
	defer influx.Close()

	sessiondHandle := handle.SessionHandle{
		SessionRepo: repoimpl.NewSessionRepo(influx, log),
		Log:         log,
	}

	g := gin.Default()

	g.LoadHTMLGlob("templates/*")
	g.StaticFile("/record.js", "./js/record.js")

	g.Use(middlewares.CORSMiddleware())

	api := router.API{
		Gin:           g,
		SessionHandle: sessiondHandle,
	}
	api.SetupRouter()

	g.Run(":" + os.Getenv("SERVER_PORT"))
}
