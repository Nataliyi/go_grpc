package main

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var sid int = 12

	for pid := 0; pid <= 1000; pid++ {
		key := strconv.Itoa(pid) + ":" + strconv.Itoa(sid)
		err := rdb.Set(ctx, key, pid, 0).Err()
		if err != nil {
			panic(err)
		}
	}
}
