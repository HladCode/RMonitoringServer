package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/HladCode/RMonitoringServer/internal/config"
	_ "github.com/HladCode/RMonitoringServer/internal/http-server/handlers/monitoring_device_handlers/getData_old" // getData
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/monitoring_device_handlers/getTime"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/monitoring_device_handlers/isConnectionGood"
	receivedata "github.com/HladCode/RMonitoringServer/internal/http-server/handlers/monitoring_device_handlers/receiveData"
	getdevices "github.com/HladCode/RMonitoringServer/internal/http-server/handlers/user_handlers/getDevices"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/user_handlers/login"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/user_handlers/refresh_jwt"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/user_handlers/register"
	sendDataFromDay "github.com/HladCode/RMonitoringServer/internal/http-server/handlers/user_handlers/sendDataFromDay"
	api_jwt "github.com/HladCode/RMonitoringServer/internal/lib/api/jwt"
	"github.com/HladCode/RMonitoringServer/internal/middleware"
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
	api_jwt.SetSecretKey(conf.KeyJWT)

	db, err := DB.NewDatabase(conf.BD_connect_parametrs)
	if err != nil {
		log.Fatal("Can not make a db: ", err)
	}
	if err := db.InitFromFile(conf.SQLFILEPATH); err != nil {
		log.Fatal("Can not init a db: ", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", isConnectionGood.New()).Methods("GET")

	user := router.PathPrefix("/user").Subrouter()
	user.Use(middleware.JWTMiddleware)
	user.HandleFunc("/getDayData", sendDataFromDay.New(db)).Methods("Post")
	user.HandleFunc("/getDevices", getdevices.New(db)).Methods("Post") //  TODO: uncomment

	user_authentication := router.PathPrefix("/auth").Subrouter()
	user_authentication.Use(middleware.AuthenticationRateLimiter(10, 35*time.Minute))
	user_authentication.HandleFunc("/register", register.New(db)).Methods("Post")
	user_authentication.HandleFunc("/login", login.New(db)).Methods("Post")
	user_authentication.HandleFunc("/refresh", refresh_jwt.New(db)).Methods("Post")

	api := router.PathPrefix("/api").Subrouter()
	//api.Use() TODO: API TOKEN MIDDLEWARE
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
