

build:
		go build -o urlShortener -v ./cmd/main.go
run:
		go run ./cmd/main.go
clean:
		rm -rf ./urlShortener