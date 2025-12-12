package gpadstats

import (
	"database/sql"
	"os"
)

type StatsLoaderConfig struct {
	Path string
	File *os.File
	DB   *sql.DB
}

type GeneCountStats struct {
	Count    int
	EcoCount int
}
