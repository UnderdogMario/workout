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
	address := GoDotEnvVariable("ADDRESS")
	password := GoDotEnvVariable("PASSWORD")
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

func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func CreatNewUser(email string, name string, password string) {
	//Get new connection
	ctx := context.Background()
	rdb := newClient()
	rdb.Do(ctx, "HMSET", "user:"+email, "name", name, "email", email, "password", password)
}

func ValidateUserInformation(email string, password string) bool{
	ctx := context.Background()
	rdb := newClient()
	result, err := rdb.HGet(ctx, "user:"+email, "password").Result()

	if err == redis.Nil {
		fmt.Println("Login Fail")
		return false
	} else if err != nil {
		panic(err)
	} else {
		if result == password {
			fmt.Println("Login Success")
			return true
		} else {
			fmt.Println("Login Fail")
			return false
		}
	}
}


