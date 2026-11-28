package repository

import (
	"errors"

	"transportation/driver"
)

var ErrDriverNotFound = errors.New("водитель не найден")

// DriverRepository определяет операции, доступные во всех хранилищах водителей.
type DriverRepository interface {
	GetByID(id int) (*driver.Driver, error)
	GetKNShortList(k, n int) ([]*driver.DriverShort, error)
	Add(d *driver.Driver) (int, error)
	Replace(id int, d *driver.Driver) error
	Delete(id int) error
	GetCount() (int, error)
}
