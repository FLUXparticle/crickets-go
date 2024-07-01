# Stage 1: Build Angular Frontend
FROM node:22-alpine3.19 AS build-frontend
WORKDIR /dist
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.22-alpine3.19 AS build-backend
WORKDIR /dist
COPY backend/go.* ./
RUN go mod download
COPY backend/ .
RUN go build -o /dist/backend crickets.go

# Stage 3: Build the final container
FROM alpine:3.19
WORKDIR /dist
COPY --from=build-frontend /dist/app /dist/app
COPY --from=build-backend /dist/backend /dist/backend
COPY /static /dist/static

EXPOSE 8080
CMD ["./backend"]
