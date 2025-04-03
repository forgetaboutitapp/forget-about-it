package server

import (
	"embed"
	"sync"
)

//go:embed sql/migrations/*
var DDL embed.FS

//go:embed web/*
var Files embed.FS

var DBFilename = ""

var MutexUsersWaiting sync.Mutex
var UsersWaiting map[int64]struct{} = make(map[int64]struct{})
var DbLock sync.RWMutex
