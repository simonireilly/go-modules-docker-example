FOLDER=$(notdir $(shell pwd))
DOCKER_CONTAINER_NAME="${FOLDER}_go-docker"
DOCKER_CONTAINER_VOLUME_NAME="${FOLDER}_go-docker-volume"
DOCKER_CONTAINER_OPTIMIZED_NAME="${FOLDER}_go-docker-optimized"

#
# Server
#

build:
	docker build -t go-docker -f Dockerfile .

run:
	docker run -d -p 8080:8080 --name ${DOCKER_CONTAINER_NAME} go-docker

stop:
	docker container stop ${DOCKER_CONTAINER_NAME}

#
#	Server and log volume
#

build-volume:
	docker build -t go-docker-volume -f Dockerfile.volume .

run-volume:
	docker run -d -p 8080:8080 -v ${shell pwd}/logs/go-docker:/app/logs --name ${DOCKER_CONTAINER_VOLUME_NAME} go-docker-volume

stop-volume:
	docker container stop ${DOCKER_CONTAINER_VOLUME_NAME}

#
# Optomised Server
#

build-optimized:
	docker build -t go-docker-optimized -f Dockerfile.multistage .

run-optimized:
	docker run -d -p 8080:8080 --name ${DOCKER_CONTAINER_OPTIMIZED_NAME} go-docker-optimized

stop-optimized:
	docker container stop ${DOCKER_CONTAINER_OPTIMIZED_NAME}
