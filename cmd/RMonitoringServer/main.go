package main

import (
	"log"
	"net/http"

	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/exportMetrics"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/getData"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/isConnectionGood"
	storageforprometheus "github.com/HladCode/RMonitoringServer/internal/storage/storageForPrometheus"
	"github.com/gorilla/mux"
)

// https://stackoverflow.com/questions/11706215/how-can-i-fix-the-git-error-object-file-is-empty
func main() {
	dataStorage := storageforprometheus.NewStorage(7)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", isConnectionGood.New()).Methods("GET")
	router.HandleFunc("/data", getData.New(dataStorage)).Methods("POST")
	router.HandleFunc("/metrics", exportMetrics.New(dataStorage)).Methods("GET")

	//TODO: make normal log
	log.Println("Start Listening...")
	log.Fatal(http.ListenAndServe("192.168.0.103:1488", router))
}
