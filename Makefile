tests:
	go test ./...

tests-functional:
	go test

tests-functional-update-ct:
	go test --update

arch:
	go run main.go check --project-path ${PWD}

build-dev:
	docker build -t fe3dback/go-arch-lint:dev .