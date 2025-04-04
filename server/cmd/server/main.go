package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/login"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
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
	http.DefaultServeMux.Handle("/api/v0/get-token", login.Server{Db: q})
	http.DefaultServeMux.Handle("/api/v0/secure/generate-new-token", secure.Server{Db: q, Next: secure.GenerateNewToken, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/check-new-token", secure.Server{Db: q, Next: secure.CheckNewToken, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/delete-new-token", secure.Server{Db: q, Next: secure.DeleteNewToken, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/get-remote-settings", secure.Server{Db: q, Next: secure.GetRemoteSettings, OrigDB: db})

	http.DefaultServeMux.Handle("/api/v0/secure/get-all-questions", secure.Server{Db: q, Next: secure.GetAllQuestions, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/post-all-questions", secure.Server{Db: q, Next: secure.PostAllQuestions, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/get-all-tags", secure.Server{Db: q, Next: secure.GetAllTags, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/grade-question", secure.Server{Db: q, Next: secure.GradeQuestion, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/get-next-question", secure.Server{Db: q, Next: secure.GetNextQuestion, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/upload-algorithm", secure.Server{Db: q, Next: secure.UploadAlgorithm, OrigDB: db})
	http.DefaultServeMux.Handle("/api/v0/secure/set-default-algorithm", secure.Server{Db: q, Next: secure.SetDefaultAlgorithm, OrigDB: db})

	address := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting server on address %s\n", address)
	corsOptions := cors.Options{
		AllowPrivateNetwork: true,
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
	}
	log.Fatal(http.ListenAndServe(address, cors.New(corsOptions).Handler(http.DefaultServeMux)))

}
