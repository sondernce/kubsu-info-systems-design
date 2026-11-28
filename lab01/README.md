# Грузовые перевозки

Модуль `transportation` содержит модели первой лабораторной и файловые
репозитории водителей из ЛР2. Оба формата используют общую логику работы со списком.

Из папки `lab01`:

```sh
go run .
go test ./...
go vet ./...
```

`main.go` демонстрирует модели ЛР1. Проверки репозиториев находятся в
`repository/file/json_test.go` и `repository/file/yaml_test.go`.

## JSON-репозиторий

Тип `file.DriverRepJSON` находится в `repository/file/json.go`.
Чтение и запись отмечены комментариями `//a` и `//b` в файлах форматов,
операции со списком отмечены комментариями от `//c` до `//i` в `repository/file/base.go`.

| Подпункт | Метод | Поведение |
|---|---|---|
| a | `ReadAll` | Загружает список из файла с проверкой полей и уникальности ID |
| b | `WriteAll` | Записывает текущий список в JSON-файл |
| c | `GetByID` | Возвращает копию водителя по ID |
| d | `GetKNShortList` | Возвращает страницу кратких объектов |
| e | `SortByLastName` | Сортирует по фамилии по возрастанию, без учёта регистра |
| f | `Add` | Добавляет копию водителя и возвращает назначенный ID |
| g | `Replace` | Заменяет запись, сохраняя ID из аргумента метода |
| h | `Delete` | Удаляет запись по ID |
| i | `GetCount` | Возвращает число записей в памяти |

Пример использования:

```go
import (
    "transportation/driver"
    "transportation/repository/file"
)

func example() error {
    repo, err := file.NewDriverRepJSON("data/drivers.json")
    if err != nil {
        return err
    }
    if err := repo.ReadAll(); err != nil {
        return err
    }

    d, err := driver.NewDriver(1, "Сидоров", "Пётр", "Иванович", "+79161234569", 7)
    if err != nil {
        return err
    }
    id, err := repo.Add(d)
    if err != nil {
        return err
    }
    if _, err := repo.GetByID(id); err != nil {
        return err
    }
    repo.SortByLastName()
    if _, err := repo.GetKNShortList(1, 20); err != nil {
        return err
    }
    return repo.WriteAll()
}
```

Конструктор не читает файл. Для создания нового файла достаточно добавить
объекты и вызвать `WriteAll`. Папка для файла должна уже существовать.
Отсутствующий файл при `ReadAll` возвращает ошибку.

Изменения списка сохраняются только при вызове `WriteAll`. Запись выполняется
через временный файл в той же папке с последующим переименованием. Ошибка
разбора или валидации при `ReadAll` сохраняет прежний список в памяти.
Пустой список записывается как `[]`; `null` не принимается.

Нумерация страниц начинается с 1. `GetKNShortList(2, 20)` возвращает записи
с 21-й по 40-ю в текущем порядке списка. Неполная последняя страница содержит
оставшиеся записи, за пределами списка возвращается пустой срез. Нулевые и
отрицательные параметры возвращают ошибку.

При добавлении исходный ID игнорируется, входной объект не изменяется.
Новый ID больше всех загруженных и ранее выданных этим экземпляром репозитория.
После перезапуска счётчик восстанавливается по максимальному ID в файле;
история удалённых ID отдельно не хранится.

Возвращаемые модели и страницы не дают менять внутренний список напрямую.
Ошибку отсутствующего водителя можно проверить через
`errors.Is(err, repository.ErrDriverNotFound)`.

Репозиторий рассчитан на последовательные вызовы из одного процесса.

## YAML-репозиторий

Тип `file.DriverRepYAML` находится в `repository/file/yaml.go`.
Все девять методов соответствуют таблице JSON-репозитория.

Для использования YAML в примере выше замените создание репозитория:

```go
repo, err := file.NewDriverRepYAML("data/drivers.yaml")
```

Остальные вызовы остаются такими же. Нумерация страниц, сортировка по фамилии,
назначение ID и явное сохранение через `WriteAll` работают одинаково.

Чтение и запись используют [go-yaml v3](https://pkg.go.dev/go.yaml.in/yaml/v3).
Зависимость закреплена в `go.mod` и `go.sum`.

Файл содержит один YAML-документ со списком водителей. Пример находится в
`data/drivers.yaml`. Номер телефона записывается строкой в кавычках, чтобы
сохранить знак `+`.

```yaml
- id: 1
  last_name: Иванов
  first_name: Иван
  middle_name: Иванович
  phone_number: "+79161234567"
  experience: 10
```

Пустой список записывается как `[]`. Пустой файл, `null`, неизвестные поля,
повторяющиеся ключи, несколько документов, неверные значения и повторяющиеся
ID возвращают ошибку. Неудачная загрузка сохраняет прежний список и счётчик ID.

Оба репозитория встраивают `driverFileRepository` из `repository/file/base.go`.
В нём находятся файловые операции, работа со списком и общая структура
`driverRecord`. Каждый формат задаёт свои функции кодирования и разбора.
Создание моделей и проверка полей выполняются кодом первой лабораторной.

Диаграмма классов ЛР2: [lab02_class_diagram.md](lab02_class_diagram.md).

## Репозиторий PostgreSQL

`repository/postgres/repository.go` содержит `DriverRepDB` с шестью операциями
из пункта 4. Вызовы используют параметризованные SQL-запросы через `database/sql`.
Страницы нумеруются с 1 и идут в порядке ID. ID новых записей назначает PostgreSQL.

Из папки `lab01` примените миграцию к выбранной базе:

```sh
psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f migrations/001_create_drivers.sql
```

Получите общий объект через `database.GetDatabase(ctx, dsn)` и передайте его
в `postgres.NewDriverRepDB(db)`. Репозиторий делегирует выполнение запросов
объекту `Database`, который управляет пулом `*sql.DB`.

Первый успешный вызов `GetDatabase` закрепляет строку подключения за процессом.
Повторный вызов с той же строкой возвращает тот же объект. После ошибки
подключения можно повторить попытку. Вызывайте `db.Close()` при завершении
программы. После закрытия общий пул больше недоступен в этом процессе.

Интеграционный тест использует временную таблицу и не меняет постоянные
таблицы базы:

```sh
TEST_DATABASE_URL="$DATABASE_URL" go test -v ./repository/postgres
```

## Общий интерфейс репозиториев

`repository.DriverRepository` объединяет шесть операций, доступных у JSON,
YAML и PostgreSQL. Файловые операции `ReadAll`, `WriteAll` и
`SortByLastName` доступны у файловых типов.

Файловые репозитории можно передавать в функции, принимающие
`repository.DriverRepository`, напрямую. Для PostgreSQL оберните
`DriverRepDB` в адаптер:

```go
var repo repository.DriverRepository = postgres.NewDriverDBAdapter(ctx, postgres.NewDriverRepDB(db))
count, err := repo.GetCount()
```

Контекст передаётся в каждый SQL-запрос. Создавайте адаптер на время жизни
этого контекста. У всех реализаций `GetCount()` возвращает `(int, error)`.
Для файлов ошибка подсчёта всегда равна `nil`.

## Фильтрация и сортировка PostgreSQL

`NewDriverDBDecorator` оборачивает адаптер и возвращает реализацию
`repository.DriverRepository`. Например, следующая настройка выбирает
водителей с фамилией Петров и стажем от пяти лет, начиная с самых опытных.

```go
repo, err := postgres.NewDriverDBDecorator(
    postgres.NewDriverDBAdapter(ctx, postgres.NewDriverRepDB(db)),
    postgres.DriverFilter{LastName: "Петров", MinExperience: 5},
    postgres.DriverSort{Field: postgres.SortByExperience, Descending: true},
)
if err != nil {
    return err
}

count, err := repo.GetCount()
if err != nil {
    return err
}
page, err := repo.GetKNShortList(1, 20)
```

Поле `LastName` задаёт точное совпадение. Пустая фамилия и нулевой
минимальный стаж оставляют все записи в выборке. Можно сортировать по ID,
фамилии или стажу в обоих направлениях. При равенстве фамилий или стажа
записи идут по ID по возрастанию, чтобы порядок страниц был стабильным.

`GetCount` учитывает фильтр для всего списка. `GetKNShortList` применяет
фильтр и сортировку до выбора страницы. Остальные операции делегируются
адаптеру.

## Фильтрация и сортировка JSON и YAML

`file.NewDriverFileDecorator` принимает JSON- или YAML-репозиторий и
настройки из пакета `repository`. Например, чтобы на странице сначала шли
водители с наибольшим стажем:

```go
view, err := file.NewDriverFileDecorator(
    jsonRepo,
    repository.DriverFilter{LastName: "Петров", MinExperience: 5},
    repository.DriverSort{Field: repository.SortByExperience, Descending: true},
)
if err != nil {
    return err
}
count, err := view.GetCount()
if err != nil {
    return err
}
page, err := view.GetKNShortList(1, 20)
```

Декоратор применяет настройки к текущему списку в памяти. `GetCount`
считает записи по тому же фильтру. Сортировка выполняется на отдельной
выборке и сохраняет порядок исходного списка и файла. Чтение, запись
и операции по ID выполняет обёрнутый репозиторий.
