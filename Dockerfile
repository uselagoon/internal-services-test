FROM golang:1.23-alpine

WORKDIR /internal-services-test

ADD . .

RUN go get github.com/joho/godotenv

RUN go build && chmod +x ./internal-services-test

ENV STORAGE_LOCATION='/app/files'

EXPOSE 3000

CMD sleep 10 && ./internal-services-test
