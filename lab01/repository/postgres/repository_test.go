package postgres

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"transportation/driver"
)

func TestDriverRepDB(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("задайте TEST_DATABASE_URL для проверки на PostgreSQL")
	}
	ctx := t.Context()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	// Временная таблица доступна только в соединении, которое её создало.
	db.SetMaxOpenConns(1)
	migration, err := os.ReadFile("../../migrations/001_create_drivers.sql")
	if err != nil {
		t.Fatal(err)
	}
	schema := strings.Replace(string(migration), "CREATE TABLE", "CREATE TEMP TABLE", 1)
	if _, err := db.ExecContext(ctx, schema); err != nil {
		t.Fatal(err)
	}
	repo := NewDriverRepDB(db)

	d, err := driver.NewDriver(99, "Петров", "Пётр", "Иванович", "+79161234567", 7)
	if err != nil {
		t.Fatal(err)
	}
	firstID, err := repo.Add(ctx, d)
	if err != nil || firstID <= 0 || firstID == d.ID() {
		t.Fatalf("Add: id=%d, err=%v", firstID, err)
	}
	secondID, err := repo.Add(ctx, d)
	if err != nil || secondID <= firstID {
		t.Fatalf("Add: id=%d, err=%v", secondID, err)
	}
	found, err := repo.GetByID(ctx, firstID)
	if err != nil || found.ID() != firstID || found.LastName() != d.LastName() ||
		found.FirstName() != d.FirstName() || found.MiddleName() != d.MiddleName() ||
		found.PhoneNumber() != d.PhoneNumber() || found.Experience() != d.Experience() {
		t.Fatalf("GetByID: %v, %v", found, err)
	}
	page, err := repo.GetKNShortList(ctx, 2, 1)
	if err != nil || len(page) != 1 || page[0].ID() != secondID || page[0].Initials() != "П.И." {
		t.Fatalf("GetKNShortList: %v, %v", page, err)
	}
	page, err = repo.GetKNShortList(ctx, 3, 1)
	if err != nil || len(page) != 0 {
		t.Fatalf("страница за пределами списка: %v, %v", page, err)
	}
	if _, err := repo.GetKNShortList(ctx, 0, 1); err == nil {
		t.Fatal("ожидалась ошибка для нулевой страницы")
	}

	if err := d.SetExperience(12); err != nil {
		t.Fatal(err)
	}
	if err := repo.Replace(ctx, firstID, d); err != nil {
		t.Fatal(err)
	}
	found, err = repo.GetByID(ctx, firstID)
	if err != nil || found.ID() != firstID || found.Experience() != 12 {
		t.Fatalf("Replace: %v, %v", found, err)
	}
	if err := repo.Delete(ctx, firstID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, firstID); !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("GetByID удалённой записи: %v", err)
	}
	if err := repo.Replace(ctx, firstID, d); !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("Replace удалённой записи: %v", err)
	}
	if err := repo.Delete(ctx, firstID); !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("повторный Delete: %v", err)
	}
	if count, err := repo.GetCount(ctx); err != nil || count != 1 {
		t.Fatalf("GetCount: %d, %v", count, err)
	}
}
