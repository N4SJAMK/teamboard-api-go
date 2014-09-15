# Teamboard API

```
teamboard-api
```

## Dependencies

Requires `golang` and `MongoDB`.

## Installation

Install with:
```
go get github.com/N4SJAMK/teamboard-api
```

Run with:
```
teamboard-api
```
You can provide a alternate host and port using:
```
teamboard-api -bind every.day:4200
```
You can provide the following `MongoDB` config as env. vars:
- `MONGODB_URL` defaults to `mongodb://localhost`
- `MONGODB_NAME` defaults to `teamboard-dev-go`



