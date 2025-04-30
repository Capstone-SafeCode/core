PROJECT_NAME = core-api

# ========== ALL ==========

build:
	docker-compose build

up:
	docker-compose up -d

stop:
	docker-compose stop

down:
	docker-compose down

re: build up

# ========== API ==========

build-api:
	docker-compose build api

up-api:
	docker-compose up -d api

stop-api:
	docker-compsoe stop api

down-api:
	docker-compose stop api

re-api: build-api up-api

# ========== LOGS / TERMINAL ==========

logs:
	docker-compose logs -f

bash-api:
	docker exec -it core-api /bin/sh
