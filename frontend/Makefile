REGISTRY_URL=xxxxxxxxxxxxxxxxxxxxxxx/workshop
PROJECT_NAMESPACE=workshop
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE_TAG=${REGISTRY_URL}/frontend:${COMMIT_SHA}

run-dev:
	docker-compose -f docker-compose-dev.yml up --build

run-prd:
	docker-compose -f docker-compose-prd.yml up --build

build-push:
	docker build --build-arg REACT_APP_BACKEND_URL=https://backend.workshop.shipit.solutions -t ${IMAGE_TAG} .
	docker push ${IMAGE_TAG}

patch-deployment:
	kubectl patch deployment frontend --patch '{"spec": {"template": {"spec": {"containers": [{"name": "frontend","image": "${IMAGE_TAG}"}]}}}}'
