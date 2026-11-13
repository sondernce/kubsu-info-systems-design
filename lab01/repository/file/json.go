package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"transportation/driver"
)

var ErrDriverNotFound = errors.New("водитель не найден")

// DriverRepJSON хранит рабочий список в памяти; ReadAll и WriteAll синхронизируют его с файлом.
type DriverRepJSON struct {
	path    string
	drivers []*driver.Driver
	lastID  int
	mapper  *driver.DriverMapper
}

type driverRecord struct {
	ID          int    `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	PhoneNumber string `json:"phone_number"`
	Experience  int    `json:"experience"`
}

// NewDriverRepJSON создаёт пустой репозиторий. Загрузка выполняется через ReadAll.
func NewDriverRepJSON(path string) (*DriverRepJSON, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("путь к JSON-файлу не может быть пустым")
	}
	return &DriverRepJSON{
		path:    path,
		drivers: make([]*driver.Driver, 0),
		mapper:  driver.NewDriverMapper(),
	}, nil
}

//a

func (r *DriverRepJSON) ReadAll() error {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return fmt.Errorf("чтение водителей: %w", err)
	}
	var records []driverRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("разбор JSON: %w", err)
	}
	if records == nil {
		return errors.New("ожидается JSON-массив водителей")
	}

	drivers := make([]*driver.Driver, 0, len(records))
	seen := make(map[int]bool, len(records))
	lastID := r.lastID
	for i, record := range records {
		if seen[record.ID] {
			return fmt.Errorf("повторяющийся ID водителя: %d", record.ID)
		}
		d, err := driver.NewDriver(record.ID, record.LastName, record.FirstName,
			record.MiddleName, record.PhoneNumber, record.Experience)
		if err != nil {
			return fmt.Errorf("запись %d: %w", i+1, err)
		}
		seen[record.ID] = true
		drivers = append(drivers, d)
		if record.ID > lastID {
			lastID = record.ID
		}
	}
	// Заменяем список только после проверки всех записей.
	r.drivers = drivers
	r.lastID = lastID
	return nil
}

//b

func (r *DriverRepJSON) WriteAll() error {
	records := make([]driverRecord, 0, len(r.drivers))
	for _, d := range r.drivers {
		records = append(records, driverRecord{
			ID: d.ID(), LastName: d.LastName(), FirstName: d.FirstName(),
			MiddleName: d.MiddleName(), PhoneNumber: d.PhoneNumber(), Experience: d.Experience(),
		})
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("кодирование водителей: %w", err)
	}
	// Временный файл защищает предыдущие данные от частичной записи.
	tmp, err := os.CreateTemp(filepath.Dir(r.path), ".drivers-*.json")
	if err != nil {
		return fmt.Errorf("создание временного файла: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(append(data, '\n')); err != nil {
		tmp.Close()
		return fmt.Errorf("запись водителей: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("закрытие временного файла: %w", err)
	}
	if err := os.Rename(tmp.Name(), r.path); err != nil {
		return fmt.Errorf("сохранение водителей: %w", err)
	}
	return nil
}

//c

func (r *DriverRepJSON) GetByID(id int) (*driver.Driver, error) {
	for _, d := range r.drivers {
		if d.ID() == id {
			return copyDriver(d, id)
		}
	}
	return nil, fmt.Errorf("%w: %d", ErrDriverNotFound, id)
}

//d

func (r *DriverRepJSON) GetKNShortList(k, n int) ([]*driver.DriverShort, error) {
	if k <= 0 || n <= 0 {
		return nil, errors.New("номер страницы и её размер должны быть больше нуля")
	}
	result := make([]*driver.DriverShort, 0)
	// Проверяем границы до умножения, чтобы избежать переполнения int.
	if len(r.drivers) == 0 || k-1 > (len(r.drivers)-1)/n {
		return result, nil
	}
	start := (k - 1) * n
	size := min(n, len(r.drivers)-start)
	for _, d := range r.drivers[start : start+size] {
		short := *r.mapper.ToShort(d)
		result = append(result, &short)
	}
	return result, nil
}

//e

func (r *DriverRepJSON) SortByLastName() {
	sort.SliceStable(r.drivers, func(i, j int) bool {
		return strings.ToLower(r.drivers[i].LastName()) < strings.ToLower(r.drivers[j].LastName())
	})
}

//f

func (r *DriverRepJSON) Add(d *driver.Driver) (int, error) {
	if r.lastID == int(^uint(0)>>1) {
		return 0, errors.New("закончились доступные ID")
	}
	id := r.lastID + 1
	stored, err := copyDriver(d, id)
	if err != nil {
		return 0, err
	}
	r.drivers = append(r.drivers, stored)
	r.lastID = id
	return id, nil
}

//g

func (r *DriverRepJSON) Replace(id int, d *driver.Driver) error {
	for i, current := range r.drivers {
		if current.ID() == id {
			stored, err := copyDriver(d, id)
			if err != nil {
				return err
			}
			r.drivers[i] = stored
			return nil
		}
	}
	return fmt.Errorf("%w: %d", ErrDriverNotFound, id)
}

//h

func (r *DriverRepJSON) Delete(id int) error {
	for i, d := range r.drivers {
		if d.ID() == id {
			copy(r.drivers[i:], r.drivers[i+1:])
			r.drivers[len(r.drivers)-1] = nil
			r.drivers = r.drivers[:len(r.drivers)-1]
			return nil
		}
	}
	return fmt.Errorf("%w: %d", ErrDriverNotFound, id)
}

//i

func (r *DriverRepJSON) GetCount() int {
	return len(r.drivers)
}

// Копия не даёт вызывающему коду менять данные репозитория в обход его методов.
func copyDriver(d *driver.Driver, id int) (*driver.Driver, error) {
	if d == nil {
		return nil, errors.New("водитель не может быть nil")
	}
	return driver.NewDriver(id, d.LastName(), d.FirstName(), d.MiddleName(), d.PhoneNumber(), d.Experience())
}
