FRONTEND_SRC=$(shell find frontend/src -type f)
STATIC_SRC=$(shell find static -type f)
BACKEND_SRC=$(shell find backend -type f)

FRONTEND_DEST=dist/app/index.html
STATIC_DEST=$(STATIC_SRC:%=dist/%)
BACKEND_DEST=dist/backend

.PHONY: all backend frontend docker-container
all: backend frontend static

.PHONY: clean
clean:
	rm -rfv dist

.PHONY: backend
backend: $(BACKEND_DEST)
$(BACKEND_DEST): $(BACKEND_SRC)
	cd backend && go build -o ../$(BACKEND_DEST) .

.PHONY: frontend
frontend: $(FRONTEND_DEST)
$(FRONTEND_DEST): $(FRONTEND_SRC)
	cd frontend && npm install && ng build --output-path ../dist/app

.PHONY: static
static: $(STATIC_DEST)
dist/static/%: static/% dist/static
	cp $< $@

dist/static:
	mkdir -p $@

.PHONY: docker-container
docker-container:
	docker buildx build -t crickets:latest .
