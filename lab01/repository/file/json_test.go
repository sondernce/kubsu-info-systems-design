package file

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"transportation/driver"
)

func newRepository(t *testing.T) *DriverRepJSON {
	t.Helper()
	r, err := NewDriverRepJSON(filepath.Join(t.TempDir(), "drivers.json"))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func newDriver(t *testing.T, name string) *driver.Driver {
	t.Helper()
	d, err := driver.NewDriver(99, name, "Иван", "Петрович", "+79161234567", 10)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestCRUDAndGeneratedIDs(t *testing.T) {
	r := newRepository(t)
	d := newDriver(t, "Иванов")
	id, err := r.Add(d)
	if err != nil || id != 1 || d.ID() != 99 {
		t.Fatalf("Add: id=%d, исходный ID=%d, err=%v", id, d.ID(), err)
	}
	if err := d.SetLastName("Сидоров"); err != nil {
		t.Fatal(err)
	}
	got, err := r.GetByID(id)
	if err != nil || got.LastName() != "Иванов" {
		t.Fatalf("репозиторий должен хранить отдельную копию: %v, %v", got, err)
	}
	if err := got.SetLastName("Петров"); err != nil {
		t.Fatal(err)
	}
	again, err := r.GetByID(id)
	if err != nil || again.LastName() != "Иванов" {
		t.Fatalf("GetByID должен возвращать копию: %v, %v", again, err)
	}
	if err := r.Replace(id, d); err != nil {
		t.Fatal(err)
	}
	got, err = r.GetByID(id)
	if err != nil || got.ID() != id || got.LastName() != "Сидоров" || r.GetCount() != 1 {
		t.Fatalf("Replace: %v, %v, count=%d", got, err, r.GetCount())
	}
	if err := r.Delete(id); err != nil {
		t.Fatal(err)
	}
	if r.GetCount() != 0 {
		t.Fatal("после удаления список должен быть пустым")
	}
	if _, err := r.GetByID(id); !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("ожидалась ErrDriverNotFound: %v", err)
	}
	if err := r.Replace(id, d); !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("замена отсутствующего ID: %v", err)
	}
	if err := r.Delete(id); !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("повторное удаление: %v", err)
	}
	if id, err := r.Add(d); err != nil || id != 2 {
		t.Fatalf("удалённый ID не должен использоваться снова: %d, %v", id, err)
	}
	if _, err := r.Add(nil); err == nil {
		t.Fatal("ожидалась ошибка для nil")
	}
	if err := r.Replace(2, nil); err == nil {
		t.Fatal("ожидалась ошибка для nil при замене")
	}
}

func TestSortAndPagination(t *testing.T) {
	r := newRepository(t)
	for _, name := range []string{"Сидоров", "Иванов", "Петров", "Абрамов", "Яковлев"} {
		if _, err := r.Add(newDriver(t, name)); err != nil {
			t.Fatal(err)
		}
	}
	r.SortByLastName()
	page, err := r.GetKNShortList(2, 2)
	if err != nil || len(page) != 2 {
		t.Fatalf("вторая страница: %v, %v", page, err)
	}
	if page[0].LastName() != "Петров" || page[1].LastName() != "Сидоров" || page[0].Initials() != "И.П." {
		t.Fatalf("неверный порядок: %v", page)
	}
	if err := page[0].SetLastName("Другой"); err != nil {
		t.Fatal(err)
	}
	got, err := r.GetByID(page[0].ID())
	if err != nil || got.LastName() != "Петров" {
		t.Fatalf("страница не должна менять хранимую модель: %v, %v", got, err)
	}
	for _, tc := range []struct{ k, n, count int }{
		{1, 2, 2}, {3, 2, 1}, {4, 2, 0}, {1, 100, 5},
		{int(^uint(0) >> 1), int(^uint(0) >> 1), 0},
	} {
		page, err := r.GetKNShortList(tc.k, tc.n)
		if err != nil || len(page) != tc.count {
			t.Fatalf("страница %+v: len=%d, err=%v", tc, len(page), err)
		}
	}
	for _, args := range [][2]int{{0, 1}, {1, 0}, {-1, 2}, {1, -2}} {
		if _, err := r.GetKNShortList(args[0], args[1]); err == nil {
			t.Fatalf("ожидалась ошибка для %v", args)
		}
	}
}

func TestJSONRoundTrip(t *testing.T) {
	r := newRepository(t)
	if err := r.WriteAll(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(r.path)
	if err != nil || string(data) != "[]\n" {
		t.Fatalf("пустой список: %q, %v", data, err)
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
	loaded, err := NewDriverRepJSON(r.path)
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
	if loaded.GetCount() != 2 || got.LastName() != d.LastName() ||
		got.FirstName() != d.FirstName() || got.MiddleName() != d.MiddleName() ||
		got.PhoneNumber() != d.PhoneNumber() || got.Experience() != d.Experience() ||
		got.Initials() != "И.П." {
		t.Fatalf("поля потеряны при сохранении: %v", got)
	}
	if id, err := loaded.Add(d); err != nil || id != 3 {
		t.Fatalf("ID после загрузки: %d, %v", id, err)
	}
}

func TestReadErrorsPreserveState(t *testing.T) {
	r := newRepository(t)
	if err := r.ReadAll(); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("отсутствующий файл: %v", err)
	}
	if _, err := r.Add(newDriver(t, "Иванов")); err != nil {
		t.Fatal(err)
	}
	valid := `{"id":8,"last_name":"Петров","first_name":"Иван","middle_name":"Петрович","phone_number":"+79161234567","experience":10}`
	for _, data := range []string{"", "null", "{}", "[", "[{}]", "[" + valid + "," + valid + "]", "[" + valid + ",{}]"} {
		if err := os.WriteFile(r.path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		if err := r.ReadAll(); err == nil {
			t.Fatalf("ожидалась ошибка для %q", data)
		}
		if r.GetCount() != 1 {
			t.Fatal("ошибочная загрузка изменила список")
		}
		if _, err := r.GetByID(1); err != nil {
			t.Fatal(err)
		}
	}
	if id, err := r.Add(newDriver(t, "Сидоров")); err != nil || id != 2 {
		t.Fatalf("ошибочная загрузка изменила счётчик ID: %d, %v", id, err)
	}
	if err := os.WriteFile(r.path, []byte("["+valid+"]"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := r.ReadAll(); err != nil {
		t.Fatal(err)
	}
	if id, err := r.Add(newDriver(t, "Сидоров")); err != nil || id != 9 {
		t.Fatalf("новый ID должен быть больше максимального загруженного: %d, %v", id, err)
	}
}

func TestInvalidPath(t *testing.T) {
	if _, err := NewDriverRepJSON(" "); err == nil {
		t.Fatal("ожидалась ошибка для пустого пути")
	}
	r, err := NewDriverRepJSON(filepath.Join(t.TempDir(), "missing", "drivers.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := r.WriteAll(); err == nil {
		t.Fatal("запись в отсутствующую папку должна вернуть ошибку")
	}
}
