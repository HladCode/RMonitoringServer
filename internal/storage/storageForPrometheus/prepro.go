package storageforprometheus

import (
	"container/list"

	st "github.com/HladCode/RMonitoringServer/internal/storage"
)

type promStorage struct {
	maxDataCountInBuf int
	buf               map[string]*list.List
}

// TODO: delete maxDataCountInBuf
func NewStorage(maxDataCountInBuf int) promStorage {
	return promStorage{
		maxDataCountInBuf: maxDataCountInBuf,
		buf:               make(map[string]*list.List),
	}
}

func (s promStorage) SaveTemperature(timestamp, refrigeratorPath string, temp float64) error { // s *promStorage
	newObj := st.NewObjectData(timestamp, refrigeratorPath, temp)

	if s.buf[refrigeratorPath] == nil {
		s.buf[refrigeratorPath] = list.New()
	}

	if s.buf[refrigeratorPath].Len() < s.maxDataCountInBuf {
		s.buf[refrigeratorPath].PushBack(newObj)
	} else {
		s.buf[refrigeratorPath].Remove(s.buf[refrigeratorPath].Front())
		s.buf[refrigeratorPath].PushBack(newObj)
	}

	return nil
}

func (s promStorage) GetTempreature() list.List {
	resList := list.New()

	for k := range s.buf {
		resList.PushBackList(s.buf[k])

		delete(s.buf, k)
	}

	return *resList
}
