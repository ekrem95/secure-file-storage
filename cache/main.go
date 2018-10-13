package cache

import "github.com/go-redis/redis"

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "pass",
	DB:       0, // use default DB
})

// pong, err := client.Ping().Result()
// fmt.Println(pong, err)
