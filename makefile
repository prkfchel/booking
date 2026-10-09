include .env
export
export PROJECT_ROOT=${shell pwd}


ps:
	@docker compose ps

env-up:
	docker compose up -d booking-postgres

env-down:
	docker compose down booking-postgres

env-cleanup:
	read -p "Очистить все volume файлы окружения? опасность утери данных [y, N]: " ans; \
    if [ "$$ans" = "y" ]; then \
            docker compose down booking-postgres && \
            sudo rm -rf "${PROJECT_ROOT}/out/pgdata"&& \
            echo "Файлы окружения очищены"; \
    else \
            echo "очистка окружения отменена"; \
    fi
	
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "отсутствует параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi

	@docker compose run --rm booking-migrate-postgres \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"
	@sudo chown -R $$USER:$$USER "${PROJECT_ROOT}/migrations"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "отсутствует параметр action"; \
		exit 1; \
	fi

	@docker compose run --rm booking-migrate-postgres \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@booking-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"


migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

env-port-forwarder:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

app-run:
	@export LOGGER_FOLDER="${PROJECT_ROOT}/out/logs" && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run "${PROJECT_ROOT}/cmd/booking/main.go"



logs-cleanup:
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf "${PROJECT_ROOT}/out/logs" && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi


app-deploy:
	@docker compose up -d --build booking


app-undeploy:
	@docker compose down booking


