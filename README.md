# Caching Gateway

A CLI-based proxy server project that uses Gin + Redis to forward requests to an origin server and cache the responses.
The purpose of this project is to return the response from the Redis cache without reaching the origin server when the same request comes in a second time.

---

## Purpose

The core functionality of this application is:

1. The user starts the server via CLI.
2. The proxy server begins listening on a specified port.
3. Incoming requests are looked up in the cache.
4. If found in cache, the response is returned directly from the cache.
5. If not found, the request is forwarded to the origin server.
6. The response from the origin is serialized and written to Redis.
7. On the same request arriving again, the response is served from cache.
8. All cached data in Redis can be cleared using `--clear-cache`.

---

## Technologies Used

* **Go** – primary language
* **Gin** – HTTP router and request handling
* **Redis** – storing cached responses
* **Docker** – can be used to run Redis and Prometheus servers
* **CLI flags** – server configuration
* **Prometheus** – metrics

---

## Overall Architecture

```text
Client
  ↓
Gin Router
  ↓
Proxy Handler
  ↓
Cache Layer (Redis)
  ↓
Origin Server
```

### How the Flow Works

#### On Cache MISS

* Request arrives.
* Proxy generates a cache key for that request.
* Looks up the key in Redis.
* If the key is not found, forwards the request to the origin server.
* Retrieves response body, headers, and status code.
* Writes the response to the cache.
* Returns to the client with the `X-Cache: MISS` header.

#### On Cache HIT

* Request arrives.
* The same cache key is found in Redis.
* Does not go to the origin server.
* The cached response is returned directly to the client.
* The `X-Cache: HIT` header is added for the client.

---

## File Structure

```text
caching-gateway/
├── cmd/
│   └── main.go
├── internal/
│   ├── cache/
│   │   └── redis.go
│   ├── config/
│   │   └── config.go
│   ├── model/
│   │   └── cache_item.go
│   ├── middleware/
│   │   ├── metrics/
│   │   │   └── metrics_middleware.go
│   │   └── rate_limiter.go
│   ├── proxy/
│   │   └── handler.go
│   ├── redis/
│   │   └── client.go
│   └── server/
│       └── server.go
├── go.mod
└── go.sum
```

---

## Flags

### `--port`

The port the proxy server will listen on.

Example:
```bash
--port 3000
```
In this case the server runs on `localhost:3000`

### `--origin`

The URL of the main server to which requests will be forwarded.

Example:
```bash
--origin http://example.com
```
The proxy forwards requests to this server.

### `--redis`

The address of the Redis server.

Example:
```bash
--redis localhost:6379
```

### `--clear-cache`

If this flag is provided, the program does not start as a server. It simply clears all cached data in Redis.

---

## Getting Started

### 1. Start Redis

With Docker:
```bash
docker run -d --name redis -p 6379:6379 redis
```

### 2. Run the project

```bash
go run ./cmd --port 3000 --origin http://example.com --redis localhost:6379 --redis-password mypassword
```

### 3. Send a test request

```bash
curl http://localhost:3000/products
```

The first response will be a MISS; the second is expected to be a HIT.

#### Customization

I have enabled GET, POST, and DELETE requests in line with my own server.
You can customize this to match your own server's needs.


#### Extra
This Repo serves as a solution to Roadmap.sh Caching Server Problem

