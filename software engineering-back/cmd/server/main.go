package main

import "software_engineering/internal/app"

func main() {
	application := app.New()
	application.Initialize()
	application.Run()
}
