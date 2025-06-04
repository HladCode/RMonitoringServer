package timescaledb

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/HladCode/RMonitoringServer/internal/lib/e"
	"github.com/HladCode/RMonitoringServer/internal/storage"
)

func (db *Database) AddNewData(readings []storage.Data_unit) error {
	base := `INSERT INTO sensor_readings (timestamp, device_id, sensor_index, value) VALUES `
	args := []interface{}{}
	placeholders := []string{}

	for i, r := range readings {
		idx := i * 4 // бо 5 полів у кожному рядку
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4))
		args = append(args, r.Timestamp, r.ID, r.SensorPinNumber, r.Data)
	}

	query := base + strings.Join(placeholders, ",")
	_, err := db.pool.Exec(db.cntxt, query, args...)
	if err != nil {
		return e.Wrap("Can not add data", err)
	}

	return nil
}

func (db *Database) GetDataFromDay(ID string, sensor_ID, day, month, year int) (string, error) {
	rows, err := db.pool.Query(db.cntxt, `SELECT timestamp, value FROM get_sensor_data_for_day($1, $2, $3);`,
		ID, sensor_ID, fmt.Sprintf("%d-%d-%dT00:00:00+00:00", year, month, day))
	if err != nil {
		return "", e.Wrap("Can not get day data", err)
	}

	defer rows.Close()

	data := make(map[string]float64)

	for rows.Next() {
		var timestamp time.Time
		var value float64
		if err := rows.Scan(&timestamp, &value); err != nil {
			return "", e.Wrap("failed to scan row", err)
		}

		data[timestamp.Format(time.RFC3339)] = value // TODO: на esp32 при передачи данных указывай в таймштампе UTC+3
	}

	if err := rows.Err(); err != nil {
		return "", e.Wrap("rows iteration error", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", e.Wrap("failed to marshal data to JSON", err)
	}

	return string(jsonData), nil
}

func (db *Database) GetDevicesFromUserInJson(username string) (string, error) {
	var companyName string

	query := `
		SELECT company_name
		FROM users_companies AS uc
		WHERE uc.user_id = (
			SELECT id FROM users AS u WHERE u.username = $1
		)
	`

	err := db.pool.QueryRow(db.cntxt, query, username).Scan(&companyName)
	if err != nil {
		return "", err
	}

	rows, err := db.pool.Query(db.cntxt, `SELECT  place,  array_agg(device_id) AS device_ids FROM devices_place_companies WHERE  company_name = $1 GROUP BY place`,
		companyName)
	if err != nil {
		return "", e.Wrap("Can not get data", err)
	}

	defer rows.Close()

	devices := make(map[string][]string)
	for rows.Next() {
		var device_ids []string
		var place string
		if err := rows.Scan(&place, &device_ids); err != nil {
			return "", e.Wrap("failed to scan row", err)
		}

		devices[place] = device_ids
	}

	if err := rows.Err(); err != nil {
		return "", e.Wrap("rows iteration error", err)
	}

	jsonData, err := json.Marshal(devices)
	if err != nil {
		return "", e.Wrap("failed to marshal data to JSON", err)
	}

	return string(jsonData), nil
}

func (db *Database) GetSensorsFromDeviceIdInJson(device_id string) (string, error) {
	rows, err := db.pool.Query(db.cntxt, `SELECT  sensor_index,  measurement_type, meaning FROM sensor_index_meaning WHERE  device_id = $1`,
		device_id)
	if err != nil {
		return "", e.Wrap("Can not get data", err)
	}

	defer rows.Close()

	devices := make(map[int]string)
	for rows.Next() {
		var sensor_id int
		var measurement_type string
		var meaning string
		if err := rows.Scan(&sensor_id, &measurement_type, &meaning); err != nil {
			return "", e.Wrap("failed to scan row", err)
		}

		devices[sensor_id] = measurement_type + "\n" + meaning
	}

	if err := rows.Err(); err != nil {
		return "", e.Wrap("rows iteration error", err)
	}

	jsonData, err := json.Marshal(devices)
	if err != nil {
		return "", e.Wrap("failed to marshal data to JSON", err)
	}

	return string(jsonData), nil
}
