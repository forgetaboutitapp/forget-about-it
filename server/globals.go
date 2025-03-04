package server

import (
	"embed"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
)

//go:embed sql/migrations/*
var DDL embed.FS

var DBFilename = filepath.Join(xdg.StateHome, "forget-about-it.sqlite3")

var MutexUsersWaiting sync.Mutex
var UsersWaiting map[uuid.UUID]struct{} = make(map[uuid.UUID]struct{})
