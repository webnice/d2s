# d2s

[![PkgGoDev](https://pkg.go.dev/badge/github.com/webnice/d2s/v2.svg)](https://pkg.go.dev/github.com/webnice/d2s/v2)
[![Coverage Status](https://coveralls.io/repos/github/webnice/d2s/badge.svg?branch=v2)](https://coveralls.io/github/webnice/d2s?branch=v2)

Генератор Golang структур по структуре таблиц базы данных

## Описание

`d2s` — это генератор Golang структур по структуре таблиц базы данных.
Утилита `db2struct` подключается к MySQL или MariaDB, загружает информацию о
таблице из `INFORMATION_SCHEMA` и создаёт Golang файл со структурой, тегами
`db` (для `database/sql` и sqlx) и тегами `gorm` (по документации GORM v2).

Поля, допускающие значение NULL, генерируются с типами-обёртками библиотеки
[`github.com/webnice/nix/v2`](https://github.com/webnice/nix), реализующими
интерфейсы `sql.Scanner` и `driver.Valuer`.

### Особенности

- Поддержка всех типов данных MySQL 8.0 и MariaDB: целочисленные, числа с
  плавающей точкой, битовые, дата и время, строковые, бинарные, `json` и
  пространственные типы;
- Нативные типы MariaDB: `uuid` → `nix.UUID`, `inet4`/`inet6` → `nix.IP`,
  `vector` → `[]byte`;
- Точное сопоставление типов: `tinyint(1)` → `bool`, `tinyint(N>1)` → `int8`,
  `smallint` → `int16`, `int` → `int32`, `bigint` → `int64`,
  `float` → `float32`, `double` → `float64`, `decimal(p,0)` → `nix.BigInt`,
  `time` → `time.Duration`, `set` → `[]string`, `json` → `json.RawMessage`;
- Nullable-поля через типы-обёртки `nix/v2`: значение не теряется, NULL
  читается и пишется корректно;
- Переопределение типов полей двумя способами: аргументом командной строки и
  маркером `{тип}` в комментарии колонки базы данных;
- Автоматические импорты сгенерированного файла: `time`, `encoding/json`,
  `github.com/webnice/nix/v2`;
- Интеграция с `go generate`.

## Требования

- Go 1.26 или новее;
- MySQL 5.7 или новее, MariaDB 10.5 или новее (используется `INFORMATION_SCHEMA`);
- Для сгенерированного кода — `github.com/webnice/nix/v2` v2.1.0 или новее.

## Установка

Утилита `db2struct`:

```bash
go install github.com/webnice/d2s/v2/db2struct@latest
```

Подключение как библиотеки:

```bash
go get -u github.com/webnice/d2s/v2
```

#### Зависимости

| Зависимость                             | Назначение                  |
|-----------------------------------------|-----------------------------|
| `github.com/go-sql-driver/mysql`        | драйвер MySQL               |
| `github.com/webnice/kit/v4`             | логирование                 |
| `github.com/alecthomas/kingpin/v2`      | разбор командной строки     |
| `github.com/webnice/nix/v2`             | типы nullable-полей (для сгенерированного кода) |

## Использование

Команда генерации добавляется в Golang файл как `go:generate`-директива:

```go
//go:generate db2struct create "" "message" "types" "Message" "message_model.go"
```

Аргументы команды `create`:

```
db2struct create [database] [table] [package] [struct] [file] [types]
```

| Аргумент   | Описание                                                        |
|------------|-----------------------------------------------------------------|
| `database` | Имя базы данных. Если пусто — определяется из DSN или соединения |
| `table`    | Имя таблицы                                                     |
| `package`  | Имя пакета сгенерированного файла                               |
| `struct`   | Имя создаваемой структуры                                       |
| `file`     | Имя создаваемого файла                                          |
| `types`    | Переопределение типов полей: `"поле1:тип;поле2:тип"` (необязательно) |

Глобальные флаги:

| Флаг           | Переменная окружения | Описание                                  |
|----------------|----------------------|-------------------------------------------|
| `-d`, `--debug`| `DB2STRUCT_DEBUG`    | Включение режима отладки                  |
| `-b`, `--drv`  | `DB2STRUCT_DRV`      | Драйвер базы данных (по умолчанию `mysql`)|
| `-u`, `--dsn`  | `DB2STRUCT_DSN`      | Имя источника данных (DSN)                |

Пример DSN для соединения по TCP:

```
user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true
```

## Сопоставление типов данных

Полная таблица сопоставления приведена в документации пакета (`go doc` или
[doc.go](doc.go)). Краткая сводка (NOT NULL / NULL):

| Тип MySQL                    | NOT NULL        | NULL              |
|------------------------------|-----------------|-------------------|
| `tinyint(1)`, `boolean`      | `bool`          | `nix.Bool`        |
| `tinyint(N>1)`               | `int8`          | `nix.Int8`        |
| `smallint`                   | `int16`         | `nix.Int16`       |
| `mediumint`, `int`           | `int32`         | `nix.Int32`       |
| `bigint`                     | `int64`         | `nix.Int64`       |
| Целочисленные `unsigned`     | `uint16/32/64`  | `nix.Uint16/32/64`|
| `float`                      | `float32`       | `nix.Float32`     |
| `double`                     | `float64`       | `nix.Float64`     |
| `decimal(p,0)`               | `nix.BigInt`    | `nix.BigInt`      |
| `decimal(p,s>0)`             | `float64`       | `nix.Float64`     |
| `bit(1)`                     | `bool`          | `nix.Bool`        |
| `bit(N>1)`                   | `[]byte`        | `nix.Bytes`       |
| `date`, `datetime`, `timestamp` | `time.Time`  | `nix.Time`        |
| `time`                       | `time.Duration` | `nix.Duration`    |
| `year`                       | `int16`         | `nix.Int16`       |
| `char`, `varchar`, `text`    | `string`        | `nix.String`      |
| `enum`                       | `string`        | `nix.String`      |
| `set`                        | `[]string`      | `nix.StringSlice` |
| `binary`, `blob`             | `[]byte`        | `nix.Bytes`       |
| `json`                       | `json.RawMessage` | `nix.RawMessage`|
| Пространственные типы        | `[]byte`        | `nix.Bytes`       |
| `uuid` (MariaDB 10.7+)       | `nix.UUID`      | `nix.UUID`        |
| `inet4`, `inet6` (MariaDB 10.5+) | `nix.IP`    | `nix.IP`          |
| `vector` (MariaDB 11.7+)     | `[]byte`        | `nix.Bytes`       |

Особенность MariaDB: `JSON` — псевдоним `LONGTEXT` с проверкой `json_valid()`,
поэтому `INFORMATION_SCHEMA` сообщает `DATA_TYPE longtext` и колонка
сопоставляется с `string`/`nix.String`. Для работы с JSON примените
переопределение типа `{nix.RawMessage}` в комментарии колонки или аргументом
командной строки.

## Переопределение типов полей

### Через командную строку

Шестой аргумент команды `create` — список `"поле1:тип;поле2:тип"`, где поле —
имя колонки базы данных (без учёта регистра), тип — стандартный тип Golang
или тип `nix/v2`:

```go
//go:generate db2struct create "" "message" "types" "Message" "message_model.go" "id:uint64;payload:nix.RawMessage;amount:nix.BigInt;created_at:time.Time"
```

Неизвестное имя типа — ошибка генерации со списком допустимых типов. Поле,
отсутствующее в схеме таблицы, игнорируется.

### Через комментарий в базе данных

Если в комментарии колонки (`COLUMN_COMMENT`) в любом месте встречается группа
в фигурных скобках, совпадающая с допустимым типом данных, — это
переопределение типа поля:

```sql
ALTER TABLE message MODIFY COLUMN ttl TIME NOT NULL COMMENT 'Время жизни {nix.Duration}';
ALTER TABLE message MODIFY COLUMN amount DECIMAL(20,0) NOT NULL COMMENT 'Сумма {nix.BigInt} в копейках';
```

При нескольких группах `{...}` берётся первая совпавшая с типом данных.
Маркер `{...}` удаляется из комментария в сгенерированном коде.

### Приоритет переопределений

1. Командная строка — наивысший приоритет.
2. Комментарий к полю базы данных.
3. Базовый маппинг — сопоставление по умолчанию.

### Nullable-правила

- Указан базовый тип (`bool`, `string`, `int64`, `time.Time`, ...) и колонка
  nullable — тип автоматически заменяется на nix-обёртку (`bool` → `nix.Bool`,
  `time.Time` → `nix.Time` и т.д.);
- Указан тип nix явно (`nix.BigInt`, `nix.RawMessage`, ...) — используется
  как есть, для nullable и not null колонок.

### Допустимые типы

Стандартные типы Go: `bool`, `string`, `int`, `int8`–`int64`, `uint`,
`uint8`–`uint64`, `uintptr`, `float32`, `float64`, `complex64`, `complex128`,
`[]byte`, `[]string`, `time.Time`, `time.Duration`, `json.RawMessage`.

Типы `github.com/webnice/nix/v2`: `nix.Bool`, `nix.String`, `nix.Bytes`,
`nix.Time`, `nix.Duration`, `nix.BigInt`, `nix.RawMessage`, `nix.StringSlice`,
`nix.UUID`, `nix.IP`, все числовые типы `nix.*` и дженерик `nix.Nullable[T]`.

## Пример сгенерированной структуры

Таблица `message` (MySQL):

```sql
CREATE TABLE message (
  id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  payload BLOB DEFAULT NULL,
  amount DECIMAL(20,0) NOT NULL COMMENT 'Сумма {nix.BigInt} в копейках',
  ttl TIME NOT NULL COMMENT 'Время жизни {nix.Duration}',
  tags SET('a','b') DEFAULT NULL,
  created_at DATETIME NOT NULL,
  attr JSON DEFAULT NULL,
  ratio FLOAT NOT NULL,
  cnt TINYINT(4) NOT NULL,
  PRIMARY KEY (id)
);
```

Команда генерации:

```go
//go:generate db2struct create "" "message" "types" "Message" "message_model.go" "payload:nix.UUID"
```

Результат (фрагмент):

```go
type Message struct {
	ID        uint64          `db:"id"         gorm:"column:id;primaryKey;autoIncrement;not null;precision:20;type:BIGINT(20) UNSIGNED"` // 
	Payload   nix.UUID        `db:"payload"    gorm:"column:payload;default:null;type:BLOB"` // 
	Amount    nix.BigInt      `db:"amount"     gorm:"column:amount;not null;type:DECIMAL(20,0)"` // Сумма в копейках
	TTL       nix.Duration    `db:"ttl"        gorm:"column:ttl;not null;type:TIME"` // Время жизни
	Tags      nix.StringSlice `db:"tags"       gorm:"column:tags;default:null;type:SET('a','b')"` // 
	CreatedAt time.Time       `db:"created_at" gorm:"column:created_at;not null;type:DATETIME"` // 
	Attr      nix.RawMessage  `db:"attr"       gorm:"column:attr;default:null;type:JSON"` // 
	Ratio     float32         `db:"ratio"      gorm:"column:ratio;not null;precision:12;type:FLOAT"` // 
	Cnt       int8            `db:"cnt"        gorm:"column:cnt;not null;type:TINYINT(4)"` // 
}
```

## Тестирование

```bash
# Стандартный прогон.
go test ./...

# Прогон через Makefile.
make test
```

## Лицензия

© WEB NICE, распространяется под [лицензией MIT](LICENSE).
