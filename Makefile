test:
	go test ./...

arch:
	go-arch-lint check --project-path ${PWD}

build-dev:
	docker build -t fe3dback/go-arch-lint:dev .