TCLIENT:=test-gclient
TSERVER:=test-gserver
DOCKER_REPO:=ikarpukhin

.PHONY: all clean proto client server build docker docker-publish docker-tag

all: clean proto build docker docker-tag docker-publish

build: client server

proto:
	protoc --go_out=plugins=grpc:${GOPATH}/src grpc/*.proto

clean:
	rm -rf server/${TSERVER} client/${TCLIENT}

server:
ifndef TGT_LINUX
	cd server && CGO_ENABLED=0 go build -o ${TSERVER} *.go
else
	cd server && CGO_ENABLED=0 GOOS=linux go build -o ${TSERVER} *.go
endif

client:
ifndef TGT_LINUX
	cd client && CGO_ENABLED=0 go build -o ${TCLIENT} *.go
else
	cd client && CGO_ENABLED=0 GOOS=linux go build -o ${TCLIENT} *.go
endif

docker-client:
	cd client && docker build --tag=${TCLIENT} .

docker-server:
	cd server && docker build --tag=${TSERVER} .

docker: docker-client docker-server

docker-tag-client:
	docker tag ${TCLIENT} albertocsm/${TCLIENT}

docker-tag-server:
	docker tag ${TSERVER} albertocsm/${TSERVER}

docker-tag: docker-tag-client docker-tag-server

docker-publish-client:
	docker push albertocsm/${TCLIENT}

docker-publish-server:
	docker push albertocsm/${TSERVER}

publish-client: client docker-client docker-tag-client docker-publish-client
