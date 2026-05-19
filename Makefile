include .env

run:
	@go run cmd/app/main.go

migrate-up:
	@migrate \
		-path migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		up

migrate-down:
	@migrate \
		-path migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		down