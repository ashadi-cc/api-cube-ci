package api

import "os"

//DbConfig represent database config model
type DbConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

//AppConfig represent app config model
type AppConfig struct {
	Port   string
	XMLUrl string
}

//Config represent config model
type Config struct {
	Db  *DbConfig
	App *AppConfig
}

//LoadConfig set config variables from current process environment
func LoadConfig() *Config {
	dbConfig := DbConfig{
		Username: getEnv("MYSQL_USERNAME", "root"),
		Password: getEnv("MYSQL_PASSWORD", "password"),
		Host:     getEnv("MYSQL_HOST", "localhost"),
		Port:     getEnv("MYSQL_PORT", "3308"),
		Database: getEnv("MYSQL_DB", "cubes"),
	}

	appConfig := AppConfig{
		Port:   getEnv("APP_PORT", "3000"),
		XMLUrl: getEnv("XML_URL", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"),
	}

	return &Config{Db: &dbConfig, App: &appConfig}
}

//get Env value from current process
func getEnv(key, fallBack string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallBack
}
