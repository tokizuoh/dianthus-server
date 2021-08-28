dc-build:
	docker-compose up --build -d
run:
	docker-compose exec app go run main.go