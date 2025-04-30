package main

import (
	"log"
	"gopr/config"
	"gopr/server"
	_ "github.com/lib/pq"
)

func main(){
	log.Println("Initializing app")
	config := config.InitConfig("runners")
	log.Println("Initializing database")
	dbhandler := server.InitDatabase(config)
	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbhandler)
	httpServer.Start()
}