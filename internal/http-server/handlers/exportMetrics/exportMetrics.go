package exportMetrics

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
)

type dataExporter interface {
	GetFormatedTempreature() list.List
}

func New(exporter dataExporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: make normal log
		log.Println("New request: ", r.RequestURI)
		responseText := exporter.GetFormatedTempreature()
		//var responseText string

		// TODO: make response formating for prometeus at storage side (data must formating right after it reach server)
		// TODO id: 1
		// for e := dataList.Front(); e != nil; e = e.Next() {
		// 	element := e.Value.(storage.ObjectData)
		// 	b := strings.Split(element.Timestamp, " ")
		// 	responseText += fmt.Sprintf("temperature{date=\"%s\", time=\"%s\", object=\"%s\"} %f\n", b[0], b[1], element.RefrigeratorPath, element.Tempreature)
		// }
		fmt.Fprint(w, responseText)
	}
}
