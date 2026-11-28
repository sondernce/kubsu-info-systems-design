package file

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"transportation/driver"
)

type testDriverRepository interface {
	GetByID(int) (*driver.Driver, error)
	GetKNShortList(int, int) ([]*driver.DriverShort, error)
	SortByLastName()
	Add(*driver.Driver) (int, error)
	Replace(int, *driver.Driver) error
	Delete(int) error
	GetCount() (int, error)
}

func newYAMLRepository(t *testing.T) *DriverRepYAML {
	t.Helper()
	r, err := NewDriverRepYAML(filepath.Join(t.TempDir(), "drivers.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func testFileRepositories(t *testing.T, check func(*testing.T, testDriverRepository)) {
	t.Helper()
	t.Run("JSON", func(t *testing.T) { check(t, newRepository(t)) })
	t.Run("YAML", func(t *testing.T) { check(t, newYAMLRepository(t)) })
}

func TestYAMLRoundTrip(t *testing.T) {
	r := newYAMLRepository(t)
	if err := r.WriteAll(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(r.path)
	if err != nil || string(data) != "[]\n" {
		t.Fatalf("пустой YAML-список: %q, %v", data, err)
	}
	if err := r.ReadAll(); err != nil {
		t.Fatal(err)
	}
	d := newDriver(t, "Иванов")
	if _, err := r.Add(d); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Add(newDriver(t, "Петров")); err != nil {
		t.Fatal(err)
	}
	if err := r.WriteAll(); err != nil {
		t.Fatal(err)
	}
	loaded, err := NewDriverRepYAML(r.path)
	if err != nil {
		t.Fatal(err)
	}
	if err := loaded.ReadAll(); err != nil {
		t.Fatal(err)
	}
	got, err := loaded.GetByID(1)
	if err != nil {
		t.Fatal(err)
	}
	if countDrivers(t, loaded) != 2 || got.LastName() != d.LastName() ||
		got.FirstName() != d.FirstName() || got.MiddleName() != d.MiddleName() ||
		got.PhoneNumber() != d.PhoneNumber() || got.Experience() != d.Experience() ||
		got.Initials() != "И.П." {
		t.Fatalf("поля потеряны при сохранении YAML: %v", got)
	}
	if id, err := loaded.Add(d); err != nil || id != 3 {
		t.Fatalf("ID после загрузки YAML: %d, %v", id, err)
	}
	if err := loaded.Replace(1, newDriver(t, "Сидоров")); err != nil {
		t.Fatal(err)
	}
	if err := loaded.Delete(2); err != nil {
		t.Fatal(err)
	}
	if err := loaded.WriteAll(); err != nil {
		t.Fatal(err)
	}
	if err := r.ReadAll(); err != nil {
		t.Fatal(err)
	}
	got, err = r.GetByID(1)
	if err != nil || got.LastName() != "Сидоров" || countDrivers(t, r) != 2 {
		t.Fatalf("изменения не сохранились: %v, %v", got, err)
	}
	if _, err := r.GetByID(2); !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("удаление не сохранилось: %v", err)
	}
}

const validYAMLDriver = `- id: 8
  last_name: Петров
  first_name: Иван
  middle_name: Петрович
  phone_number: "+79161234567"
  experience: 10
`

func TestYAMLReadErrorsPreserveState(t *testing.T) {
	r := newYAMLRepository(t)
	if err := r.ReadAll(); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("отсутствующий файл: %v", err)
	}
	if _, err := r.Add(newDriver(t, "Иванов")); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{
		"empty":           "",
		"null":            "null",
		"mapping":         "id: 8",
		"syntax":          "[",
		"invalid fields":  "- {}",
		"duplicate ID":    validYAMLDriver + validYAMLDriver,
		"invalid second":  validYAMLDriver + "- {}",
		"unknown field":   validYAMLDriver + "  typo: 10\n",
		"duplicate key":   validYAMLDriver + "  id: 9\n",
		"invalid phone":   strings.ReplaceAll(validYAMLDriver, "+79161234567", "123"),
		"second document": validYAMLDriver + "---\n[]\n",
		"broken second":   validYAMLDriver + "---\n[",
	} {
		t.Run(name, func(t *testing.T) {
			if err := os.WriteFile(r.path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			if err := r.ReadAll(); err == nil {
				t.Fatalf("ожидалась ошибка для %q", data)
			}
			if countDrivers(t, r) != 1 {
				t.Fatal("ошибочная загрузка изменила список")
			}
			if _, err := r.GetByID(1); err != nil {
				t.Fatal(err)
			}
		})
	}
	if id, err := r.Add(newDriver(t, "Сидоров")); err != nil || id != 2 {
		t.Fatalf("ошибочная загрузка изменила счётчик ID: %d, %v", id, err)
	}
	if err := os.WriteFile(r.path, []byte("# Список водителей\n"+validYAMLDriver), 0600); err != nil {
		t.Fatal(err)
	}
	if err := r.ReadAll(); err != nil {
		t.Fatal(err)
	}
	if id, err := r.Add(newDriver(t, "Сидоров")); err != nil || id != 9 {
		t.Fatalf("ID после загрузки: %d, %v", id, err)
	}
}

func TestYAMLInvalidPath(t *testing.T) {
	if _, err := NewDriverRepYAML(" "); err == nil {
		t.Fatal("ожидалась ошибка для пустого пути")
	}
	r, err := NewDriverRepYAML(filepath.Join(t.TempDir(), "missing", "drivers.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := r.WriteAll(); err == nil {
		t.Fatal("запись в отсутствующую папку должна вернуть ошибку")
	}
}

func countDrivers(t *testing.T, r testDriverRepository) int {
	t.Helper()
	count, err := r.GetCount()
	if err != nil {
		t.Fatal(err)
	}
	return count
}
