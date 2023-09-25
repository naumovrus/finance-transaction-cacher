migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up


run_docker:
	docker run --name=finance-api-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres

run_redis:
	docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 -d --rm redis/redis-stack:latest 

run:
	go run cmd/main.go
