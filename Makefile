.PHONY: all build test clean install lint generate

all: test build

build:
	go build -o bin/prefix ./cmd/prefix

test:
	go test -v -race -coverprofile=coverage.out ./...

coverage: test
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/ coverage.out coverage.html

install:
	go install ./cmd/prefix

lint:
	golangci-lint run

docker:
	docker build -t prefix .

generate: install
	mkdir -p .tmp
	@command -v generate-fake-data >/dev/null 2>&1 || go install moul.io/generate-fake-data@latest

	echo 'foo@bar:~$$ prefix -h' > .tmp/usage.txt
	(prefix -h || true) 2>> .tmp/usage.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix' > .tmp/default.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr | prefix >> .tmp/default.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format="#{{.LineNumber3}} {{.ShortUptime}} {{.ShortDuration}} | "' > .tmp/example-1.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr | prefix --format="#{{.LineNumber3}} {{.ShortUptime}} {{.ShortDuration}} | " >> .tmp/example-1.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format="{{.LineNumber3}} "' > .tmp/example-2.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr --sleep-max=0 | prefix --format="{{.LineNumber3}} " >> .tmp/example-2.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format=">>> "' > .tmp/example-3.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr --sleep-max=0 | prefix --format=">>> " >> .tmp/example-3.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format="{{SLOW_LINES}} up={{.ShortUptime}} | "' > .tmp/example-4.txt
	generate-fake-data --seed=4242 --lines=10 --sleep-max=1.5s --dict=lorem-ipsum --no-stderr | prefix --format="{{SLOW_LINES}} up={{.ShortUptime}} | " >> .tmp/example-4.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format="{{SHORT_DATE}} "' > .tmp/example-5.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr --sleep-max=1.5s | prefix --format="{{SHORT_DATE}} " >> .tmp/example-5.txt

	@command -v embedmd >/dev/null 2>&1 && embedmd -w README.md || echo "embedmd not found, skipping README update"

.DEFAULT_GOAL := all