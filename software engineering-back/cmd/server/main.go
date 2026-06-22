package main

import "software_engineering/internal/app"

// main 程序入口，创建应用实例并启动 HTTP 服务
func main() {
	application := app.New()
	application.Initialize()
	application.Run()
}
