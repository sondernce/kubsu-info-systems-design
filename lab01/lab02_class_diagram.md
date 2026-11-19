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

    class DriverRepYAML {
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
    DriverRepYAML "1" *-- "0..*" Driver : хранит копии
    DriverRepYAML "1" *-- "1" DriverMapper : использует
    DriverRepYAML ..> driverRecord : читает и записывает YAML
    DriverRepYAML ..> DriverShort : возвращает страницы
    DriverMapper ..> Driver : принимает
    DriverMapper ..> DriverShort : возвращает указатель
```

Диаграмма отражает пункты 1 и 2 ЛР2. Методы моделей первой лабораторной
показаны в сокращённом виде; полная диаграмма ЛР1 находится в
[driver_class_diagram.md](driver_class_diagram.md).

`NewDriverRepJSON(path)` и `NewDriverRepYAML(path)` создают пустые репозитории.
Это функции пакета `file`, возвращающие указатель на соответствующий тип
и ошибку. Для загрузки файла нужно вызвать `ReadAll`.

`driverRecord` содержит экспортируемые поля с JSON- и YAML-тегами. Оба
репозитория используют эту структуру из `repository/file/record.go`. При загрузке
каждая запись превращается в `Driver` через конструктор с валидацией.

На этом этапе репозитории являются отдельными типами. Объединение в общую
иерархию относится к пункту 3 ЛР2.

Маппер возвращает указатель на встроенную модель. Репозиторий копирует
полученный `DriverShort`, чтобы изменение результата не затрагивало его список.
