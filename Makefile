tests:
	go test ./...

tests-functional:
	go test

tests-functional-update-ct:
	go test --update

arch:
	docker run --rm \
		-v ${PWD}:/app \
		fe3dback/go-arch-lint:latest-stable-release check --project-path /app

build-dev:
	docker build -t fe3dback/go-arch-lint:dev .

release-dry:
	@echo "check config.."
	goreleaser check
	@echo "build dry release.."
	goreleaser --snapshot --skip-publish --rm-dist
