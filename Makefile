.PHONY: all build-backend build-frontend docker-build
all: build-backend build-frontend docker-build

build-backend:
	cd backend && go build -o ../bin/backend main.go

build-frontend:
	cd frontend && npm install && npm run build

docker-build:
	docker build -t crickets:latest .
