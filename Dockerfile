FROM golang:1.19

COPY . go/src/app

WORKDIR go/src/app

RUN go build -o urlShortener ./cmd/api/main.go

RUN chmod +x ./wait-for-it.sh

EXPOSE 8080

CMD ["./urlShortener"]