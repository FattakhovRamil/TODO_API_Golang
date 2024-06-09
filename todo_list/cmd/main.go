package main

import (
	"os"
	"os/signal"
	"syscall"
	_ "todo_list/docs"
	_ "todo_list/models"
	server "todo_list/server"
)

// @title Todo List API
// @version 1.0
// @description This is a simple todo list API.
// @host localhost:3001
// @BasePath /
func main() {
	go server.ServerStart()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

}
