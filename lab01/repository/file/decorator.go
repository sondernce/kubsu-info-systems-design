package file

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"transportation/driver"
	"transportation/repository"
)

type fileDriverRepository interface {
	repository.DriverRepository
	ReadAll() error
	WriteAll() error
	SortByLastName()
	fileBase() *driverFileRepository
}

var _ repository.DriverRepository = (*DriverFileDecorator)(nil)

// DriverFileDecorator задаёт фильтр и порядок выдачи для JSON и YAML.
type DriverFileDecorator struct {
	fileDriverRepository
	filter  repository.DriverFilter
	sorting repository.DriverSort
}

func NewDriverFileDecorator(repo fileDriverRepository, filter repository.DriverFilter, sorting repository.DriverSort) (*DriverFileDecorator, error) {
	if err := driver.ValidateExperience(filter.MinExperience); err != nil {
		return nil, err
	}
	switch sorting.Field {
	case repository.SortByID, repository.SortByLastName, repository.SortByExperience:
	default:
		return nil, fmt.Errorf("неизвестное поле сортировки: %d", sorting.Field)
	}
	filter.LastName = strings.TrimSpace(filter.LastName)
	return &DriverFileDecorator{fileDriverRepository: repo, filter: filter, sorting: sorting}, nil
}

func (d *DriverFileDecorator) GetKNShortList(k, n int) ([]*driver.DriverShort, error) {
	base := d.fileBase()
	var matches []*driver.Driver
	for _, item := range base.drivers {
		if d.matches(item) {
			matches = append(matches, item)
		}
	}
	slices.SortFunc(matches, d.compare)
	return base.shortPage(matches, k, n)
}

func (d *DriverFileDecorator) GetCount() (int, error) {
	count := 0
	for _, item := range d.fileBase().drivers {
		if d.matches(item) {
			count++
		}
	}
	return count, nil
}

func (d *DriverFileDecorator) matches(item *driver.Driver) bool {
	return item.Experience() >= d.filter.MinExperience &&
		(d.filter.LastName == "" || item.LastName() == d.filter.LastName)
}

func (d *DriverFileDecorator) compare(a, b *driver.Driver) int {
	var order int
	switch d.sorting.Field {
	case repository.SortByID:
		order = cmp.Compare(a.ID(), b.ID())
	case repository.SortByLastName:
		order = cmp.Compare(a.LastName(), b.LastName())
	case repository.SortByExperience:
		order = cmp.Compare(a.Experience(), b.Experience())
	}
	if d.sorting.Descending {
		order = -order
	}
	if order == 0 {
		return cmp.Compare(a.ID(), b.ID())
	}
	return order
}
