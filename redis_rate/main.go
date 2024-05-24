package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// If you want to flush the database, uncomment the following line
	// err := rdb.FlushDB(ctx).Err()
	// if err != nil {
	// 	panic(err)
	// }

	limiter := redis_rate.NewLimiter(rdb)
	res, err := limiter.AllowN(ctx, "project:123", redis_rate.PerMinute(10), 11)
	if err != nil {
		panic(err)
	}
	fmt.Println("limit", res.Limit, "allowed", res.Allowed, "remaining", res.Remaining, "retry_after", res.RetryAfter)
}
