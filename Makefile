up: down
	docker compose up --build -d

down:
	docker compose down

clean:
	docker compose down -v

build-docs:
	swag init -g testtask/main.go --output testtask/adapters/rest/docs --parseDependency --parseInternal