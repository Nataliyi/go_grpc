package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var default_cluster = 11

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var sid int = 12

	for pid := 0; pid <= 1000; pid++ {
		var key = strconv.Itoa(pid) + ":" + strconv.Itoa(sid)
		cluster, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			cluster = default_cluster
		} else if err != nil {
			panic(err)
		} else {
			fmt.Println(key, cluster)
		}
	}

	fmt.Println("Done")
}
