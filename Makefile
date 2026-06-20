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


# Создать миграцию автоматически (после изменения модели)
# alembic revision --autogenerate -m "add users table"

# Применить все миграции
# alembic upgrade head