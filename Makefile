FRONTEND_SRC=$(shell find frontend/src -type f)
BACKEND_SRC=$(shell find backend -type f)

FRONTEND_DEST=dist/static/app/index.html
BACKEND_DEST=dist/backend

.PHONY: all backend frontend docker-container
all: backend frontend

.PHONY: clean
clean:
	rm -rfv dist/

backend: $(BACKEND_DEST)
$(BACKEND_DEST): $(BACKEND_SRC)
	cd backend && go build -o ../$(BACKEND_DEST) main.go

frontend: $(FRONTEND_DEST)
$(FRONTEND_DEST): $(FRONTEND_SRC)
	cd frontend && npm install && npm run build

docker-container:
	docker build -t crickets:latest .
