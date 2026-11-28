package postgres

import (
	"errors"
	"fmt"
	"strings"

	"transportation/driver"
	"transportation/repository"
)

// DriverFilter ограничивает выборку по точной фамилии и минимальному стажу.
type DriverFilter struct {
	LastName      string
	MinExperience int
}

type SortField uint8

const (
	SortByID SortField = iota
	SortByLastName
	SortByExperience
)

type DriverSort struct {
	Field      SortField
	Descending bool
}

var _ repository.DriverRepository = (*DriverDBDecorator)(nil)

// DriverDBDecorator применяет фильтр к списку и счётчику и задаёт порядок списка.
type DriverDBDecorator struct {
	*DriverDBAdapter
	filter  DriverFilter
	orderBy string
}

func NewDriverDBDecorator(adapter *DriverDBAdapter, filter DriverFilter, sorting DriverSort) (*DriverDBDecorator, error) {
	if adapter == nil {
		return nil, errors.New("адаптер PostgreSQL не задан")
	}
	if err := driver.ValidateExperience(filter.MinExperience); err != nil {
		return nil, err
	}
	filter.LastName = strings.TrimSpace(filter.LastName)
	orderBy, err := driverOrderBy(sorting)
	if err != nil {
		return nil, err
	}
	return &DriverDBDecorator{DriverDBAdapter: adapter, filter: filter, orderBy: orderBy}, nil
}

func (d *DriverDBDecorator) GetKNShortList(k, n int) ([]*driver.DriverShort, error) {
	return d.repo.getKNShortList(d.ctx, k, n, d.filter, d.orderBy)
}

func (d *DriverDBDecorator) GetCount() (int, error) {
	return d.repo.getCount(d.ctx, d.filter)
}

func driverOrderBy(sorting DriverSort) (string, error) {
	var column string
	switch sorting.Field {
	case SortByID:
		column = "id"
	case SortByLastName:
		column = "last_name"
	case SortByExperience:
		column = "experience"
	default:
		return "", fmt.Errorf("неизвестное поле сортировки: %d", sorting.Field)
	}

	direction := " ASC"
	if sorting.Descending {
		direction = " DESC"
	}
	orderBy := column + direction
	if sorting.Field != SortByID {
		orderBy += ", id ASC"
	}
	return orderBy, nil
}
