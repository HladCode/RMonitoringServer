package storage

import "container/list"

// type ObjectData struct {
// 	Timestamp        string
// 	RefrigeratorPath string
// 	Tempreature      float64
// }

type Storage interface {
	// TODO id: 1
	//
	SaveTemperature(timestamp, phoneNumber, refrigeratorPath string, temp float64) error
	GetTempreature(refrigeratorPath string) list.List
	// GetRefrigerators() []string
}

// func NewObjectData(timestamp, refrigeratorPath string, temp float64) ObjectData {
// 	return ObjectData{
// 		Timestamp:        timestamp,
// 		RefrigeratorPath: refrigeratorPath,
// 		Tempreature:      temp,
// 	}
// }
