# passport · [![Docker](https://github.com/wanderfusion/passport/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/wanderfusion/passport/actions/workflows/docker-publish.yml)
## Running
### Local Build
```
docker build -t passport:latest .
docker run -p 8080:8080 -v $(pwd)/config.yml:/app/config.yml passport:latest
```

### GHCR
```
docker run -d -p 8081:8081 -v $(pwd)/config.yml:/app/config.yml ghcr.io/wanderfusion/passport:main
```

## DB
- Postgres 15
- Migrations using https://github.com/jackc/tern
