package config

const configFileName = ".burrow.json"

type Config struct {
	DbFile      string `json:"db_file"`
	DefaultPort string `json:"default_port"`
}
