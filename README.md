# Начало работы

1. Копируем репозиторий и переходим в него:
```sh
git clone https://github.com/Egor213/PullReqService.git
```
```sh
cd PullReqService
```

2. Создаем в нем файл `.env` и копируем содержимое файла `.env.example`. Затем выставляем необходимые на ваш вкус настройки.

3. Для запуска проекта можно воспользовать коммандой:  
+ C логами:
```sh
docker-compose up --build -d && docker-compose logs -f  
```
+ Без логов:
```sh
docker-compose up --build -d  
```

После этого сервис будет доступен на `http://localhost:8080`
> При этом `api` находится по адресу: `http://localhost:8080/api/v1`

4. Для завершения работы программы можно использовать команды:
+ Для удаления данных бд:
```sh
docker compose down -v
```
+ Без удаления данных бд:
```sh
docker compose down
```

5. Для запуска тестов можно воспользоваться командой:  
```go
go test -v ./...
```

6. Для запуска линтера можно воспользоваться командой:  
```go
golangci-lint run
```

7. Для запуска интеграционных тестов воспользуйтесь коммандой:
```sh
docker compose --profile tests up --build
```
