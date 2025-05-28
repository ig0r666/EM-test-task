## Запуск проетка
Запустить проект:
```Makefile 
make up
или
docker compose up --build -d
```
Остановить проект:
```Makefile 
make down
или
docker compose down
```
Полная очистка и остановка:
```Makefile 
make clean
или
docker compose down -v
```
Сгенерировать Swagger документацию:
```Makefile 
make build-docs
```

После запуска проекта документация доступна по эндпоинту /swagger/doc.json
