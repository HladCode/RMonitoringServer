package simpleStorage

import (
	"fmt"
	"os"
	"strings"

	"github.com/HladCode/RMonitoringServer/internal/lib/e"
)

type SimpleStorage struct {
}

func NewStorage() *SimpleStorage {
	return &SimpleStorage{}
}

// SaveTemperature(ID, Purpose, SensorPinNumber, Timestamp string, data float64) error
// GetTempreature(ID, Purpose, SensorPinNumber, data string) string

// TODO: change request esp32 json timestamp so all delimeters will be " "
func (s SimpleStorage) SaveData(ID, Purpose, SensorPinNumber, Timestamp string, data float64) error {
	splitedTimestamp := strings.Split(Timestamp, " ")
	path := fmt.Sprintf("%s/%s/%s/%s/%s/%s/", ID, Purpose, SensorPinNumber, splitedTimestamp[0], splitedTimestamp[1], splitedTimestamp[2])
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return e.WrapIfErr("can not make directory", err)
		}
	}

	file, err := os.OpenFile(fmt.Sprintf("%s%s!%s!%s.txt", path, splitedTimestamp[0], splitedTimestamp[1], splitedTimestamp[2]),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return e.WrapIfErr("can not make data file", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s:%s:%s %f\n", splitedTimestamp[3], splitedTimestamp[4], splitedTimestamp[5], data))
	if err != nil {
		return e.WrapIfErr("can not write to data file", err)
	}

	return nil
}

// func (s SimpleStorage) GetTempreature(ID, Purpose, SensorPinNumber, data string) string {

// }
