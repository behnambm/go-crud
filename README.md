### RUN

#### run tests
```bash
make test 
```

#### Create and seed the database
```bash
make init
```

#### run the server
```bash
make run
```


#### Note that the Postman collection is provided in the repo

# Docker 

```bash
make up
```

```bash 
docker compose exec app bash 
```

```bash 
./app.out --initdb 
```

`App:` localhost:8080

`Metrics:` localhost:9000

`Prometheus:` localhost:9090

`Grafana:` localhost:3000