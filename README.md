# Currency Service


[![build](https://github.com/edmarfelipe/currency-service/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/edmarfelipe/currency-service/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/edmarfelipe/currency-service/graph/badge.svg?token=oZfYQLbFHH)](https://codecov.io/gh/edmarfelipe/currency-service)

## Description

This is a simple service that converts currency using the [Currency API](https://currencyapi.com/).

## Table of Contents
- [Overview](#overview)
  * [Architecture](#architecture) 
  * [Caching strategy](#caching-strategy)
- [Installing](#installing)
- [Usage](#usage)
- [REST API](#rest-api)
- [Architecture](#architecture)
- [Caching strategy](#caching-strategy)

### Overview

#### Architecture

```mermaid
C4Container
    Container_Ext(client, "Clients", "Web App, Mobile, Service, etc", "External client that uses the API") 

    Boundary(c1, "Currency Service") {
        Container(api, "API Application", "Container: Go", "API that converts currency")
        ContainerDb(cache, "Cache", "Container: Redis", "Stores the currency rates")
    }

    System_Ext(currency_api, "Currency API",  "External API that provides the currency rates")

    Rel(client, api, "Uses", "sync, JSON/HTTPS")
    Rel_Back(cache, api, "Reads from and writes to", "sync, TCP/6379")
    Rel(api, currency_api, "Reads from", "sync, JSON/HTTPS")

    UpdateRelStyle(cache, api,$textColor="#ced4da", $lineColor="#ced4da", $offsetX="-50", $offsetY="25")
    UpdateRelStyle(api, currency_api,$textColor="#ced4da", $lineColor="#ced4da",  $offsetX="10", $offsetY="-30")
    UpdateRelStyle(client, api, $textColor="#ced4da", $lineColor="#ced4da", $offsetX="5", $offsetY="-40")
```

#### Caching strategy

We want to always get from the cache as match as possible, to do so, we will update the cache in the background, while we return the value from the cache.
```mermaid
flowchart TD
    NewRequest((New Request)) -->
    IsOnCache{Is in the cache?}
    IsAsync{Is Async?}
    IsCacheFinishing{Cache is almost expire?}
    Request[Request from API]
    Return[Return value]

    IsOnCache -->|Yes| IsCacheFinishing
    IsCacheFinishing-.->|Yes|Request
    IsCacheFinishing-->|No|Return
    IsCacheFinishing-->|Yes|Return
    IsOnCache -->|No| Request
    Request--> SaveInCache
    SaveInCache[Save in the cache] --> IsAsync
    IsAsync -->|No|Return
    IsAsync -->|Yes|End
    
```

### Installing

To work with this project, you need to have installed:
* [Go 1.20](https://golang.org/doc/install)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/)

### Usage

#### Running the project

```shell
docker-compose up -d
```

#### Building the project

```shell
make build
```

#### Running the tests

```shell
make test
```

### REST API

[Open API Definition](./swagger.yaml)

| Description   | Verb   | Path                            |
|---------------|--------|---------------------------------|
| Convert Value | GET    | /api/convert/{currency}/{value} |
| Metrics       | GET    | /api/metrics                    |
| Ready         | GET    | /api/ready                      |

Example of request:

```shell
curl --request GET \
  --url http://localhost:3000/api/convert/BRL/543.34
```

### Project Key Features

- Structure Logging with Request ID
- Graceful Shutdown
- Healthcheck route that waits for the startup to finish
- Usage of Redis [Client-side Caching](https://redis.io/topics/client-side-caching)
- Prometheus metrics with latency and request count
- Tiny docker image with distroless base image


