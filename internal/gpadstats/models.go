package gpadstats

import (
	"database/sql"
	"io"
	"os"
)

type StatsLoaderConfig struct {
	Path   string
	File   *os.File
	Reader io.Reader
	DB     *sql.DB
}

type GeneCountStats struct {
	Count    int
	EcoCount int
}