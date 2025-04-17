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

# ========== MONGO ==========

build-mongo:
	docker-compose build mongo

up-mongo:
	docker-compose up -d mongo

stop-mongo:
	docker-compsoe stop mongo

down-mongo:
	docker-compose stop mongo

re-mongo: build-mongo up-mongo

# ========== LOGS / TERMINAL ==========

logs:
	docker-compose logs -f

bash-api:
	docker exec -it core-api /bin/sh

mongo-shell:
	docker exec -it mongo mongosh
