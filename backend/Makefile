REGISTRY_URL=xxxxxxxxxxxxxxxxxxxxxxx/workshop
PROJECT_NAMESPACE=workshop
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE_TAG=${REGISTRY_URL}/backend:${COMMIT_SHA}

export GO111MODULE=on

all: clean build run

build: *.go
	go mod download
	go mod tidy
	go mod vendor
	go build -o backend

run:
	BROKER_URL=tcp://localhost:1883 \
			   AUTH=workshop-auth \
			   PORT=8888 \
			   SUBSCRIBE_TOPIC=mqtt-k8s-workshop \
			   ./backend

clean:
	rm -rf ./backend

build-push:
	docker build -t ${IMAGE_TAG} .
	docker push ${IMAGE_TAG}

patch-deployment:
	kubectl patch deployment backend --patch '{"spec": {"template": {"spec": {"containers": [{"name": "backend","image": "${IMAGE_TAG}"}]}}}}'
