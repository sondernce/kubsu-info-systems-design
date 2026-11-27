package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"transportation/database"
	"transportation/driver"
)

var ErrDriverNotFound = errors.New("водитель не найден")

type DriverRepDB struct {
	db *database.Database
}

func NewDriverRepDB(db *database.Database) *DriverRepDB {
	return &DriverRepDB{db: db}
}

//a

func (r *DriverRepDB) GetByID(ctx context.Context, id int) (*driver.Driver, error) {
	const query = `
		SELECT last_name, first_name, middle_name, phone_number, experience
		FROM drivers WHERE id = $1`
	var lastName, firstName, middleName, phoneNumber string
	var experience int
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&lastName, &firstName, &middleName, &phoneNumber, &experience,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %d", ErrDriverNotFound, id)
	}
	if err != nil {
		return nil, fmt.Errorf("получение водителя %d: %w", id, err)
	}
	d, err := driver.NewDriver(id, lastName, firstName, middleName, phoneNumber, experience)
	if err != nil {
		return nil, fmt.Errorf("данные водителя %d: %w", id, err)
	}
	return d, nil
}

//b

func (r *DriverRepDB) GetKNShortList(ctx context.Context, k, n int) ([]*driver.DriverShort, error) {
	if k <= 0 || n <= 0 {
		return nil, errors.New("номер страницы и её размер должны быть больше нуля")
	}
	if int64(k-1) > math.MaxInt64/int64(n) {
		return nil, errors.New("смещение страницы слишком велико")
	}
	offset := int64(k-1) * int64(n)

	const query = `
		SELECT id, last_name, first_name, middle_name
		FROM drivers ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, n, offset)
	if err != nil {
		return nil, fmt.Errorf("получение страницы водителей: %w", err)
	}
	defer rows.Close()

	result := make([]*driver.DriverShort, 0)
	for rows.Next() {
		var id int
		var lastName, firstName, middleName string
		if err := rows.Scan(&id, &lastName, &firstName, &middleName); err != nil {
			return nil, fmt.Errorf("чтение строки водителя: %w", err)
		}
		short, err := driver.NewDriverShort(id, lastName, firstName, middleName)
		if err != nil {
			return nil, fmt.Errorf("данные водителя %d: %w", id, err)
		}
		result = append(result, short)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("чтение страницы водителей: %w", err)
	}
	return result, nil
}

//c

func (r *DriverRepDB) Add(ctx context.Context, d *driver.Driver) (int, error) {
	if err := validateDriver(d); err != nil {
		return 0, err
	}
	const query = `
		INSERT INTO drivers (last_name, first_name, middle_name, phone_number, experience)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query,
		d.LastName(), d.FirstName(), d.MiddleName(), d.PhoneNumber(), d.Experience(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("добавление водителя: %w", err)
	}
	return id, nil
}

//d

func (r *DriverRepDB) Replace(ctx context.Context, id int, d *driver.Driver) error {
	if err := validateDriver(d); err != nil {
		return err
	}
	const query = `
		UPDATE drivers
		SET last_name = $1, first_name = $2, middle_name = $3, phone_number = $4, experience = $5
		WHERE id = $6`
	result, err := r.db.ExecContext(ctx, query,
		d.LastName(), d.FirstName(), d.MiddleName(), d.PhoneNumber(), d.Experience(), id,
	)
	if err != nil {
		return fmt.Errorf("замена водителя %d: %w", id, err)
	}
	return checkDriverAffected(result, id)
}

//e

func (r *DriverRepDB) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM drivers WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("удаление водителя %d: %w", id, err)
	}
	return checkDriverAffected(result, id)
}

//f

func (r *DriverRepDB) GetCount(ctx context.Context) (int, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM drivers").Scan(&count); err != nil {
		return 0, fmt.Errorf("подсчёт водителей: %w", err)
	}
	return count, nil
}

func validateDriver(d *driver.Driver) error {
	if d == nil {
		return errors.New("водитель не может быть nil")
	}
	_, err := driver.NewDriver(d.ID(), d.LastName(), d.FirstName(), d.MiddleName(), d.PhoneNumber(), d.Experience())
	return err
}

func checkDriverAffected(result sql.Result, id int) error {
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка результата для водителя %d: %w", id, err)
	}
	if count == 0 {
		return fmt.Errorf("%w: %d", ErrDriverNotFound, id)
	}
	return nil
}
