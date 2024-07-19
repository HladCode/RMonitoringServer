package exportMetrics

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HladCode/RMonitoringServer/internal/storage"
)

type dataExporter interface {
	GetTempreature() list.List
}

func New(exporter dataExporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: make normal log
		log.Println("New request: ", r.RequestURI)
		dataList := exporter.GetTempreature()
		var responseText string

		for e := dataList.Front(); e != nil; e = e.Next() {
			element := e.Value.(storage.ObjectData)
			b := strings.Split(element.Timestamp, " ")
			responseText += fmt.Sprintf("temperature{date=\"%s\", time=\"%s\", object=\"%s\"} %f\n", b[0], b[1], element.RefrigeratorPath, element.Tempreature)
		}
		fmt.Fprint(w, responseText)
	}
}
