# Lagoon Internal Services Test

This project was designed to perform a functional test on a given native Lagoon service.

It works by
* taking a service type (e.g. solr, redis, postgres etc)
* taking a service name (e.g. solr-9, redis-6, postgres-15) that matches the docker compose service name
* writing a set of predefined variables into the selected services data store
* reading this data back to a web page for easy validation

## Supported services

Internal Services Test currently supports
* [MariaDB](https://docs.lagoon.sh/docker-images/mariadb/) - use `mariadb`
* [MongoDB](https://docs.lagoon.sh/docker-images/mongodb/) - use `mongodb`
* [OpenSearch](https://docs.lagoon.sh/docker-images/opensearch/) - use `opensearch`
* [PostgreSQL](https://docs.lagoon.sh/docker-images/postgres/) - use `postgres`
* [Redis](https://docs.lagoon.sh/docker-images/redis/) - use `redis`
* [Solr](https://docs.lagoon.sh/docker-images/solr/) - use `solr`
* [Persistent Storage](https://docs.lagoon.sh/concepts-basics/docker-compose-yml/#persistent-storage/) - use `storage`

As we add more service types to Lagoon, this list may expand

## Usage in Docker Compose locally

This project can be cloned, and the services started with:

```
docker compose build
docker compose up -d
```

You can then access the main test interface at `http://0.0.0.0:3000/` - note that there is no default page (yet)

To check a service, use this format http://0.0.0.0:3000/{service-type}?service={service-name}

```
internal-services-test$ curl -kL http://0.0.0.0:3000/opensearch?service=opensearch-2
"SERVICE_HOST=opensearch-2"
"LAGOON_ENVIRONMENT_TYPE=development"
"LAGOON_GIT_SHA=SHA256"
"LAGOON_TEST_VAR=internal-services-test"
```
and to check a storage path
```
internal-services-test$ curl -kL http://0.0.0.0:3000/storage?path=/app/files
"STORAGE_PATH=/app/files/storage.txt"
"LAGOON_GIT_SHA=SHA256"
"LAGOON_ENVIRONMENT_TYPE=development"
"LAGOON_TEST_VAR=internal-services-test"
```

### Additional or custom variables

In the docker-compose.yml file, you can update the environment variables in the `web` service and restart the docker-compose service. This will provide the updated variables to the service next time it is called.

## Usage in Lagoon

This project can also be cloned and deployed into a Lagoon cluster, using the same patterns as above for verification.

Note that by default, all services in the docker-compose.yml will be deployed



