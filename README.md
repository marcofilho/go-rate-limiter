# Go Rate Limiter

This project implements a **Rate Limiter** in Go, which can be configured to limit the maximum number of requests per second based on an IP address or an access token.

## Project Structure

```plaintext
go-rate-limiter/
├── cmd/
│   ├── main.go
│   └── .env
├── configs/
│   └── config.go
├── internal/
│   ├── limiter/
│   │   ├── limiter.go
│   │   ├── storage.go
│   │   └── redis_storage.go
│   ├── middleware/
│   │   └── middleware.go
│   └── server/
│       └── server.go
├── tests/
    └── token/
│   │   ├── test_script_rate_limiter_token.go   # Script for testing rate limiting with a valid token
│   └── test_script_rate_limiter.go # Script for testing rate limiting
├── Dockerfile           # Dockerfile for building the project
├── docker-compose.yml   # Docker Compose configuration
└── README.md            # Project documentation
```


## Features

- Request limiting by **IP** or **Access Token**.
- Configurable limits and block durations via `.env` file.
- Data persistence in Redis.
- Middleware for seamless integration with web servers.
- Support for different limits for specific tokens.
- Docker and Docker Compose for easy deployment.

---

## Prerequisites

- [Go](https://golang.org/) (version 1.20 or higher)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## How to Use

### 1. Build and Run with Docker

#### **1.1 Build the Docker Image**
In the root directory of the project, run:

```bash
docker build -t go-rate-limiter .
```

#### **1.2 Run the Container**
Run the container with:

```bash
docker run -d --name go-rate-limiter -p 8080:8080 --env-file ./cmd/.env go-rate-limiter
```

- The server will be available at `http://localhost:8080`.

---

### 2. Use Docker Compose

If you are using Redis, you can run the server and Redis together with Docker Compose.

#### **2.1 Start the Services**
Run:

```bash
docker-compose up --build
```

---

### 3. Test the System

#### **3.1 `/ping` Endpoint**
Send requests to the `/ping` endpoint:

```bash
curl http://localhost:8080/ping
```

#### **3.2 Test with Token**
Send requests with the `API_KEY` header:

```bash
curl -X GET http://localhost:8080/ping -H "API_KEY: my-token"
```

#### **3.3 Test Rate Limiting**
Use the test script to send multiple requests:

```bash
go run tests/token/test_script_rate_limiter_token.go 15
```

---