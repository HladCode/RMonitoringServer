package exportMetrics

import (
	"fmt"
	"log"
	"net/http"
)

type dataExporter interface {
	GetTempreature() string
}

func New(exporter dataExporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("New request: ", r.RequestURI)
		responseText := exporter.GetTempreature()
		fmt.Fprint(w, responseText)
	}
}

//var responseText string
// for e := dataList.Front(); e != nil; e = e.Next() {
// 	element := e.Value.(storage.ObjectData)
// 	b := strings.Split(element.Timestamp, " ")
// 	responseText += fmt.Sprintf("temperature{date=\"%s\", time=\"%s\", object=\"%s\"} %f\n", b[0], b[1], element.RefrigeratorPath, element.Tempreature)
// }
