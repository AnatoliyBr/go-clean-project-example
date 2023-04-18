## Хранилище Счетчиков
### Мотивация
Это мой первый проект на Go, его смысл в том, чтобы потренироваться:
* Организовывать код проекта согласно [чистой архитектуре](https://github.com/evrone/go-clean-template)
* Использовать **каналы** для синхронизации доступа к общим данным
* Использовать **инъекцию зависимостей** для обеспечения низкой связанности слоев

### Описание
Проект представляет из себя REST API сервер, который реализует **хранилище счетчиков**. В терминологии Go, это map'а, то есть набор из ключей (названий счетчиков) и их текущих значений.

Пользователям доступны следующие функции:
1. Создание счетчика по запросу: `/set?name=<имя>&val=<число>`
2. Просмотр текущего состояния счетчика по запросу: `/get?name=<имя>`
3. Увеличение счетчика на 1 по запросу: `/inc?name=<имя>`
4. Уменьшение счетчика на 1 по запросу: `/dec?name=<имя>`

## Сервер
### Алгоритм работы
Для обеспечения безопасного одновременного доступа обработчиков к данным можно использовать каналы или мьютексы. В данном проекте использовались каналы.

Обработчики запросов формируют команды и отправляют их в канал `cmds`.

Горутина, работающая в фоновом режиме, считывает команды из канала `cmds` и вызывает соответствующий `UseCase`, который в свою очередь вызывает соответствующий метод хранилища счетчиков `CounterStore`.

Метод хранилища счетчиков `CounterStore` возвращает ответ, который проходит через `UseCase` и попадает в горутину, которая отправляет его по каналу `replyChan`, которым обладает каждая команда.

Обработчики считывают информацию из `replyChan` и формируют ответ на запрос пользователя.

Данный алгоритм взаимодействия слоев соответствует идеологии чистой архитектуры:

```
HTTP > usecase
        usecase > repository
        usecase < repository
HTTP < usecase
```

### Архитектура
![architecture diagram](/assets/architecture.png "architecture diagram")

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