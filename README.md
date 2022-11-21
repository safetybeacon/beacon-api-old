# beacon-api


## Quick Start 

```bash
export POSTGRES_USER=<value>
export POSTGRES_PASSWORD=<value>
export POSTGRES_DATABASE=<value>
export POSTGRES_HOST=<value>
export PORT=<value>
```

```bash
go build -o server cmd/main.go
```

```bash
./server
```

## Swagger docs

you can find swagger docs on

```
http://localhost:<port>/v1/swagger/index.html
```

![swagger_docs](./swagger_docs.png)


## Docker

- run `docker-compose.yml` file

```bash
export db_password=<value>
```

```bash
docker-compose up --build -d
```

- to check the status of the pods

```bash
docker-compose ps
```

- you can test the API using swagger navigate to `http://localhost:8080/v1/swagger/index.html`
