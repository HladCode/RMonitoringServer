package storage

type Storage interface {
	SaveTemperature(timestamp, phoneNumber, refrigeratorPath string, temp float64) error
	GetTempreature() string
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
