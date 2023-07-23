export DOCKER_BUILDKIT = 1
compose = docker compose -f dockerfiles/docker-compose.yml

test_db:
	$(compose) up -d postgres --build
	$(compose) run --rm wait -c postgres:5432

test_db_migrate: test_db
	docker run --rm -v "$(PWD):/app" -w /app --network=host arigaio/atlas:0.12.0-alpine migrate apply --env test

integration_tests: test_db_migrate
	docker build --network=host -t team_app:integration_tests --target integration_tests -f dockerfiles/Dockerfile .

unit_tests:
	docker build -t team_app:unit_tests --target unit_tests -f dockerfiles/Dockerfile .

component_tests: test_db_migrate
	docker build --network=host -t team_app:component_tests --target component_tests -f dockerfiles/Dockerfile .

down:
	$(compose) down

up:
	$(compose) up -d --build
