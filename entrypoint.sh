#!/bin/sh
set -euo pipefail

# 1. Читаем секрет из файла, который Docker автоматически создает
# Убедитесь, что имя файла соответствует имени секрета в docker-compose.yml
# (у вас это 'db_password')
DB_PASS="$(tr -d '\r\n' < /run/secrets/db_password)"

# 2. Собираем строку подключения с использованием переменных окружения,
# которые уже определены в .env, и прочитанного пароля.
# $DB_USER, $DB_NAME, $DB_PORT должны быть определены в .env
DB_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Running database migrations..."
echo DB_URL

# 3. Выполняем команду 'migrate' с собранной строкой подключения
exec /migrate -path /migrations -database "$DB_URL" up