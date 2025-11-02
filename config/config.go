package config

import (
	"database/sql"
	"net/http"
	"time"
)

type ClientConfig struct {
	Db            *sql.DB
	GlobalTimeout time.Duration
	Client        *http.Client
}
