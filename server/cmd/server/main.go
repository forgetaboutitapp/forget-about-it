package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/login"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
	dbUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/db_utils"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/rs/cors"
)

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	db, err := dbUtils.OpenDatabase(ctx)
	if err != nil {
		panic(err)
	}
	q := sql_queries.New(db)
	realSub, err := fs.Sub(server.Files,"web")
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

	fmt.Println("Starting server")
	corsOptions := cors.Options{
		AllowPrivateNetwork: true,
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
	}
	log.Fatal(http.ListenAndServe(":8080", cors.New(corsOptions).Handler(http.DefaultServeMux)))

}
