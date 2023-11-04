FROM golang:1.21

COPY . src/app

WORKDIR src/app

RUN make build

RUN chmod +x ./scripts/wait-for-it.sh

EXPOSE 8080

CMD ["./urlShortener"]