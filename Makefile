export ORGANIZATION=flowup
export PROJECT_NAME=owl
export ESCAPED_PROJECT_NAME=owl
export PROJECT_CMD_NAME=owl

export PROJECT_GOPATH=github.com/${ORGANIZATION}/${PROJECT_NAME}
export PROJECT_PATH=/go/src/${PROJECT_GOPATH}
export CMD_GOPATH=${PROJECT_GOPATH}/cmd/${PROJECT_CMD_NAME}

init:
	docker-compose up init

build:
	docker-compose build

dev:
	docker-compose up dev

lib:
	docker-compose up -d empty && docker exec ${ESCAPED_PROJECT_NAME}_empty_1 glide get --non-interactive $(LIB) && docker-compose stop empty

stop:
	docker-compose stop

clean:
	docker-compose down

attach:
	docker-compose up -d empty && docker exec -it ${ESCAPED_PROJECT_NAME}_empty_1 /bin/bash && docker-compose stop empty
