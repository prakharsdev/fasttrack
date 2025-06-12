# ------------------------------------
# CONFIGURATION
# ------------------------------------
DOCKER_COMPOSE = docker-compose -f build/docker-compose.yml
APP_NAME = fasttrack-app
LOCAL_ENTRY = cmd/app/main.go

# ------------------------------------
# DEVELOPMENT COMMANDS
# ------------------------------------

# Build Go binary locally 
build:
	go build -o fasttrack $(LOCAL_ENTRY)

# Run app locally (without Docker)
run-local:
	go run $(LOCAL_ENTRY)

# Run unit tests
test:
	go test ./...

# ------------------------------------
# DOCKER COMMANDS
# ------------------------------------

# Build & start everything (fresh build)
up:
	$(DOCKER_COMPOSE) up --build

# Start existing containers (skip rebuild)
start:
	$(DOCKER_COMPOSE) up

# Stop containers
down:
	$(DOCKER_COMPOSE) down

# View container logs
logs:
	$(DOCKER_COMPOSE) logs -f

# Rebuild only the app container
rebuild-app:
	$(DOCKER_COMPOSE) build app

# Full cleanup (containers + volumes + networks)
clean:
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

# ------------------------------------
# HELP
# ------------------------------------

help:
	@echo ""
	@echo "FastTrack Makefile Commands"
	@echo ""
	@echo "  make up            => Build & start full stack"
	@echo "  make start         => Start existing containers"
	@echo "  make rebuild-app   => Rebuild only app container"
	@echo "  make down          => Stop containers"
	@echo "  make clean         => Full Docker cleanup"
	@echo "  make logs          => View live logs"
	@echo ""
	@echo "  make build         => Build Go binary"
	@echo "  make run-local     => Run app locally"
	@echo "  make test          => Run tests"
	@echo ""
