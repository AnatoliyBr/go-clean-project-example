## Хранилище Счетчиков
Мой первый проект с чистой архитектурой на Go.

## Сервер
### Функции

### Архитектура

## REST API
### Пример использования
```bash
$ curl "localhost:8080/set?name=a&val=19"
Counter 'a' with val '19' was set
$ curl "localhost:8080/get?name=a"
a: 7
$ curl "localhost:8080/inc?name=a"
ok, a: 8
$ curl "localhost:8080/get?name=a"
a: 8
```

## Тестирование параллельных подключений
Для имитации множества одновременных подключений используем **JMeter (аналог ApacheBench)** - программу для тестирования веб-сервера.

```bash
$ ab -n 20000 -c 200 "localhost:8080/inc?name=i"
```

* `-n` - общее число запросов
* `-c` - число запросов, отправляемых одновременно

## Запуск и отладка
Все команды, используемые в процессе разработки и тестирования, фиксировались в **Makefile**. 

## Полезные ссылки
* [On concurrency in Go HTTP servers](https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers/)
* [Разработка REST-серверов на Go. Часть 1: стандартная библиотека](https://habr.com/ru/companies/ruvds/articles/559816/)