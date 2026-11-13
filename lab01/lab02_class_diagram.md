```mermaid
classDiagram
    class Driver {
        +DriverShort DriverShort
        -string firstName
        -string middleName
        -int experience
        -string phoneNumber
        +FirstName() string
        +MiddleName() string
        +Experience() int
        +PhoneNumber() string
    }

    class DriverShort {
        -int id
        -string lastName
        -string initials
        +ID() int
        +LastName() string
        +Initials() string
    }

    class DriverMapper {
        +ToShort(d *Driver) *DriverShort
    }

    class DriverRepJSON {
        -string path
        -drivers []*Driver
        -int lastID
        -mapper *DriverMapper
        +ReadAll() error
        +WriteAll() error
        +GetByID(id int) (*Driver, error)
        +GetKNShortList(k int, n int) ([]*DriverShort, error)
        +SortByLastName()
        +Add(d *Driver) (int, error)
        +Replace(id int, d *Driver) error
        +Delete(id int) error
        +GetCount() int
    }

    class driverRecord {
        +int ID
        +string LastName
        +string FirstName
        +string MiddleName
        +string PhoneNumber
        +int Experience
    }

    Driver "1" *-- "1" DriverShort : встраивает
    DriverRepJSON "1" *-- "0..*" Driver : хранит копии
    DriverRepJSON "1" *-- "1" DriverMapper : использует
    DriverRepJSON ..> driverRecord : читает и записывает JSON
    DriverRepJSON ..> DriverShort : возвращает страницы
    DriverMapper ..> Driver : принимает
    DriverMapper ..> DriverShort : возвращает указатель
```

Диаграмма отражает первый пункт ЛР2. Методы моделей первой лабораторной
показаны в сокращённом виде; полная диаграмма ЛР1 находится в
[driver_class_diagram.md](driver_class_diagram.md).

`NewDriverRepJSON(path)` создаёт пустой репозиторий. Это функция пакета
`file`, возвращающая `(*DriverRepJSON, error)`.

`driverRecord` содержит экспортируемые поля с JSON-тегами. При загрузке
каждая запись превращается в `Driver` через конструктор с валидацией.

Маппер возвращает указатель на встроенную модель. Репозиторий копирует
полученный `DriverShort`, чтобы изменение результата не затрагивало его список.
