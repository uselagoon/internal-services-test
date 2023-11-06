FROM golang:1.21-alpine

WORKDIR /internal-services-test

ADD . .

RUN go get github.com/joho/godotenv

RUN go build && chmod +x ./internal-services-test

ENV SOLR_HOST=solr \
    REDIS_HOST=redis \
    OPENSEARCH_HOST=opensearch-2

EXPOSE 3000

CMD sleep 10 && ./internal-services-test
