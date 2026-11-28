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

    class driverFileRepository {
        -string path
        -drivers []*Driver
        -int lastID
        -mapper *DriverMapper
        -readAll(decode) error
        -writeAll(encode) error
        +GetByID(id int) (*Driver, error)
        +GetKNShortList(k int, n int) ([]*DriverShort, error)
        +SortByLastName()
        +Add(d *Driver) (int, error)
        +Replace(id int, d *Driver) error
        +Delete(id int) error
        +GetCount() (int, error)
    }

    class DriverRepJSON {
        -driverFileRepository driverFileRepository
        +ReadAll() error
        +WriteAll() error
    }

    class DriverRepYAML {
        -driverFileRepository driverFileRepository
        +ReadAll() error
        +WriteAll() error
    }

    class DriverRepository {
        <<interface>>
        +GetByID(id int) (*Driver, error)
        +GetKNShortList(k int, n int) ([]*DriverShort, error)
        +Add(d *Driver) (int, error)
        +Replace(id int, d *Driver) error
        +Delete(id int) error
        +GetCount() (int, error)
    }

    class DriverDBAdapter {
        -ctx context.Context
        -repo *DriverRepDB
        +GetByID(id int) (*Driver, error)
        +GetKNShortList(k int, n int) ([]*DriverShort, error)
        +Add(d *Driver) (int, error)
        +Replace(id int, d *Driver) error
        +Delete(id int) error
        +GetCount() (int, error)
    }

    class DriverDBDecorator {
        -filter DriverFilter
        -orderBy string
        +GetKNShortList(k int, n int) ([]*DriverShort, error)
        +GetCount() (int, error)
    }

    class DriverFilter {
        +string LastName
        +int MinExperience
    }

    class DriverSort {
        +SortField Field
        +bool Descending
    }

    class DriverRepDB {
        -db *Database
        +GetByID(ctx, id int) (*Driver, error)
        +GetKNShortList(ctx, k int, n int) ([]*DriverShort, error)
        +Add(ctx, d *Driver) (int, error)
        +Replace(ctx, id int, d *Driver) error
        +Delete(ctx, id int) error
        +GetCount(ctx) (int, error)
    }

    class Database {
        <<Singleton>>
        -pool *sql.DB
        +QueryRowContext(ctx, query, args) *sql.Row
        +QueryContext(ctx, query, args) (*sql.Rows, error)
        +ExecContext(ctx, query, args) (sql.Result, error)
        +Close() error
    }

    class SQLDB["sql.DB"]

    class driverRecord {
        +int ID
        +string LastName
        +string FirstName
        +string MiddleName
        +string PhoneNumber
        +int Experience
    }

    Driver "1" *-- "1" DriverShort : встраивает
    DriverRepJSON *-- driverFileRepository : встраивает
    DriverRepYAML *-- driverFileRepository : встраивает
    driverFileRepository "1" *-- "0..*" Driver : хранит копии
    driverFileRepository "1" *-- "1" DriverMapper : использует
    driverFileRepository ..> driverRecord : преобразует записи
    driverFileRepository ..> DriverShort : возвращает страницы
    DriverRepJSON ..|> DriverRepository
    DriverRepYAML ..|> DriverRepository
    DriverDBAdapter ..|> DriverRepository
    DriverDBDecorator ..|> DriverRepository
    DriverDBDecorator --> DriverDBAdapter : оборачивает
    DriverDBDecorator --> DriverFilter : хранит фильтр
    DriverDBDecorator ..> DriverSort : принимает при создании
    DriverDBDecorator ..> DriverRepDB : запрашивает выборку с фильтром
    DriverDBAdapter --> DriverRepDB : передаёт контекст и вызовы
    DriverRepDB --> Database : делегирует запросы
    Database *-- SQLDB : управляет пулом
    DriverRepDB ..> Driver : читает и сохраняет
    DriverRepDB ..> DriverShort : возвращает страницы
    DriverMapper ..> Driver : принимает
    DriverMapper ..> DriverShort : возвращает указатель
```

Диаграмма отражает пункты с 1 по 7 ЛР2. Методы моделей первой лабораторной
показаны в сокращённом виде; полная диаграмма ЛР1 находится в
[driver_class_diagram.md](driver_class_diagram.md).

`NewDriverRepJSON(path)` и `NewDriverRepYAML(path)` создают пустые репозитории.
Это функции пакета `file`, возвращающие указатель на соответствующий тип
и ошибку. Для загрузки файла нужно вызвать `ReadAll`.

`DriverRepJSON` и `DriverRepYAML` встраивают `driverFileRepository`. Методы
поиска, сортировки, пагинации и изменения списка определены в общем типе и
доступны у обоих репозиториев. Это аналог наследования поведения через
встраивание в Go. На диаграмме показана композиция.

Общий тип читает и записывает файл. Каждый репозиторий передаёт ему функции
преобразования `driverRecord` в свой формат и обратно. При загрузке общий
тип создаёт `Driver` через конструктор с валидацией.

Маппер возвращает указатель на встроенную модель. Репозиторий копирует
полученный `DriverShort`, чтобы изменение результата не затрагивало его список.

`database.GetDatabase(ctx, dsn)` открывает общий пул PostgreSQL и возвращает
один экземпляр `Database` на процесс. `NewDriverRepDB(db)` принимает этот
экземпляр. Репозиторий составляет SQL и преобразует результаты в модели.
Выполнение запросов он делегирует `Database`.

`DriverRepJSON` и `DriverRepYAML` реализуют `DriverRepository` через методы
встроенного типа. `DriverDBAdapter` реализует тот же интерфейс и передаёт
контекст из конструктора в каждый вызов `DriverRepDB`.

`DriverDBDecorator` оборачивает `DriverDBAdapter` и меняет поведение
`GetKNShortList` и `GetCount`. Остальные вызовы выполняет адаптер.
PostgreSQL применяет фильтр и сортировку перед выбором страницы.
Счётчик использует тот же фильтр.
