package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
)

func string(client *redis.Client) {
	key := "name"
	val := "大脸猫"
	err := client.Set(key, val, 30*time.Second).Err()
	checkError(err)

	v2, err := client.Get(key).Result()
	checkError(err)
	fmt.Println(v2)

	//client.Del(key)
}
func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "39.106.49.188:6379",
		Password: "liujun", // no password set
		DB:       0,        // use default DB
	})
	//ctx := context.TODO()
	string(client)

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
