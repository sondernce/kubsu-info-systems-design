package driver

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var nameRegex = regexp.MustCompile(`^[А-Яа-яЁёA-Za-z-]+$`)

type Driver struct {
	id         int
	lastName   string
	firstName  string
	middleName string
	experience int
}

// валидаторы

func ValidateID(id int) error {
	if id <= 0 {
		return errors.New("ID должен быть больше 0")
	}
	return nil
}

// универсальный валидатор для first/middle/last name 
func ValidateName(value, fieldName string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fmt.Errorf("%s не может быть пустым", fieldName)
	}
	if !nameRegex.MatchString(trimmed) {
		return fmt.Errorf("%s должно содержать только буквы", fieldName)
	}
	return nil
}

func ValidateExperience(exp int) error {
	if exp < 0 || exp > 70 {
		return errors.New("стаж должен быть в диапазоне от 0 до 70 лет")
	}
	return nil
}

func ValidateAll(id int, lastName, firstName, middleName string, exp int) error {
	if err := ValidateID(id); err != nil {
		return err
	}
	if err := ValidateName(lastName, "Фамилия"); err != nil {
		return err
	}
	if err := ValidateName(firstName, "Имя"); err != nil {
		return err
	}
	if err := ValidateName(middleName, "Отчество"); err != nil {
		return err
	}
	if err := ValidateExperience(exp); err != nil {
		return err
	}
	return nil
}

// констрктор

func NewDriver(id int, lastName, firstName, middleName string, exp int) (*Driver, error) {
	if err := ValidateAll(id, lastName, firstName, middleName, exp); err != nil {
		return nil, err
	}
	return &Driver{id, lastName, firstName, middleName, exp}, nil
}

// геттеры
func (d *Driver) ID() int          { return d.id }
func (d *Driver) LastName() string   { return d.lastName }
func (d *Driver) FirstName() string  { return d.firstName }
func (d *Driver) MiddleName() string { return d.middleName }
func (d *Driver) Experience() int   { return d.experience }