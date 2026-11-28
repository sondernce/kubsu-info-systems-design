```mermaid
classDiagram
    class Printer {
        <<interface>>
        +Print() string
    }

    class DriverShort {
        -int id
        -string lastName
        -string initials
        +NewDriverShort(id, lastName, firstName, middleName) (*DriverShort, error)
        +ID() int
        +LastName() string
        +Initials() string
        +SetLastName(lastName) error
        +ShortPrint() string
        +Print() string
    }

    class Driver {
        +DriverShort DriverShort
        -string firstName
        -string middleName
        -int experience
        -string phoneNumber
        +NewDriver(args ...interface) (*Driver, error)
        -newDriverFull(id, lastName, firstName, middleName, phoneNumber, exp) (*Driver, error)
        -newDriverFromString(s) (*Driver, error)
        -newDriverFromJSON(jsonData) (*Driver, error)
        +FirstName() string
        +MiddleName() string
        +Experience() int
        +PhoneNumber() string
        +SetFirstName(firstName) error
        +SetMiddleName(middleName) error
        +SetExperience(exp) error
        +SetPhoneNumber(phoneNumber) error
        -recalcInitials()
        +FullPrint() string
        +Print() string
        +Equals(other) bool
    }

    class DriverMapper {
        +NewDriverMapper() *DriverMapper
        +ToShort(d *Driver) *DriverShort
    }

    Driver *-- DriverShort : embeds (composition)
    DriverShort ..|> Printer : realizes
    Driver ..|> Printer : realizes (переопределяет)
    DriverMapper ..> Driver : принимает
    DriverMapper ..> DriverShort : возвращает указатель
```

`DriverMapper.ToShort(d)` возвращает `&d.DriverShort` — указатель на встроенную
краткую модель, а не копию. Поэтому изменение имени или отчества через сеттеры
`Driver` обновляет инициалы и в ранее полученном `DriverShort`.

`DriverShort` встроен в `Driver` как экспортируемое поле без отдельного имени
`short`. Преобразование `ToShort` принадлежит мапперу, а не `Driver`.

Функции создания (`NewDriver`, `NewDriverShort`, `NewDriverMapper` и
`newDriver…`) показаны рядом с соответствующими типами для удобства чтения;
в Go это функции пакета `driver`, а не методы структур.
