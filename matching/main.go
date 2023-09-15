package main

import (
	"github.com/spf13/viper"
	"matching/engine"
	"matching/handler"
	"matching/log"
	"matching/middleware"
	"matching/process"
	"net/http"
)

func init() {
	//initViper()
	initLog()

	engine.Init()
	middleware.Init()
	process.Init()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/openMatching", handler.OpenMatching)
	mux.HandleFunc("/closeMatching", handler.CloseMatching)
	mux.HandleFunc("/handleOrder", handler.HandleOrder)

	viper.Set("server.port", ":18080")
	log.Info("HTTP ListenAndServe at port %s", viper.GetString("server.port"))
	if err := http.ListenAndServe(viper.GetString("server.port"), mux); err != nil {
		panic(err)
	}
}

func initLog() {
	err := log.Init("./", "matching.log", "m", log.DEBUG)
	if err != nil {
		panic(err)
	}
}
