FROM golang:1.19

COPY . src/app

WORKDIR src/app

RUN make build

RUN chmod +x ./wait-for-it.sh

EXPOSE 8080

CMD ["./urlShortener"]