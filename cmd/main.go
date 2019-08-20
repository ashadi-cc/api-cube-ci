package main

import (
	"api"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	api := api.New()
	api.Init()
	api.Run()
}
