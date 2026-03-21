set dotenv-load

version := `cat VERSION`
build_time := `date +%s`

run: build
    ./bin/web-starter

build:
    rm -rf bin/
    go generate ./...
    go build -ldflags "-X main.Version={{version}} -X main.BuildTime={{build_time}}" -o bin/web-starter .

verify:
    go mod tidy
    go fmt ./...
    go vet ./...
    go test ./...
    go build ./...

start-dev-env:
    docker compose -f extra/dev/docker-compose.yml  down -v
    docker compose -f extra/dev/docker-compose.yml  up --build
