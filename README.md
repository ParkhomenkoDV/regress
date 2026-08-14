# Regress = инструмент сравнения JSON файлов

**Regress** — мощная утилита для сравнения JSON файлов, которая помогает находить различия между двумя версиями данных и экспортировать результаты в удобном формате Excel.

![](./assets/images/testing.png)

## Requirements

- Go 1.25 или выше

## Usage

1. Создать директорию `before` и положи туда JSON файлы
1. Создай директорию `after` и положи туда **одноименные** JSON файлы
1. Запусти скрипт командой

```bash
go run regress.go [-before=S] [-after=S] [-all] [-workers=N]
```

Флаг       | Описание                                      | По умолчанию
-----------|-----------------------------------------------|-------------
`-before`  | Директория с JSON-файлами **ДО** изменений    | `before`
`-after`   | Директория с JSON-файлами **ПОСЛЕ** изменений | `after`
`-all`     | Показывать все поля даже без изменений        | `false`
`-workers` | Количество параллельных воркеров обработки    | количетсво ядер CPU


## Build

```bash
go build -o regress regress.go
./regress
```

## Install

```bash
go install
```

## Schema
```
before/      after/
| 1.json -   | 1.json -
| 2.json  |  | 2.json  |
| 3.json  |  | 3.json  |
| ...     |  | ...     |
          |            |  
          v            v
         ----------------
        | i.json  i.json |
        |    |       |   |
        |    v       v   |
        |    UnMarshal   | gorutine
        |      |   |     |
        |      v   v     |
        |     Compare    |
         ----------------
                 |
                 v
         comparison.xlsx
```

