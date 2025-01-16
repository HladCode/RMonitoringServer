package storage

type Storage interface {
	SaveTemperature(ID, Purpose, SensorPinNumber, Timestamp string, data float64) error
	GetTempreature(ID, Purpose, SensorPinNumber, data string) string
}

// type ObjectData struct {
// 	Timestamp        string
// 	RefrigeratorPath string
// 	Tempreature      float64
// }
// func NewObjectData(timestamp, refrigeratorPath string, temp float64) ObjectData {
// 	return ObjectData{
// 		Timestamp:        timestamp,
// 		RefrigeratorPath: refrigeratorPath,
// 		Tempreature:      temp,
// 	}
// }
