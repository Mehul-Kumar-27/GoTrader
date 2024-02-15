SCRAPER_FILE=scraperApp
STORE_FILE=storeApp
BROKER_FILE=serverApp

TAG := $(shell git rev-parse --short HEAD)

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


swarm:
	@echo "Deploying the containers"
	docker stack deploy -c swarm.yaml gotrader


build_scraper_image:
	@echo "Building scraper image"
	docker build -t mehulkumar27/gotrader_scraper:${TAG} -f ./scraper/cmd/scraper.dockerfile ./scraper/cmd

build_server_image:
	@echo "Building server image"
	docker build -t mehulkumar27/gotrader_server:${TAG} -f ./server/cmd/server.dockerfile ./server/cmd

build_store_image:
	@echo "Building store image"
	docker build -t mehulkumar27/gotrader_store:${TAG} -f ./store/cmd/store.dockerfile ./store/cmd

image: build_scraper_image build_server_image build_store_image
	@echo "Pushing the images"
	docker push mehulkumar27/gotrader_scraper:${TAG}
	docker push mehulkumar27/gotrader_server:${TAG}
	docker push mehulkumar27/gotrader_store:${TAG}