package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/do"
	dbUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/db_utils"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1/client_serverv1connect"
	"github.com/hashicorp/mdns"
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

	port := 80
	if _, portStr, err := net.SplitHostPort(address); err == nil {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	} else if p, err := strconv.Atoi(strings.TrimPrefix(address, ":")); err == nil {
		port = p
	}

	host, _ := os.Hostname()
	info := []string{"Forget About It Server"}
	service, err := mdns.NewMDNSService(host, "_forgetaboutit._tcp", "", "", port, nil, info)
	if err != nil {
		log.Printf("Failed to create mDNS service: %v", err)
	} else {
		mdnsServer, err := mdns.NewServer(&mdns.Config{Zone: service})
		if err != nil {
			log.Printf("Failed to start mDNS server: %v", err)
		} else {
			defer mdnsServer.Shutdown()
			fmt.Printf("Registered mDNS service _forgetaboutit._tcp for %s on port %d\n", host, port)
		}
	}

	fmt.Printf("Starting server on address %s\n", address)
	corsOptions := cors.Options{
		AllowPrivateNetwork: true,
		AllowedOrigins:      []string{"*"},
		AllowedHeaders:      []string{"*"},
	}
	log.Fatal(http.ListenAndServe(address, h2c.NewHandler(cors.New(corsOptions).Handler(http.DefaultServeMux), &http2.Server{})))

}
