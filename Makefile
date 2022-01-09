tests:
	go test ./...

tests-functional:
	go test

tests-functional-update-ct:
	go test --update

arch-next:
	@echo "-version:"
	go run main.go version
	@echo "-status:"
	go run main.go check

arch-prev:
	@echo "-version:"
	docker run --rm fe3dback/go-arch-lint:latest-stable-release version
	@echo "-status:"
	docker run --rm \
		-v ${PWD}:/app \
		fe3dback/go-arch-lint:latest-stable-release check --project-path /app

release-dry:
	@echo "check config.."
	goreleaser check
	@echo "build dry release.."
	goreleaser --snapshot --skip-publish --rm-dist

generate-std:
	docker build -t go-std-extractor \
		./build/stdextractor
	docker run --rm \
		-v ${PWD}/internal/generated/stdgo:/dest/stdgo \
		go-std-extractor -output "/dest/stdgo/packages.go" -goVersion "$$(go version)"
