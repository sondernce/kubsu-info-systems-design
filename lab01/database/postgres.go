package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Database выполняет запросы через общий пул PostgreSQL.
type Database struct {
	pool *sql.DB
}

var (
	mu       sync.Mutex
	instance *Database
	dsnInUse string
	closed   bool
)

// GetDatabase возвращает первый успешно открытый экземпляр на весь процесс.
func GetDatabase(ctx context.Context, dsn string) (*Database, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("строка подключения PostgreSQL пуста")
	}

	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		if dsn != dsnInUse {
			return nil, errors.New("PostgreSQL уже настроен с другой строкой подключения")
		}
		if closed {
			return nil, errors.New("пул PostgreSQL закрыт")
		}
		return instance, nil
	}

	pool, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("открытие пула PostgreSQL: %w", err)
	}
	if err := pool.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("подключение к PostgreSQL: %w", errors.Join(err, pool.Close()))
	}

	instance = &Database{pool: pool}
	dsnInUse = dsn
	return instance, nil
}

func (d *Database) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.pool.QueryRowContext(ctx, query, args...)
}

func (d *Database) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.pool.QueryContext(ctx, query, args...)
}

func (d *Database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.pool.ExecContext(ctx, query, args...)
}

// Close закрывает общий пул при завершении программы.
func (d *Database) Close() error {
	mu.Lock()
	defer mu.Unlock()
	closed = true
	return d.pool.Close()
}
