package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/HladCode/RMonitoringServer/internal/config"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/exportMetrics"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/getData"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/isConnectionGood"

	"github.com/HladCode/RMonitoringServer/internal/storage/prometheusRespStorage"
	"github.com/gorilla/mux"
)

// https://stackoverflow.com/questions/11706215/how-can-i-fix-the-git-error-object-file-is-empty
func main() {
	ConfigPath := flag.String(
		"ConfigPath",
		"",
		"",
	)

	flag.Parse()

	switch *ConfigPath {
	case "":
		log.Fatal("ConfigPath is not specified")
	case ".":
		*ConfigPath = "configs/prod.json"
	case "..":
		*ConfigPath = "configs/local.json"
	}

	conf := config.MustRead(*ConfigPath)

	dataBuffer := prometheusRespStorage.NewStorage()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", isConnectionGood.New()).Methods("GET")
	router.HandleFunc("/data", getData.New(dataBuffer)).Methods("POST")
	router.HandleFunc("/metrics", exportMetrics.New(dataBuffer)).Methods("GET")

	//TODO: make normal log
	log.Println("Start Listening...")
	log.Fatal(http.ListenAndServe(conf.Host+":"+conf.Port, router))
}
