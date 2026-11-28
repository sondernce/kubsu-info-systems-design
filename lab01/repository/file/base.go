package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"transportation/driver"
	"transportation/repository"
)

var ErrDriverNotFound = repository.ErrDriverNotFound

var (
	_ repository.DriverRepository = (*DriverRepJSON)(nil)
	_ repository.DriverRepository = (*DriverRepYAML)(nil)
)

// driverFileRepository хранит список водителей и выполняет общие операции с файлом.
type driverFileRepository struct {
	path    string
	drivers []*driver.Driver
	lastID  int
	mapper  *driver.DriverMapper
}

// driverRecord задаёт одинаковые имена полей в JSON и YAML.
type driverRecord struct {
	ID          int    `json:"id" yaml:"id"`
	LastName    string `json:"last_name" yaml:"last_name"`
	FirstName   string `json:"first_name" yaml:"first_name"`
	MiddleName  string `json:"middle_name" yaml:"middle_name"`
	PhoneNumber string `json:"phone_number" yaml:"phone_number"`
	Experience  int    `json:"experience" yaml:"experience"`
}

func newDriverFileRepository(path string) driverFileRepository {
	return driverFileRepository{path: path, mapper: driver.NewDriverMapper()}
}

func (r *driverFileRepository) readAll(decode func([]byte) ([]driverRecord, error)) error {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return fmt.Errorf("чтение водителей: %w", err)
	}
	records, err := decode(data)
	if err != nil {
		return err
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

func (r *driverFileRepository) writeAll(encode func([]driverRecord) ([]byte, error)) error {
	records := make([]driverRecord, 0, len(r.drivers))
	for _, d := range r.drivers {
		records = append(records, driverRecord{
			ID: d.ID(), LastName: d.LastName(), FirstName: d.FirstName(),
			MiddleName: d.MiddleName(), PhoneNumber: d.PhoneNumber(), Experience: d.Experience(),
		})
	}
	data, err := encode(records)
	if err != nil {
		return fmt.Errorf("кодирование водителей: %w", err)
	}
	// Временный файл защищает предыдущие данные от частичной записи.
	tmp, err := os.CreateTemp(filepath.Dir(r.path), ".drivers-*")
	if err != nil {
		return fmt.Errorf("создание временного файла: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
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

func (r *driverFileRepository) GetByID(id int) (*driver.Driver, error) {
	for _, d := range r.drivers {
		if d.ID() == id {
			return copyDriver(d, id)
		}
	}
	return nil, fmt.Errorf("%w: %d", ErrDriverNotFound, id)
}

//d

func (r *driverFileRepository) GetKNShortList(k, n int) ([]*driver.DriverShort, error) {
	return r.shortPage(r.drivers, k, n)
}

func (r *driverFileRepository) fileBase() *driverFileRepository {
	return r
}

func (r *driverFileRepository) shortPage(drivers []*driver.Driver, k, n int) ([]*driver.DriverShort, error) {
	if k <= 0 || n <= 0 {
		return nil, errors.New("номер страницы и её размер должны быть больше нуля")
	}
	result := make([]*driver.DriverShort, 0)
	// Проверяем границы до умножения, чтобы избежать переполнения int.
	if len(drivers) == 0 || k-1 > (len(drivers)-1)/n {
		return result, nil
	}
	start := (k - 1) * n
	size := min(n, len(drivers)-start)
	for _, d := range drivers[start : start+size] {
		short := *r.mapper.ToShort(d)
		result = append(result, &short)
	}
	return result, nil
}

//e

func (r *driverFileRepository) SortByLastName() {
	sort.SliceStable(r.drivers, func(i, j int) bool {
		return strings.ToLower(r.drivers[i].LastName()) < strings.ToLower(r.drivers[j].LastName())
	})
}

//f

func (r *driverFileRepository) Add(d *driver.Driver) (int, error) {
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

func (r *driverFileRepository) Replace(id int, d *driver.Driver) error {
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

func (r *driverFileRepository) Delete(id int) error {
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

func (r *driverFileRepository) GetCount() (int, error) {
	return len(r.drivers), nil
}

// Копия не даёт вызывающему коду менять данные репозитория в обход его методов.
func copyDriver(d *driver.Driver, id int) (*driver.Driver, error) {
	if d == nil {
		return nil, errors.New("водитель не может быть nil")
	}
	return driver.NewDriver(id, d.LastName(), d.FirstName(), d.MiddleName(), d.PhoneNumber(), d.Experience())
}
