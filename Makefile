include .env
export

dev:
	go run cmd/main.go
copy_env:
	sed 's/=.*$/=/' .env > .env.example

migration:
	goose -dir migrations create $(name) sql

migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down

pre-commit:
	pre-commit run --all-files --hook-stage pre-commit

pre-push:
	pre-commit run --all-files --hook-stage pre-push

show-coverage:
	go test ./internal/service/... ./internal/repository/... -coverprofile=coverage.out
	go tool cover -func=coverage.out | grep -v "100.0%"
# Создать миграцию автоматически (после изменения модели)
# alembic revision --autogenerate -m "add users table"

# Применить все миграции
# alembic upgrade head