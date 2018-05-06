MUTABLE_VERSION ?= canary
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || echo "git-$$(git rev-parse --short HEAD)")

IMAGE_PREFIX ?= teamhephy
IMAGE := ${DEV_REGISTRY}${IMAGE_PREFIX}/${SHORT_NAME}:${VERSION}
MUTABLE_IMAGE := ${DEV_REGISTRY}${IMAGE_PREFIX}/${SHORT_NAME}:${MUTABLE_VERSION}

info:
	@echo "Build tag:       ${VERSION}"
	@echo "Registry:        ${DEV_REGISTRY}"
	@echo "Immutable tag:   ${IMAGE}"
	@echo "Mutable tag:     ${MUTABLE_IMAGE}"

.PHONY: docker-push
docker-push: docker-immutable-push docker-mutable-push

.PHONY: docker-immutable-push
docker-immutable-push:
	docker push ${IMAGE}

.PHONY: docker-mutable-push
docker-mutable-push:
	docker push ${MUTABLE_IMAGE}
