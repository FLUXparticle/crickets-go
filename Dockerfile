# Stage 1: Build Angular Frontend
FROM node:22-alpine AS build-frontend
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.22-alpine AS build-backend
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend/ .
RUN go build -o /app/backend main.go

# Stage 3: Build the final container
FROM alpine:latest
WORKDIR /app
COPY --from=build-frontend /app/static /app/static
COPY --from=build-backend /app/backend /app/backend

EXPOSE 8080
CMD ["./backend"]
