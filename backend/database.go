package backend

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
)

// get a connection
func initCache() redis.Conn{
	address := goDotEnvVariable("ADDRESS")
	password := goDotEnvVariable("PASSWORD")
	conn, err := redis.Dial("tcp",address, redis.DialPassword(password))
	if err != nil {
		panic(err)
	}
	return conn
}

// read env data
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

//takes email and password, and return true if the password is correct
func ValidateUserInformation(email string, password string) bool{
	conn := initCache()
	result, err := redis.String(conn.Do("HGET", "user:"+email, "password"))
	if err != nil {
		log.Fatal(err)
		return false
	} else if result == ""  {
		fmt.Println("Account Don't Exist")
		return false
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

// this takes in a redis connection and a email to generate a session-id for that user
func CreateSessionToken(email string) string{
	sessionToken := uuid.NewV4().String()
	_, err := initCache().Do("SETEX", sessionToken, "1200", email)
	if err != nil {
		return ""
	}
	return sessionToken
}
