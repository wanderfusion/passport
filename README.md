# passport
[![Docker](https://github.com/wanderfusion/passport/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/wanderfusion/passport/actions/workflows/docker-publish.yml)
## Running
```
docker build -t passport:latest .
docker run -p 8080:8080 -v $(pwd)/config.yml:/app/config.yml passport:latest
```

## DB
- Postgres 15
- Migrations using https://github.com/jackc/tern
