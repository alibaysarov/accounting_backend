
copy_env:
	sed 's/=.*$/=/' .env > .env.example


# Создать миграцию автоматически (после изменения модели)
# alembic revision --autogenerate -m "add users table"

# Применить все миграции
# alembic upgrade head