package driver

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"encoding/json"
)

var nameRegex = regexp.MustCompile(`^[А-Яа-яЁёA-Za-z-]+$`)

type Driver struct {
	id         int
	lastName   string
	firstName  string
	middleName string
	experience int
}

// дтошка для джейсона
type driverDTO struct {
	ID         int    `json:"id"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Experience int    `json:"experience"`
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

// в гошке нет перегрузки в привычном виде поэтому вот так
// перегрузка 1 - создание из строки вида id;lastName;firstName;middleName;experience
func NewDriverFromString(s string) (*Driver, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 5 {
		return nil, errors.New("неверный формат строки, ожидается 5 элементов через ';'")
	}

	id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, errors.New("некорректный ID")
	}

	exp, err := strconv.Atoi(strings.TrimSpace(parts[4]))
	if err != nil {
		return nil, errors.New("некорректный стаж")
	}

	return NewDriver(id, parts[1], parts[2], parts[3], exp)
}
// перегрузка 2 создание из джейсончика
func NewDriverFromJSON(jsonData []byte) (*Driver, error) {
	var dto driverDTO
	if err := json.Unmarshal(jsonData, &dto); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	return NewDriver(dto.ID, dto.LastName, dto.FirstName, dto.MiddleName, dto.Experience)
}

// геттеры
func (d *Driver) ID() int          { return d.id }
func (d *Driver) LastName() string   { return d.lastName }
func (d *Driver) FirstName() string  { return d.firstName }
func (d *Driver) MiddleName() string { return d.middleName }
func (d *Driver) Experience() int   { return d.experience }