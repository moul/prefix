GOPKG ?=	moul.io/prefix
DOCKER_IMAGE ?=	moul/prefix
GOBINS ?=	./cmd/prefix
NPM_PACKAGES ?=	.

include rules.mk

generate: install
	mkdir -p .tmp
	GO111MODULE=off go get moul.io/generate-fake-data

	echo 'foo@bar:~$$ prefix -h' > .tmp/usage.txt
	prefix -h 2>> .tmp/usage.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix' > .tmp/default.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr | prefix >> .tmp/default.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format="#{{.LineNumber3}} {{.Uptime}} {{.Duration}} | "' > .tmp/example-1.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr | prefix --format="#{{.LineNumber3}} {{.Uptime}} {{.Duration}} | " >> .tmp/example-1.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format="{{.LineNumber3}} "' > .tmp/example-2.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr --sleep-max=0 | prefix --format="{{.LineNumber3}} " >> .tmp/example-2.txt

	echo 'foo@bar:~$$ generate-fake-data | prefix -format=">>> "' > .tmp/example-3.txt
	generate-fake-data --seed=42 --lines=10 --dict=lorem-ipsum --no-stderr --sleep-max=0 | prefix --format=">>> " >> .tmp/example-3.txt

	embedmd -w README.md
	#rm -rf .tmp
