TCLIENT:=test-gclient
TSERVER:=test-gserver
DOCKER_REPO:=ikarpukhin

.PHONY: all clean proto client server build docker docker-publish docker-tag

all: clean proto build docker docker-tag docker-publish

build: client server

proto:
	protoc --go_out=plugins=grpc:${GOPATH}/src grpc/*.proto

server:
ifndef TGT_LINUX
	cd server && go build -o ${TSERVER} *.go
else
	cd server && GOOS=linux go build -o ${TSERVER} *.go
endif

client:
ifndef TGT_LINUX
	cd client && go build -o ${TCLIENT} *.go
else
	cd client && GOOS=linux go build -o ${TCLIENT} *.go
endif

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