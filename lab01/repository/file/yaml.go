package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"go.yaml.in/yaml/v3"
)

// DriverRepYAML читает и записывает список водителей в формате YAML.
type DriverRepYAML struct {
	driverFileRepository
}

// NewDriverRepYAML создаёт пустой репозиторий. Загрузка выполняется через ReadAll.
func NewDriverRepYAML(path string) (*DriverRepYAML, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("путь к YAML-файлу не может быть пустым")
	}
	return &DriverRepYAML{driverFileRepository: newDriverFileRepository(path)}, nil
}

//a

func (r *DriverRepYAML) ReadAll() error {
	return r.readAll(decodeYAML)
}

//b

func (r *DriverRepYAML) WriteAll() error {
	return r.writeAll(func(records []driverRecord) ([]byte, error) {
		return yaml.Marshal(records)
	})
}

func decodeYAML(data []byte) ([]driverRecord, error) {
	var records []driverRecord
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&records); err != nil {
		return nil, fmt.Errorf("разбор YAML: %w", err)
	}
	if records == nil {
		return nil, errors.New("ожидается YAML-список водителей")
	}

	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("разбор конца YAML: %w", err)
		}
		return nil, errors.New("ожидается один YAML-документ со списком водителей")
	}
	return records, nil
}
