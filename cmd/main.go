package main

import "api"

func main() {
	api := api.New()
	api.Init()
	api.Run()
}
