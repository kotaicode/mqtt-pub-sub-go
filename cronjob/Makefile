REGISTRY_URL=xxxxxxxxxxxxxxxxxxxxxxx/workshop
PROJECT_NAMESPACE=workshop
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE_TAG=${REGISTRY_URL}/publisher-cronjob:${COMMIT_SHA}

export GO111MODULE=on

all: clean build run

run:
	BROKER_URL=tcp://localhost:1883 \
			   AUTH=workshop-auth \
			   PUBLISHER_ID=pub-cronjob \
			   PUBLISH_TOPIC=mqtt-k8s-workshop \
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
	kubectl patch cronjob pub-cronjob --patch '{"spec": {"template": {"spec": {"containers": [{"name": "pub-cronjob","image": "${IMAGE_TAG}"}]}}}}'
