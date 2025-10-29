package config

import "database/sql"

type ClientConfig struct {
	Db *sql.DB
}
