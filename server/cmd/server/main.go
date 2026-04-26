package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/do"
	dbUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/db_utils"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1/client_serverv1connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	dbLocation := flag.String("location", filepath.Join(xdg.StateHome, "forget-about-it.sqlite3"), "sqlite3 file location")
	flag.Parse()
	server.DBFilename = *dbLocation
	q, db := dbUtils.GetDB()
	realSub, err := fs.Sub(server.Files, "web")
	if err != nil {
		panic(err)
	}
	handler := http.FileServer(http.FS(realSub))
	http.DefaultServeMux.Handle("/", http.StripPrefix("/", handler))

	path, connectHandler := client_serverv1connect.NewForgetAboutItServiceHandler(&do.Server{Db: q, OrigDB: db})
	http.DefaultServeMux.Handle(path, connectHandler)

	address, err := q.GetConfigValue(context.Background(), "host")
	if err != nil {
		if err == sql.ErrNoRows {
			address = ":80"
		} else {
			panic(err)
		}
	}
	address = strings.TrimSpace(address)
	address = strings.TrimPrefix(address, "http://")
	address = strings.TrimPrefix(address, "https://")
	if address == "" {
		address = ":80"
	}
	fmt.Printf("Starting server on address %s\n", address)
	corsOptions := cors.Options{
		AllowPrivateNetwork: true,
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
	}
	log.Fatal(http.ListenAndServe(address, h2c.NewHandler(cors.New(corsOptions).Handler(http.DefaultServeMux), &http2.Server{})))

}
