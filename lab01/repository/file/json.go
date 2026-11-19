package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// DriverRepJSON читает и записывает список водителей в формате JSON.
type DriverRepJSON struct {
	driverFileRepository
}

// NewDriverRepJSON создаёт пустой репозиторий. Загрузка выполняется через ReadAll.
func NewDriverRepJSON(path string) (*DriverRepJSON, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("путь к JSON-файлу не может быть пустым")
	}
	return &DriverRepJSON{driverFileRepository: newDriverFileRepository(path)}, nil
}

//a

func (r *DriverRepJSON) ReadAll() error {
	return r.readAll(decodeJSON)
}

//b

func (r *DriverRepJSON) WriteAll() error {
	return r.writeAll(func(records []driverRecord) ([]byte, error) {
		data, err := json.MarshalIndent(records, "", "  ")
		return append(data, '\n'), err
	})
}

func decodeJSON(data []byte) ([]driverRecord, error) {
	var records []driverRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("разбор JSON: %w", err)
	}
	if records == nil {
		return nil, errors.New("ожидается JSON-массив водителей")
	}
	return records, nil
}
