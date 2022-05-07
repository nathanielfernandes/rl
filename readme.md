# rl - Ratelimiter

A simple and generic moving window ratelimiter for my go projects

#### usage

```go
import (
    rl "github.com/nathanielfernandes/rl"
)
...
//        4 calls per 5 second window
rlm := rl.NewRatelimitManager(4, 5000)
...
//                      ex. request host
if rlm.IsRateLimited("ANY STRING IDENTIFIER") {
    return // IDENTIFIER is ratelimited
} else {
    // IDENTIFIER is within the ratelimit
}
```
