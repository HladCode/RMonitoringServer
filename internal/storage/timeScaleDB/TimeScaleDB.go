package timescaledb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/HladCode/RMonitoringServer/internal/lib/e"
	"github.com/HladCode/RMonitoringServer/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool  *pgxpool.Pool
	cntxt context.Context
}

func NewDatabase(connStr string) (*Database, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, e.Wrap("can not open DB", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, e.Wrap("can not connect to DB", err)
	}

	return &Database{pool: pool, cntxt: ctx}, nil
}

func (db *Database) InitFromFile(path string) error {
	// TODO просто сделать проверку на наличие таблиц,
	// а то некоторые запросы плохо исполняются по много раз
	content, err := os.ReadFile(path)
	if err != nil {
		return e.Wrap("Can not read file", err)
	}

	_, err = db.pool.Exec(db.cntxt, string(content))
	if err != nil {
		return e.Wrap("Exec schema error", err)
	}

	return nil
}

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
