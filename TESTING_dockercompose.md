Docker Compose test internal services test
==========================================

This is a docker-compose testing suite for the internal services:

Start up tests
--------------

Run the following commands to get up and running with this example.

```bash
# should remove any previous runs and poweroff
docker compose down --volumes --remove-orphans

# should start up our services successfully
docker compose build && docker compose up -d
```

Verification commands
---------------------

Run the following commands to validate things are rolling as they should.

```bash
# Ensure services are ready to connect
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://mariadb-10-5:3306 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://mariadb-10-11:3306 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://postgres-12:5432 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://postgres-16:5432 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://opensearch-2:9200 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://mongo-4:27017 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://redis-6:6379 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://redis-7:6379 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://solr-8:8983 -timeout 1m
docker run --rm --net internal-services-test_default jwilder/dockerize dockerize -wait tcp://solr-9:8983 -timeout 1m

# commons should be running Alpine Linux
docker compose exec -T commons sh -c "cat /etc/os-release" | grep "Alpine Linux"

# mariadb-10-5 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mariadb?service=mariadb-10-5" | grep "SERVICE_HOST=10.5"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mariadb?service=mariadb-10-5" | grep "LAGOON_TEST_VAR=internal-services-test"

# mariadb-10-11 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mariadb?service=mariadb-10-11" | grep "SERVICE_HOST=10.11"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mariadb?service=mariadb-10-11" | grep "LAGOON_TEST_VAR=internal-services-test"

# mysql-8-0 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mysql?service=mysql-8-0" | grep "SERVICE_HOST=8.0"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mysql?service=mysql-8-0" | grep "LAGOON_TEST_VAR=internal-services-test"

# mysql-8-4 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mysql?service=mysql-8-4" | grep "SERVICE_HOST=8.4"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mysql?service=mysql-8-4" | grep "LAGOON_TEST_VAR=internal-services-test"

# postgres-12 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/postgres?service=postgres-12" | grep "SERVICE_HOST=PostgreSQL 12"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/postgres?service=postgres-12" | grep "LAGOON_TEST_VAR=internal-services-test"

# postgres-16 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/postgres?service=postgres-16" | grep "SERVICE_HOST=PostgreSQL 16"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/postgres?service=postgres-16" | grep "LAGOON_TEST_VAR=internal-services-test"

# opensearch-2 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/opensearch?service=opensearch-2" | grep "SERVICE_HOST=opensearch-2"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/opensearch?service=opensearch-2" | grep "LAGOON_TEST_VAR=internal-services-test"

# mongo-4 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mongo?service=mongo-4" | grep "SERVICE_HOST=mongo-4"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mongo?service=mongo-4" | grep "LAGOON_TEST_VAR=internal-services-test"

# redis-6 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/redis?service=redis-6" | grep "SERVICE_HOST=redis-6"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/redis?service=redis-6" | grep "LAGOON_TEST_VAR=internal-services-test"

# redis-7 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/redis?service=redis-7" | grep "SERVICE_HOST=redis-7"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/redis?service=redis-7" | grep "LAGOON_TEST_VAR=internal-services-test"

# solr-8 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/solr?service=solr-8" | grep "SERVICE_HOST=solr-8"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/solr?service=solr-8" | grep "LAGOON_TEST_VAR=internal-services-test"

# solr-9 should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/solr?service=solr-9" | grep "SERVICE_HOST=solr-9"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/solr?service=solr-9" | grep "LAGOON_TEST_VAR=internal-services-test"

# persistent storage should be able to read/write data
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/storage?path=/app/files" | grep "STORAGE_PATH=/app/files/storage.txt"
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/storage?path=/app/files" | grep "LAGOON_TEST_VAR=internal-services-test"

# Incorrect service should be caught & error output
docker compose exec -T commons sh -c "curl -kL http://go-web:3000/mariadb?service=incorrect-service" | grep "mariadb is not a compatible driver with service: incorrect-service"
```

Destroy tests
-------------

Run the following commands to trash this app like nothing ever happened.

```bash
# should be able to destroy our services with success
docker compose down --volumes --remove-orphans
```
