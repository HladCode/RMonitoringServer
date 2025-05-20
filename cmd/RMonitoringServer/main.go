package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/HladCode/RMonitoringServer/internal/config"
	_ "github.com/HladCode/RMonitoringServer/internal/http-server/handlers/getData_old" // getData
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/getTime"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/isConnectionGood"
	receivedata "github.com/HladCode/RMonitoringServer/internal/http-server/handlers/receiveData"
	senddatafromday "github.com/HladCode/RMonitoringServer/internal/http-server/handlers/sendDataFromDay"
	DB "github.com/HladCode/RMonitoringServer/internal/storage/timeScaleDB"

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

	db, err := DB.NewDatabase(conf.BD_connect_parametrs)
	if err != nil {
		log.Fatal("Can not make a db: ", err)
	}
	if err := db.InitFromFile("configs/init.sql"); err != nil {
		log.Fatal("Can not init a db: ", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", isConnectionGood.New()).Methods("GET")
	router.HandleFunc("/getDayData", senddatafromday.New(db)).Methods("Get")

	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/time", getTime.New()).Methods("Get")
	api.HandleFunc("/sendData", receivedata.New(db)).Methods("POST")

	//TODO: make normal log
	log.Println("Start Listening...")
	log.Fatal(http.ListenAndServe(conf.Host+":"+conf.Port, router))
}

/*

	//dataSaver := simpleStorage.NewStorage()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", isConnectionGood.New()).Methods("GET")
	//router.HandleFunc("/data", getData.New(dataSaver)).Methods("POST")
	//router.HandleFunc("/metrics", exportMetrics.New(dataBuffer)).Methods("GET")
	router.HandleFunc("/time", getTime.New()).Methods("Get")

	//TODO: make normal log
	log.Println("Start Listening...")
	log.Fatal(http.ListenAndServe(conf.Host+":"+conf.Port, router))

*/
