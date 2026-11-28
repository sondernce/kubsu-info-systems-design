package postgres

import (
	"context"

	"transportation/driver"
	"transportation/repository"
)

var _ repository.DriverRepository = (*DriverDBAdapter)(nil)

// DriverDBAdapter передаёт контекст в методы PostgreSQL.
type DriverDBAdapter struct {
	ctx  context.Context
	repo *DriverRepDB
}

// NewDriverDBAdapter связывает репозиторий с контекстом одной операции или запроса.
func NewDriverDBAdapter(ctx context.Context, repo *DriverRepDB) *DriverDBAdapter {
	return &DriverDBAdapter{ctx: ctx, repo: repo}
}

func (a *DriverDBAdapter) GetByID(id int) (*driver.Driver, error) {
	return a.repo.GetByID(a.ctx, id)
}

func (a *DriverDBAdapter) GetKNShortList(k, n int) ([]*driver.DriverShort, error) {
	return a.repo.GetKNShortList(a.ctx, k, n)
}

func (a *DriverDBAdapter) Add(d *driver.Driver) (int, error) {
	return a.repo.Add(a.ctx, d)
}

func (a *DriverDBAdapter) Replace(id int, d *driver.Driver) error {
	return a.repo.Replace(a.ctx, id, d)
}

func (a *DriverDBAdapter) Delete(id int) error {
	return a.repo.Delete(a.ctx, id)
}

func (a *DriverDBAdapter) GetCount() (int, error) {
	return a.repo.GetCount(a.ctx)
}
