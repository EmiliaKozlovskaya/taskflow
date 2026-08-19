#Чтобы каждый раз не писать длинные команды для запуска docker-compose, можно создать Makefile, 
#который будет содержать все необходимые команды для сборки и запуска проекта. 
#Однако чтобы не хранить все те же переменные окружения в открытом доступе в Makefile, 
#можно создать отдельный файл .env, который будет содержать все необходимые переменные. (его добавим в .gitignore, чтобы не коммитить его в репозиторий)
include .env
export

export PROJECT_ROOT=$(shell pwd)
#Устанавливаем переменную окружения PROJECT_ROOT, 
#которая будет содержать путь к корневой директории проекта.

#Чтобы запустить проект, достаточно будет выполнить команду `make env-up` в терминале, и все необходимые контейнеры будут запущены.
env-up:
	@docker compose up -d todoapp-postgres 
env-down:
	@docker compose down todoapp-postgres
#Команда `make env-cleanup` будет удалять все контейнеры, тома и сети, связанные с проектом.
env-cleanup:
	@read -p "Are you sure you want to remove all containers, volumes, and networks? (y/N) " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "All containers, volumes, and networks have been removed."; \
	else \
		echo "Cleanup cancelled."; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

#Команда `make migrate-create` будет создавать новую миграцию для базы данных, в которых уже прописываться будут sql запросы на изменение схемы базы данных.
#Имя миграции будет передаваться через аргумент `seq`, который указывается при вызове команды. Например, `make migrate-create seq=add_users_table` создаст новую миграцию с именем "add_users_table".
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Error: Migration name is required. Usage: make migrate-create seq=<migration_name>"; \
		exit 1; \
	fi
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"
#Команда `make migrate-up` будет запускать миграции для базы данных, а команда `make migrate-down` будет откатывать миграции.
#Вместо дублирования команд для запуска миграций, мы можем создать общую команду `make migrate-action`, 
#которая будет принимать аргумент `action`, определяющий, какую операцию выполнять (up или down).
migrate-up:
	@make migrate-action action=up 
migrate-down:
	@make migrate-action action=down
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Error: Action is required. Usage: make migrate-action action=<up|down>"; \
		exit 1; \
	fi
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)
#например make migrate-action action="down 3" - откатить последние 3 миграции, чтобы не писать отдельные таргеты для каждой команды.
logs-cleanup:
	@read -p "Are you sure you want to remove all logs? (y/N) " answer; \
	if [ "$$answer" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "All logs have been removed."; \
	else \
		echo "Cleanup cancelled."; \
	fi


todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

#--build для пересборки при каждом запуске
todoapp-deploy:
	@docker compose up -d --build todoapp

#чтобы обращаться к docker-compose не напрямую а через makefile сразу с .env
ps:
	@docker compose ps
