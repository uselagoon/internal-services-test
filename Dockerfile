FROM golang:alpine

WORKDIR /go-dbaas

ADD . .

RUN go get github.com/joho/godotenv

ENTRYPOINT go build

ENV MARIADB_PASSWORD=api \
    MARIADB_USER=api \
    MARIADB_DATABASE=infrastructure

EXPOSE 8080

CMD ["go-dbaas"]