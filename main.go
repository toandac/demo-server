package main

import (
	"log"
	"net/http"
	"os"

	"demo-server/database"
	"demo-server/handle"
	repoimpl "demo-server/repository/repo_impl"
	"demo-server/router"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/urfave/cli/v2"
)

func run(c *cli.Context) error {
	influx := &database.InfluxDB{
		URL:          c.String("influx-url"),
		Token:        c.String("influx-token"),
		Bucket:       c.String("influx-bucket"),
		Measurement:  c.String("influx-measurement"),
		Organization: c.String("influx-organization"),
	}
	influx.NewInfluxDB()
	defer influx.Close()

	recordHandle := handle.RecordHandle{
		RecordRepo: repoimpl.NewRecordRepo(influx),
		URL:        c.String("service-url"),
	}

	eventHandle := handle.EventHandle{
		EventRepo: repoimpl.NewEventRepo(influx),
		URL:       c.String("service-url"),
	}

	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type"},
		}))
	api := router.API{
		Chi:          r,
		RecordHandle: recordHandle,
		EventHandle:  eventHandle,
	}
	api.SetupRouter()

	return http.ListenAndServe(c.String("address")+":"+c.String("port"), r)
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "service-url", Value: "http://localhost:3000"},
			&cli.StringFlag{Name: "address", Value: "127.0.0.1"},
			&cli.StringFlag{Name: "port", Value: "3000"},

			&cli.StringFlag{Name: "influx-url", Value: "http://localhost:8086"},
			&cli.StringFlag{Name: "influx-token", Value: "162L5asvhPfFeE3hp4EWxT8Z6XXkNfzsh4XQ-R6XunRXrnYfJnd_AOlQ-dDyxcmC3OCQm829pbuWf_QNfgJOvA=="},
			&cli.StringFlag{Name: "influx-bucket", Value: "records"},
			&cli.StringFlag{Name: "influx-measurement", Value: "test"},
			&cli.StringFlag{Name: "influx-organization", Value: "tanda_organization"},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}
