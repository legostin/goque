package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func ExampleNewClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "172.17.0.3:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func main() {
	fmt.Printf("hello, world\n")
	ExampleNewClient()
}
