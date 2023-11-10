package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Storage interface {
	AddSensor( /*Country, Company, City, Refrigerator, Sensor string*/ ) error
}

// NewDatabase(storagePath string) (*Storage, error)
