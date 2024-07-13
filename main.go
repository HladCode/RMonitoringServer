package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// User represents a user structure
type Data struct {
	Tempreature string `json:"t"`
	Path        string `json:"p"`
}

// HelloWorld handles GET requests to "/"
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

// CreateUser handles POST requests to "/users"
func GetData(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print("Reading body has been failed", "\n")
		fmt.Fprintf(w, "Error")
		return
	}
	var dat Data
	err = json.Unmarshal(reqBody, &dat)
	if err != nil {
		log.Print("Marshaling has been failed", "\n", reqBody, "\n")
		fmt.Fprintf(w, "Error")
		return
	}

	log.Println(dat.Tempreature, "°C", ", ", dat.Path)
	fmt.Fprintf(w, "Okk")
}

func handleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HelloWorld).Methods("GET")
	router.HandleFunc("/data", GetData).Methods("POST")

	log.Println("Start Listening...")
	log.Fatal(http.ListenAndServe("192.168.0.103:1488", router))
}

func main() {
	handleRequests()
}
