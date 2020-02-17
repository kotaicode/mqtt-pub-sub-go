REGISTRY_URL=xxxxxxxxxxxxxxxxxxxxxxx/workshop
PROJECT_NAMESPACE=workshop
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE_TAG=${REGISTRY_URL}/publisher:${COMMIT_SHA}

export GO111MODULE=on

all: clean build run

run: build
	BROKER_URL=tcp://localhost:1883 \
			   AUTH=workshop-auth \
			   PORT=8081 \
			   PUBLISHER_ID=pub-1 \
			   PUBLISH_TOPIC=mqtt-k8s-workshop \
			   PUBLISH_INTERVAL=3 \
			   ./publisher

build: *.go
	go mod download
	go mod tidy
	go mod vendor
	go build -o publisher

clean:
	rm -rf ./publisher

build-push:
	docker build -t ${IMAGE_TAG} .
	docker push ${IMAGE_TAG}

patch-deployment:
	kubectl patch deployment pub-0 --patch '{"spec": {"template": {"spec": {"containers": [{"name": "pub-0","image": "${IMAGE_TAG}"}]}}}}'
