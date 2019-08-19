package api

import "os"

var (
	MYSQL_USERNAME = ""
	MYSQL_PASSWORD = ""
	MYSQL_HOST     = ""
	MYSQL_PORT     = ""
	MYSQL_DB       = ""
	URL            = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	APP_PORT       = ""
)

//LoadConfig load configuration from current env
func LoadConfig() {
	APP_PORT = getEnv("APP_PORT", "3000")
	MYSQL_USERNAME = getEnv("MYSQL_USERNAME", "root")
	MYSQL_PASSWORD = getEnv("MYSQL_PASSWORD", "password")
	MYSQL_HOST = getEnv("MYSQL_HOST", "localhost")
	MYSQL_PORT = getEnv("MYSQL_PORT", "3308")
	MYSQL_DB = getEnv("MYSQL_DB", "cubes")
}

func getEnv(key, fallBack string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallBack
}
