package db

import (
	"embed"
)

//go:embed schema/*.sql
var Schema embed.FS