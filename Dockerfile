FROM golang:alpine

WORKDIR /go-dbaas

ADD . .

RUN go get github.com/joho/godotenv

RUN go build && chmod +x ./go-dbaas

ENV MARIADB_PASSWORD=lagoon \
MARIADB_USER=lagoon \
MARIADB_DATABASE=lagoon \
MARIADB_HOST=mariadb \
POSTGRES_USER=lagoon \
POSTGRES_PASSWORD=lagoon \
POSTGRES_DATABASE=lagoon \
POSTGRES_HOST=postgres

EXPOSE 3000

CMD sleep 10 && ./go-dbaas