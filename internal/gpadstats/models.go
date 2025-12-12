package gpadstats

import (
	"database/sql"
	"io"
)

type StatsLoaderConfig struct {
	Path   string
	Reader io.Reader
	DB     *sql.DB
}

type GeneCountStats struct {
	Count    int
	EcoCount int
}
