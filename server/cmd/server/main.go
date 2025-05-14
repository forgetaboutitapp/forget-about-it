package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/do"

	"github.com/adrg/xdg"
	"github.com/forgetaboutitapp/forget-about-it/server"
	dbUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/db_utils"
	"github.com/rs/cors"
)

func main() {
	dbLocation := flag.String("location", filepath.Join(xdg.StateHome, "forget-about-it.sqlite3"), "sqlite3 file location")
	port := flag.Int("port", 8080, "port to host")
	flag.Parse()
	server.DBFilename = *dbLocation
	q, db := dbUtils.GetDB()
	realSub, err := fs.Sub(server.Files, "web")
	if err != nil {
		panic(err)
	}
	handler := http.FileServer(http.FS(realSub))
	http.DefaultServeMux.Handle("/", http.StripPrefix("/", handler))
	http.DefaultServeMux.Handle("/api/v1/do", do.Server{Db: q, OrigDB: db})

	address := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting server on address %s\n", address)
	corsOptions := cors.Options{
		AllowPrivateNetwork: true,
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
	}
	log.Fatal(http.ListenAndServe(address, cors.New(corsOptions).Handler(http.DefaultServeMux)))

}
