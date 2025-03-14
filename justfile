default:
    @just --list

# build fandogh binary
build:
    go build -o fandogh ./cmd/fandogh

# update go packages
update:
    @cd ./cmd/fandogh && go get -u

# run golangci-lint
lint:
    golangci-lint run -c .golangci.yml
