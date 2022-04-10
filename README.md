## Description

### Intro
The solution consists of main and worker applications:
- The main application provides 2 ports, one for grpc and one for rest APIs
- The worker just sit and wait for messages from the main application

### Stack
- Golang 1.16
- gRPC/gRPC-gateway
- AWS DynamoDB
- AWS SQS
- Docker-compose

### Project structure
- cmd/app/main.go - main application
- cmd/worker/main.go - worker application
- api/proto/* - contains a description of application proto files
- api/generated/* - contains generated proto files
- third_party/* - contains third party api's
- docker-compose.protoc.yml - contains a description of gRPC/Protocol buffer compiler container
- docker-compose.yml - contains a description application deployment
- internal/* - contains application codebase


## Starter

### Pre requirements
- Golang 1.16+
- Docker
- Docker-compose
- Make

### How to run application?
Run  `$ make run`

### How to run unit tests?
Run  `$ make unittest`

### How to generate proto api?
Run  `$ make proto`

## Usage

### Create answer via rest api:

```sh
$ curl -d '{"key":"name", "value":"John"}' -H "Content-Type: application/json" -X POST http://localhost:8000/v1/answers
```

### Update answer via rest api:

```sh
$ curl -d '{"key":"name", "value":"Sam"}' -H "Content-Type: application/json" -X PUT http://localhost:8000/v1/answers
```

### Get answer via rest api:

```sh
curl -H "Content-Type: application/json" -X GET http://localhost:8000/v1/answers?key=${KEY}
```

### Delete answer via rest api:

```sh
curl -H "Content-Type: application/json" -X DELETE http://localhost:8000/v1/answers?key=${KEY}
```

### Get answer history via rest api:

```sh
curl -H "Content-Type: application/json" -X GET http://localhost:8000/v1/answers/${KEY}/history
```