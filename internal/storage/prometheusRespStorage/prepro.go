package prometheusRespStorage

import (
	"fmt"
	"strings"
	"sync"

	"github.com/HladCode/RMonitoringServer/internal/lib/e"
)

type prometheusResponse struct {
	response string
	//maxDataCountInBuf int
}

func NewStorage() *prometheusResponse {
	return &prometheusResponse{
		response: "",
		//maxDataCountInBuf: maxDataCountInBuf,
	}
}

func (s *prometheusResponse) SaveTemperature(timestamp, phoneNumber, refrigeratorPath string, temp float64) error { // s *promStorage
	date_time := strings.Split(timestamp, " ")
	if len(date_time) != 2 {
		return e.Wrap("Timestamp is not correct", nil)
	}

	responseText := fmt.Sprintf("temperature{date=\"%s\", time=\"%s\", object=\"%s\", phone_number=\"%s\"} %f\n",
		date_time[0],
		date_time[1],
		refrigeratorPath,
		phoneNumber,
		temp)
	var mu sync.Mutex
	mu.Lock()
	s.response += responseText
	mu.Unlock()
	return nil
}

func (s *prometheusResponse) GetTempreature() string {
	defer func() {
		var mu sync.Mutex
		mu.Lock()
		s.response = ""
		mu.Unlock()
	}()
	return s.response
}
