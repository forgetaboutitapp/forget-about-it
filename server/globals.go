package server

import (
	"embed"
	"path/filepath"

	"github.com/adrg/xdg"
)

//go:embed sql/migrations/*
var DDL embed.FS

var DBFilename = filepath.Join(xdg.StateHome, "forget-about-it.sqlite3")
