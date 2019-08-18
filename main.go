package main

func init() {
	LoadConfig()
}

func main() {
	app := &App{}
	app.Init()
	app.Run()
}
