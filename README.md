# Sliding Log Rate Limiter with Go and Redis

This repository implements a **Sliding Log Rate Limiter** using Go and Redis. The sliding log algorithm ensures precise rate limiting by maintaining a log of request timestamps and supports distributed systems with atomic operations using Redis Lua scripting.

## Features

- **Accurate Sliding Log Algorithm**: Dynamically calculates the rate limit window for precise control.
- **Distributed Support**: Shared rate limits across multiple application instances.
- **Atomic Operations**: Uses Lua scripting in Redis to ensure thread-safe operations.
- **Configurable**: Easily adjust rate limits and time windows via environment variables or configuration.

## How It Works

1. Each API request logs its timestamp in a Redis sorted set.
2. The Lua script:
    - Removes outdated timestamps outside the rate limit window.
    - Counts valid timestamps within the window.
    - Adds the new timestamp if the request is allowed.
3. The result determines whether the request can proceed or should be throttled.

## Prerequisites

- **Go 1.18+**
- **Redis 6.0+**

## Installation

1. Clone the repository:
    
    ```bash
    git clone https://github.com/<your-username>/sliding-log-rate-limiter.git
    cd sliding-log-rate-limiter
    ```
    
2. Install dependencies:
    
    ```bash
    go mod tidy
    ```
    
3. Ensure Redis is running and accessible.

## Configuration

Set the following environment variables or modify the configuration file:

- `REDIS_ADDR`: Redis server address (default: `localhost:6379`)
- `REDIS_PASSWORD`: Redis password (if required)
- `RATE_LIMIT`: Maximum requests allowed in the window (e.g., `100`)
- `WINDOW_SIZE`: Rate limit window size in seconds (e.g., `1`)

## Usage

1. Start the rate limiter:
    
    ```bash
    go run main.go
    ```
    
2. Use the following Lua script in Redis for the atomic rate limiting logic. The script is located in `scripts/ratelimiter.lua`:
    
    ```lua
    local pgname = KEYS[1]
    local now = tonumber(ARGV[1])
    local window = tonumber(ARGV[2])
    local limit = tonumber(ARGV[3])
    local clearBefore = now - window
    redis.call('ZREMRANGEBYSCORE', pgname, 0, clearBefore)
    local amount = redis.call('ZCARD', pgname)
    if amount < limit then
        redis.call('ZADD', pgname, now, now)
    end
    redis.call('EXPIRE', pgname, window)
    return limit - amount
    ```
    
3. Integrate the rate limiter into your application by invoking the `CheckRateLimit` function with parameters for the payment gateway, current time, and configuration.

## Benchmarking

Benchmark the rate limiter to test its throughput:

```bash
go test -bench=.
```

Sample results (local setup):

- Avg operation time: **55.37 μs**
- Throughput: **~18,000 requests/sec**

## Contributing

1. Fork the repository.
2. Create your feature branch: `git checkout -b feature-name`.
3. Commit your changes: `git commit -m 'Add feature'`.
4. Push to the branch: `git push origin feature-name`.
5. Open a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

- https://rafaeleyng.github.io/redis-pipelining-transactions-and-lua-scripts
- https://konghq.com/blog/how-to-design-a-scalable-rate-limiting-algorithm/
- https://rafaeleyng.github.io/redis-pipelining-transactions-and-lua-scripts
- https://redis.io/commands/zadd
- https://redis.io/commands/zcard
- https://redis.io/commands/zcount
- https://engagor.github.io/blog/2017/05/02/sliding-window-rate-limiter-redis/
- https://engagor.github.io/blog/2018/09/11/error-internal-rate-limit-reached/
- https://app.diagrams.net/#G1USTEf_sVbyi0ri0NjN_xUYmXraM7bLES
- https://medium.com/@saisandeepmopuri/system-design-rate-limiter-and-data-modelling-9304b0d18250
- https://www.figma.com/blog/an-alternative-approach-to-rate-limiting/
- https://gist.github.com/ptarjan/e38f45f2dfe601419ca3af937fff574d 
- https://stripe.com/blog/rate-limiters
- https://github.com/go-redis/redis
- https://github.com/abhirockzz/redis-geo.lua-golang/blob/master/redis-geo-lua-example.go
