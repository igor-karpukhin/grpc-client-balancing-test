TCLIENT:=test-gclient
TSERVER:=test-gserver
DOCKER_REPO:=ikarpukhin

.PHONY: all clean proto client server build docker docker-publish docker-tag

all: clean proto build docker docker-tag docker-publish

build: client server

proto:
	protoc --go_out=plugins=grpc:${GOPATH}/src grpc/*.proto

server:
	cd server && go build -o ${TSERVER} *.go

client:
	cd client && go build -o ${TCLIENT} *.go

clean:
	rm -rf server/${TSERVER} client/${TCLIENT}

docker:
	cd client && docker build --tag=${TCLIENT} .
	cd server && docker build --tag=${TSERVER} .

docker-tag:
	docker tag ${TSERVER} ikarpukhin/${TSERVER}
	docker tag ${TCLIENT} ikarpukhin/${TCLIENT}

docker-publish:
	docker push ikarpukhin/${TSERVER}
	docker push ikarpukhin/${TCLIENT}