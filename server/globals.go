package server

import (
	"embed"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
)

//go:embed sql/migrations/*
var DDL embed.FS

var DBFilename = filepath.Join(xdg.StateHome, "forget-about-it.sqlite3")

var MutexUsersWaiting sync.Mutex
var UsersWaiting map[int64]struct{} = make(map[int64]struct{})
