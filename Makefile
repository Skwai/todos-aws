build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/get functions/get.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create functions/create.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/index functions/index.go