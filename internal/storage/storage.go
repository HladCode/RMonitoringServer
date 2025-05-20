package storage

type Data_unit struct {
	ID              string  `json:"ID"`
	SensorPinNumber int     `json:"sID"`
	Timestamp       string  `json:"dt"`
	Data            float64 `json:"d"`
}

type User_data struct {
	User_name       string
	Email           string
	Hashed_password string
	ID              int
}

type Storage interface {
	AddNewData(readings []Data_unit) error
	GetDataFromDay(ID string, sensor_ID, day, month, year int) error
}
