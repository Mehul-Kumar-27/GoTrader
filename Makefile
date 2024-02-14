SCRAPER_FILE=scraperApp
STORE_FILE=storeApp
BROKER_FILE=serverApp

build_scraper:
	@echo "Building scraper"
	cd ./scraper/cmd/api && env GOOS=linux CGO_ENABLED=0 go build -a -o ../$(SCRAPER_FILE) .

build_store:
	@echo "Building store"
	cd ./store/cmd/api && env GOOS=linux CGO_ENABLED=0 go build -a -o ../${STORE_FILE} .

build_server:
	@echo "Building server"
	cd ./server/cmd/api && env GOOS=linux CGO_ENABLED=0 go build -a -o ../${BROKER_FILE} .
up :
	@echo "Starting up the containers"
	docker-compose up 
down:
	@echo "Stopping the containers"
	docker-compose down
remove:
	@echo "Removing the images"
	docker rmi -f gotrader_scraper

inside:
	@echo "Entering the container"
	docker exec -it scraper /bin/sh

build: build_scraper build_store build_server
	@echo "Building the docker compose"
	docker-compose build --no-cache
