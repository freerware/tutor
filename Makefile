all: start

start: bins
	@docker compose --file ./docker/docker-compose.yaml up

start-detached: bins
	@docker compose --file ./docker/docker-compose.yaml up -d

bins:
	@GO111MODULE=on go build .

restart:
	@docker compose --file ./docker/docker-compose.yaml restart

down:
	@docker compose --file ./docker/docker-compose.yaml down 

clean: down 
	@docker image rm docker_tutor
	@GO111MODULE=on go clean -x

logs:
	@docker compose --file ./docker/docker-compose.yaml logs

local: export SERVER_HOST=0.0.0.0
local: export SERVER_PORT=8000
local: export DB_HOST=0.0.0.0
local: export DB_PORT=3306
local: export DB_USER=web_app
local: export DB_PASSWORD=web_app_password
local: export DB_PARSE_TIME=true
local: export DB_CHARSET=utf8mb4
local: export DB_NAME=tutor
local: export REPORTING_HOST=127.0.0.1
local: export REPORTING_PORT=8125
local: export REPORTING_MAX_FLUSH_INTERVAL=150
local: export REPORTING_MAX_FLUSH_BYTES=512

local: bins
	@docker compose --file ./docker/docker-compose.debug.yaml up -d
	./tutor
	# nvim main.go

debug: bins
	@docker compose --file ./docker/docker-compose.debug.yaml up -d

debug-db: debug
	@docker exec -it docker_tutor-db_1 mysql -u web_app -p
