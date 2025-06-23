# bobshop-api

## Debug in VSCode

Copy ***env/example.env -> env/config.env***

Create folder ***.vscode***, create file **launch.json** in ***.vscode***

```
{
	"version": "0.2.0",
	"configurations": [
		{
			"name": "bobshop-api",
			"type": "go",
			"request": "launch",
			"mode": "debug",
			"program": "${workspaceFolder}/main.go"
		}
	]
}
```

# Compile

```
GOOS=linux GOARCH="amd64" go build -o bobshop-api main.go
GOOS=darwin GOARCH="amd64" go build -o bobshop-api main.go
GOOS=windows GOARCH="amd64" go build -o bobshop-api_64bit.exe main.go
GOOS=windows GOARCH="386" go build -o bobshop-api_32bit.exe main.go
```

# Run

```
MacOS, Linux:
  ./bobshop-api -config=env/config.env

Windows:
  bobshop-api_32bit.exe -config=env/config.env
  bobshop-api_64bit.exe -config=env/config.env

Nohup Command
  Start:	nohup ./bobshop-api -config=env/config.env > bobshop-api.log &
  Stop:		pkill bobshop-api
```

## Deploy:

```
brew install make (macos)

apt install make (ubuntu)

make image
```

## Docker:

```
Build:		docker build -t bobshop-api .

Run:		docker run --name bobshop-api -dp 3000:3000 bobshop-api

Log:		docker logs -f bobshop-api
```

## Docker Compose:

```
Start:		docker-compose up -d

Stop:		docker-compose down

Log:		docker logs -f bobshop-api
```