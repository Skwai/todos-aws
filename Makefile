build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/get functions/get.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post functions/post.go