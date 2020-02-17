REGISTRY_URL=xxxxxxxxxxxxxxxxxxxxxxx/workshop
PROJECT_NAMESPACE=workshop
COMMIT_SHA=$(shell git rev-parse HEAD)
IMAGE_TAG=${REGISTRY_URL}/broker:${COMMIT_SHA}

build-push:
	cd ./hmq && go mod download && go mod vendor
	docker build -t ${IMAGE_TAG} ./hmq
	docker push ${IMAGE_TAG}

deploy:
	cd ops/k8s/manifests
	kustomize edit set image ${IMAGE_TAG}
	kustomize build | kubectl -n ${PROJECT_NAMESPACE} apply -f -

patch-deployment:
	kubectl patch deployment broker --patch '{"spec": {"template": {"spec": {"containers": [{"name": "broker","image": "${IMAGE_TAG}"}]}}}}'
