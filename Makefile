tests:
	go test ./...

tests-functional:
	go test

tests-functional-update-ct:
	go test --update

arch:
	go-arch-lint check --project-path ${PWD}

build-dev:
	docker build -t fe3dback/go-arch-lint:dev .