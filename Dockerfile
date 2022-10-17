FROM golang:alpine

WORKDIR /go-dbaas

ADD . .

RUN go get github.com/joho/godotenv

RUN go build && chmod +x ./go-dbaas

ENV MARIADB_PASSWORD=api \
MARIADB_USER=api \
MARIADB_DATABASE=infrastructure \
POSTGRES_USER=pqgotest \
POSTGRES_PASSWORD=pqgotest \
POSTGRES_DB=infrastructure

EXPOSE 8080

CMD sleep 10 && ./go-dbaas