# Getting started
- Выгрузите проект командой
```shell
    git clone 
```

- Запустите инфраструктуру для сервиса
```shell
    docker-compose up
```
- Примените миграции командой
```shell
    goose postgres "postgresql://postgres:postgres@127.0.0.1:5433/postgres?sslmode=disable" -dir db/migrations up 
```

- Для остановки инфрастуктуры выполните
```shell
    docker-compose down
```

- Сервис готов к запуску

## Примеры запросов

- Создание нового пользователя:
```shell
  curl --location 'http://localhost:8082/users' \
--header 'Content-Type: application/json' \
--data-raw '{
  "first_name": "Jhon",
  "last_name": "Snow",
  "login": "some_login",
  "email": "json@gmail.com",
  "password": "some_password"
}
```

- Логин под новым пользователем
```shell
  curl --location 'http://localhost:8082/login' \
--header 'Content-Type: text/plain' \
--data '{
    "login": "some_login",
    "password": "some_password"
}'

```

- Создание заметки
```shell
  curl --location 'http://localhost:8082/users/1/notes' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njc4Njg0MjksInN1YiI6IjUifQ.evPvK8bOjRDyuU4PoDr_-dObbFiSuvtfHboxwUlD9Hs' \
--data '{
    "title": "sometitle",
   "content" : "somecontent"
}'
```

- Получение списка заметок для пользователя (заменить id пользователя в path)
```shell 
  curl --location 'http://localhost:8082/users/1/notes' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njc4Njg1NDcsInN1YiI6IjUifQ._EdNI1XEfQURf8RoQQWH_NgOFUIEwWA1GwVGRMXLoBM' \
--data ''
```

- Попытка получения записок другого пользователя (заменить id на другого пользователя)
```shell 
  curl --location 'http://localhost:8082/users/2/notes' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Njc4Njg1NDcsInN1YiI6IjUifQ._EdNI1XEfQURf8RoQQWH_NgOFUIEwWA1GwVGRMXLoBM' \
--data ''
```