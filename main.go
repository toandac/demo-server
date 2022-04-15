package main

import (
	"demo-server/database"
	"demo-server/handle"
	repoimpl "demo-server/repository/repo_impl"
	"demo-server/router"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/urfave/cli/v2"
)

func run(c *cli.Context) error {
	// influx := &database.InfluxDB{
	// 	URL:    "http://localhost:8086",
	// 	Token:  "162L5asvhPfFeE3hp4EWxT8Z6XXkNfzsh4XQ-R6XunRXrnYfJnd_AOlQ-dDyxcmC3OCQm829pbuWf_QNfgJOvA==",
	// 	Bucket: "events_bucket",
	// }
	// influx.NewInfluxDB()
	// defer influx.Close()

	dbDSN, err := url.Parse(c.String("db-dsn"))
	if err != nil {
		log.Println(err)
	}

	badger := &database.BadgerDB{}
	badger.NewBadgerDB(dbDSN)
	defer badger.Close()

	recordHandle := handle.RecordHandle{
		RecordRepo: repoimpl.NewRecordRepo(badger),
		URL:        c.String("service-url"),
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	}))
	api := router.API{
		Chi:          r,
		RecordHandle: recordHandle,
	}
	api.SetupRouter()

	return http.ListenAndServe(c.String("address")+":"+c.String("port"), r)
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "service-url", Value: "http://localhost:3000"},
			&cli.StringFlag{Name: "address", Value: "127.0.0.1"},
			&cli.StringFlag{Name: "db-dsn", Value: "badger:///tmp/badgerdb_1.3"},
			&cli.StringFlag{Name: "port", Value: "3000"},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}
