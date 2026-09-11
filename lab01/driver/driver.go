package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var nameRegex = regexp.MustCompile(`^[А-Яа-яЁёA-Za-z-]+$`)

type Driver struct {
	DriverShort            // встраивание
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

// func ValidateAll(id int, lastName, firstName, middleName string, exp int) error {
// 	if err := ValidateID(id); err != nil {
// 		return err
// 	}
// 	if err := ValidateName(lastName, "Фамилия"); err != nil {
// 		return err
// 	}
// 	if err := ValidateName(firstName, "Имя"); err != nil {
// 		return err
// 	}
// 	if err := ValidateName(middleName, "Отчество"); err != nil {
// 		return err
// 	}
// 	if err := ValidateExperience(exp); err != nil {
// 		return err
// 	}
// 	return nil
// }

// констрктор

func NewDriver(id int, lastName, firstName, middleName string, exp int) (*Driver, error) {
	// создаем базовую 
	short, err := NewDriverShort(id, lastName, firstName, middleName)
	if err != nil {
		return nil, err
	}

	// валидируем специфичные для полного класса поля
	if err := ValidateExperience(exp); err != nil {
		return nil, err
	}

	return &Driver{
		DriverShort: *short,
		firstName:  firstName,
		middleName: middleName,
		experience: exp,
	}, nil
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
func (d *Driver) FirstName() string  { return d.firstName }
func (d *Driver) MiddleName() string { return d.middleName }
func (d *Driver) Experience() int    { return d.experience }

// 7
// полная версия вывода
func (d *Driver) FullPrint() string {
	return fmt.Sprintf("Водитель [#%d]: %s %s %s, стаж: %d лет", d.id, d.lastName, d.firstName, d.middleName, d.experience)
}

// краткая версия вывода по инициалам
func (d *Driver) ShortPrint() string {
	fRunes := []rune(d.firstName)
	mRunes := []rune(d.middleName)
	return fmt.Sprintf("ID: %d | %s %c.%c.", d.id, d.lastName, fRunes[0], mRunes[0])
}

// сравнение объектов на равенство по содержимому полей
func (d *Driver) Equals(other *Driver) bool {
	if other == nil {
		return false
	}
	return d.id == other.id &&
		d.lastName == other.lastName &&
		d.firstName == other.firstName &&
		d.middleName == other.middleName &&
		d.experience == other.experience
}

//8
//----------------------------------------------------------------------------

type DriverShort struct {
	id       int
	lastName string
	initials string
}

func NewDriverShort(id int, lastName, firstName, middleName string) (*DriverShort, error) {
	if err := ValidateID(id); err != nil {
		return nil, err
	}
	if err := ValidateName(lastName, "Фамилия"); err != nil {
		return nil, err
	}
	if err := ValidateName(firstName, "Имя"); err != nil {
		return nil, err
	}
	if err := ValidateName(middleName, "Отчество"); err != nil {
		return nil, err
	}

	fRunes := []rune(firstName)
	mRunes := []rune(middleName)
	initials := fmt.Sprintf("%c.%c.", fRunes[0], mRunes[0])

	return &DriverShort{
		id:       id,
		lastName: lastName,
		initials: initials,
	}, nil
}

func (ds *DriverShort) ID() int          { return ds.id }
func (ds *DriverShort) LastName() string { return ds.lastName }
func (ds *DriverShort) Initials() string { return ds.initials }

func (ds *DriverShort) FullPrint() string {
	return fmt.Sprintf("ID: %d | %s %s", ds.id, ds.lastName, ds.initials)
}
