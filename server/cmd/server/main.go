package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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
	http.DefaultServeMux.Handle("/api/v0/get-token", login.Server{Db: q})
	http.DefaultServeMux.Handle("/api/v0/secure/generate-new-token", secure.Server{Db: q, Next: secure.GenerateNewToken})
	http.DefaultServeMux.Handle("/api/v0/secure/check-new-token", secure.Server{Db: q, Next: secure.CheckNewToken})
	http.DefaultServeMux.Handle("/api/v0/secure/delete-new-token", secure.Server{Db: q, Next: secure.DeleteNewToken})
	fmt.Println("Starting server")
	corsOptions := cors.Options{
		AllowPrivateNetwork: true,
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
	}
	log.Fatal(http.ListenAndServe(":8080", cors.New(corsOptions).Handler(http.DefaultServeMux)))

}
