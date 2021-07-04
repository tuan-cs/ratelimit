## Simple ratelimiter

### Design
Based on the idea from http.HandleFunc, the limiter works as a middleware of request, allow chaining multiple limiter such as limiter based on total request and limiter based on number of request from individual phone number.
If one limiter denies request, return error to client.

### How to use
- First start the mock sms server
```
  go run server.go
```

- Change config in the examples/examples.go to suitable with requirement such as rps, refresh duration.
- Run examples to view results
```
  go run examples/examples.go
```