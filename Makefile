DEFAULT_GOAL := build

.PHONY: check style build run dev docs

check:
	go mod tidy && go mod verify && go vet ./...

style:
	goimports-reviser -rm-unused -separate-named -set-alias \
	-imports-order std,general,company,project,blanked,dotted \
	-project-name bobshop \
	-format ./...

build: check style
	go build -o bin/app cmd/server/main.go

run: build
	./bin/app

dev:
	air

docs:
	swag init -g cmd/server/main.go