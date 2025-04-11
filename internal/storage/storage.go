package storage

type Data_unit struct {
	ID              string  `json:"ID"`
	SensorPinNumber int     `json:"sID"`
	Timestamp       string  `json:"dt"`
	Data            float64 `json:"d"`
	//	Purpose         string убрать также с ESP32 эта инфа будет в другой таблице в бд
}

type Storage interface {
	AddNewData(readings []Data_unit) error
	GetDataFromDay(ID string, sensor_ID, day, month, year int) error
}
