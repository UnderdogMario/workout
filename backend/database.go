package backend

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func newClient() *redis.Client {
	address := goDotEnvVariable("ADDRESS")
	password := goDotEnvVariable("PASSWORD")
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
	return rdb
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func CreatNewUser(email string, name string) {
	//Get new connection
	ctx := context.Background()
	rdb := newClient()
	rdb.Do(ctx, "HMSET", "user:"+email, "name", name, "email", email)
}

func ValidateUserInformation(email string) {
	ctx := context.Background()
	rdb := newClient()
	result := rdb.Do(ctx, "HGET", "user:"+email, "name")
	fmt.Println(result)
}


