

build:
		go build -o urlShortener -v ./cmd/api/main.go
run:
		go run ./cmd/api/main.go
clean:
		rm -rf ./urlShortener