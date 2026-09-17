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
var phoneRegex = regexp.MustCompile(`^\+7\d{10}$`)

// явная защита от SQL-инъекций (кавычки, разделители команд, ключевые слова)
var sqlInjectionPattern = regexp.MustCompile(`(?i)('|"|;|--|/\*|\*/|union\s+select|drop\s+table|insert\s+into|delete\s+from)`)

func containsForbiddenPattern(value string) bool {
	return sqlInjectionPattern.MatchString(value)
}

type Driver struct {
	DriverShort // встраивание
	firstName   string
	middleName  string
	experience  int
	phoneNumber string
}

// дтошка для джейсона
type driverDTO struct {
	ID          int    `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	Experience  int    `json:"experience"`
	PhoneNumber string `json:"phone_number"`
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
	if containsForbiddenPattern(trimmed) {
		return fmt.Errorf("%s содержит недопустимые символы", fieldName)
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

func ValidatePhoneNumber(phone string) error {
	trimmed := strings.TrimSpace(phone)
	if trimmed == "" {
		return errors.New("номер телефона не может быть пустым")
	}
	if containsForbiddenPattern(trimmed) {
		return errors.New("номер телефона содержит недопустимые символы")
	}
	if !phoneRegex.MatchString(trimmed) {
		return errors.New("номер телефона должен быть в формате +7XXXXXXXXXX")
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
// в гошке нет перегрузки в привычном виде, поэтому эмулируем её одним именем NewDriver:
// диспетчеризация по количеству и типам аргументов через interface{}

// NewDriver - "перегруженный" конструктор:
//   - NewDriver(id int, lastName, firstName, middleName, phoneNumber string, exp int) - обычные поля
//   - NewDriver(s string)                                                              - строка "id;lastName;firstName;middleName;phoneNumber;experience"
//   - NewDriver(jsonData []byte)                                                       - JSON
func NewDriver(args ...interface{}) (*Driver, error) {
	switch len(args) {
	case 1:
		switch v := args[0].(type) {
		case string:
			return newDriverFromString(v)
		case []byte:
			return newDriverFromJSON(v)
		default:
			return nil, errors.New("неверный тип аргумента: ожидается string или []byte")
		}
	case 6:
		id, ok := args[0].(int)
		if !ok {
			return nil, errors.New("аргумент 1 (id) должен быть int")
		}
		lastName, ok := args[1].(string)
		if !ok {
			return nil, errors.New("аргумент 2 (lastName) должен быть string")
		}
		firstName, ok := args[2].(string)
		if !ok {
			return nil, errors.New("аргумент 3 (firstName) должен быть string")
		}
		middleName, ok := args[3].(string)
		if !ok {
			return nil, errors.New("аргумент 4 (middleName) должен быть string")
		}
		phoneNumber, ok := args[4].(string)
		if !ok {
			return nil, errors.New("аргумент 5 (phoneNumber) должен быть string")
		}
		exp, ok := args[5].(int)
		if !ok {
			return nil, errors.New("аргумент 6 (experience) должен быть int")
		}
		return newDriverFull(id, lastName, firstName, middleName, phoneNumber, exp)
	default:
		return nil, errors.New("неверное количество аргументов для NewDriver")
	}
}

// создание из явных полей
func newDriverFull(id int, lastName, firstName, middleName, phoneNumber string, exp int) (*Driver, error) {
	// создаем базовую
	short, err := NewDriverShort(id, lastName, firstName, middleName)
	if err != nil {
		return nil, err
	}

	// валидируем специфичные для полного класса поля
	if err := ValidateExperience(exp); err != nil {
		return nil, err
	}
	trimmedPhone := strings.TrimSpace(phoneNumber)
	if err := ValidatePhoneNumber(trimmedPhone); err != nil {
		return nil, err
	}

	return &Driver{
		DriverShort: *short,
		firstName:   firstName,
		middleName:  middleName,
		experience:  exp,
		phoneNumber: trimmedPhone,
	}, nil
}

// создание из строки вида id;lastName;firstName;middleName;phoneNumber;experience
func newDriverFromString(s string) (*Driver, error) {
	parts := strings.Split(s, ";")
	if len(parts) != 6 {
		return nil, errors.New("неверный формат строки, ожидается 6 элементов через ';'")
	}

	id, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, errors.New("некорректный ID")
	}

	exp, err := strconv.Atoi(strings.TrimSpace(parts[5]))
	if err != nil {
		return nil, errors.New("некорректный стаж")
	}

	return newDriverFull(id, parts[1], parts[2], parts[3], parts[4], exp)
}

// создание из джейсончика
func newDriverFromJSON(jsonData []byte) (*Driver, error) {
	var dto driverDTO
	if err := json.Unmarshal(jsonData, &dto); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	return newDriverFull(dto.ID, dto.LastName, dto.FirstName, dto.MiddleName, dto.PhoneNumber, dto.Experience)
}

// геттеры
func (d *Driver) FirstName() string   { return d.firstName }
func (d *Driver) MiddleName() string  { return d.middleName }
func (d *Driver) Experience() int     { return d.experience }
func (d *Driver) PhoneNumber() string { return d.phoneNumber }

// сеттеры (с валидацией)
func (d *Driver) SetFirstName(firstName string) error {
	if err := ValidateName(firstName, "Имя"); err != nil {
		return err
	}
	d.firstName = firstName
	d.recalcInitials()
	return nil
}

func (d *Driver) SetMiddleName(middleName string) error {
	if err := ValidateName(middleName, "Отчество"); err != nil {
		return err
	}
	d.middleName = middleName
	d.recalcInitials()
	return nil
}

func (d *Driver) SetExperience(exp int) error {
	if err := ValidateExperience(exp); err != nil {
		return err
	}
	d.experience = exp
	return nil
}

func (d *Driver) SetPhoneNumber(phoneNumber string) error {
	trimmed := strings.TrimSpace(phoneNumber)
	if err := ValidatePhoneNumber(trimmed); err != nil {
		return err
	}
	d.phoneNumber = trimmed
	return nil
}

// пересчёт инициалов при смене имени/отчества, чтобы initials не рассинхронизировались
func (d *Driver) recalcInitials() {
	fRunes := []rune(d.firstName)
	mRunes := []rune(d.middleName)
	d.initials = fmt.Sprintf("%c.%c.", fRunes[0], mRunes[0])
}

// 7
// полная версия вывода
func (d *Driver) FullPrint() string {
	return fmt.Sprintf("Водитель [#%d]: %s %s %s, стаж: %d лет, тел.: %s", d.id, d.lastName, d.firstName, d.middleName, d.experience, d.phoneNumber)
}

// сравнение объектов на равенство по содержимому полей
func (d *Driver) Equals(other *Driver) bool {
	if other == nil {
		return false
	}
	return d.id == other.id &&
		d.phoneNumber == other.phoneNumber
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

func (ds *DriverShort) SetLastName(lastName string) error {
	if err := ValidateName(lastName, "Фамилия"); err != nil {
		return err
	}
	ds.lastName = lastName
	return nil
}

func (ds *DriverShort) ShortPrint() string {
	return fmt.Sprintf("ID: %d | %s %s", ds.id, ds.lastName, ds.initials)
}

// полиморфизм: единый интерфейс Printer с одинаковой сигнатурой Print() string,
// но разным поведением у DriverShort и Driver (Driver перекрывает встроенный метод)
type Printer interface {
	Print() string
}

func (ds *DriverShort) Print() string {
	return ds.ShortPrint()
}

func (d *Driver) Print() string {
	return d.FullPrint()
}