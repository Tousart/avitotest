### Стек

- Golang
- PostgreSQL
- Docker


### Запуск

Для запуска можно использовать команду:

```
docker-compose up
```

Также можно вызвать команду:

```
make run
```

Команда выше собирает проект (`docker compose build`) и запускает (`docker compose up`).


### База данных

<img width="1098" height="990" alt="Снимок экрана 2025-11-16 214705" src="https://github.com/user-attachments/assets/71767eb0-ca84-4b89-a2b3-cb1600e8a5c5" />


Дополнительная таблица `pr_reviewers` нужна для того, чтобы отслеживать, кто ревьюит конкретный pull request.

Ограничений типы данных в таблице не было, поэтому все string как VARCHAR(64).

Нагрузки по условию небольшие, поэтому индексов по id (у primary keys) достаточно.

Миграции запускаются отдельным контейнером (написано в docker-compose.yaml).

### Дополнительно

Реализовал метод `/users/getActivity` (без параметров), который возвращает:
- id пользователя
- имя пользователя
- общее количество ПР, которые он ревьюит
- количество ПР (от общего), которые MERGED
- количество ПР (от общего), которые OPEN
Данные отсортированы по убыванию количества ПР.

Также добавил коды ошибок на `Internal Error` и `Bad Request`.

### Примеры запросов

Создание команды:

```
curl -X POST http://localhost:8080/team/add -d '{"team_name": "nambavan", "members": [{"user_id": "u1", "username": "Alice", "is_active": true}, {"user_id": "u2", "username": "Bob", "is_active": true}, {"user_id": "u3", "username": "Victor", "is_active": true}, {"user_id": "u4", "username": "Maria", "is_active": true}]}'
```

Ответ:

```
{"team_name":"nambavan","members":[{"user_id":"u1","username":"Alice","is_active":true},{"user_id":"u2","username":"Bob","is_active":true},{"user_id":"u3","username":"Victor","is_active":true},{"user_id":"u4","username":"Maria","is_active":true}]}
```

Изменение активности пользователя:

```
curl -X POST http://localhost:8080/users/setIsActive -d '{"user_id": "u2", "is_active": false}'
```

Ответ:

```
{"user_id":"u2","username":"Bob","team_name":"nambavan","is_active":false}
```

Создание ПР:

```
curl -X POST http://localhost:8080/pullRequest/create -d '{"pull_request_id": "pr-1228", "pull_request_name": "Bobs PR", "author_id": "u2"}'
```

Ответ:

```
{"pull_request_id":"pr-1228","pull_request_name":"Bobs PR","author_id":"u2","status":"OPEN","assigned_reviewers":["u3","u1"],"created_at":"2025-11-16T19:05:40Z","merged_at":""}
```

Переназначение (результат может быть разным в зависимости от имеющихся на ПР ревьюерах):

```
curl -X POST http://localhost:8080/pullRequest/reassign -d '{"pull_request_id":"pr-1228","old_user_id":"u3"}'
```

Ответ:

```
{"pull_request_id":"pr-1228","pull_request_name":"Bobs PR","author_id":"u2","status":"OPEN","assigned_reviewers":["u1","u4"],"created_at":"","merged_at":""}
```

Merge ПР (идемпотентный):

```
curl -X POST http://localhost:8080/pullRequest/merge -d '{"pull_request_id":"pr-1228"}'
```

Ответ:

```
{"pull_request_id":"pr-1228","pull_request_name":"Bobs PR","author_id":"u2","status":"MERGED","assigned_reviewers":null,"created_at":"2025-11-16T19:05:40Z","merged_at":"2025-11-16T19:08:07Z"}
```

Получение ПР ревьюера:

```
curl -X GET http://localhost:8080/users/getReview?user_id=u1
```

Ответ:

```
[{"pull_request_id":"pr-1228","pull_request_name":"Bobs PR","author_id":"u2","status":"MERGED"}]
```

Вывод статистики ревьюеров (мой метод):

```
curl -X GET http://localhost:8080/users/getActivity
```

Ответ:

```
[{"user_id":"u4","username":"Maria","pull_requests":1,"merged_pr":1,"open_pr":0},{"user_id":"u1","username":"Alice","pull_requests":1,"merged_pr":1,"open_pr":0}]
```
