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
func initConnection() redis.Conn{
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
func ValidateUserInformation(email string, password string) (bool, UserInfo){
	conn := initConnection()
	defer conn.Close()
	var userInfo = UserInfo{}
	result, err := redis.String(conn.Do("HGET", "user:"+email, "password"))
	if err != nil {
		return false, userInfo
	} else if result == ""  {
		fmt.Println("Account Don't Exist")
		return false, userInfo
	} else {
		if result == password {
			fmt.Println("Login Success")
			userInfo := getUserInfo(email, conn, userInfo)
			return true, userInfo
		} else {
			fmt.Println("Login Fail")
			return false, userInfo
		}
	}
}

// this takes in a redis connection and a email to generate a session-id for that user
func CreateSessionToken(email string) string{
	sessionToken := uuid.NewV4().String()
	_, err := initConnection().Do("SETEX", sessionToken, "99999", email)
	if err != nil {
		return ""
	}
	return sessionToken
}

// create new user
func CreateNewUser(email string, password string) bool {
	if email == "" || password == "" {
		return false
	}
	conn := initConnection()
	defer conn.Close()
	redisEmail, _ := redis.String(conn.Do("HGET", "user:"+email, "email"))

	if redisEmail != "" {
		fmt.Println("user already exist")
		return false
	}
	conn.Do("HMSET", "user:"+email, "email", email, "password", password)
	return true
}

// this will return all user information
func getUserInfo(email string, conn redis.Conn, userInfo UserInfo) UserInfo{
	v, err := redis.Values(conn.Do("HGETALL", "user:"+email))
	if err != nil {
		fmt.Println(err)
	}
	if err := redis.ScanStruct(v,&userInfo); err != nil {
		fmt.Println(err)
	}
	return userInfo
}

// this will validate whether the session expire
func ValidateSessionID(sid string) bool {
	conn := initConnection()
	defer conn.Close()
	v, _ := redis.String(conn.Do("GET", sid))
	if v == "" {
		fmt.Println("Session ID invalid")
		return false
	} else {
		return true
	}
}

// this takes the new userInfo and update the database
func UserProfileUpdate(userInfo UserInfo) {
	conn := initConnection()
	defer conn.Close()
	conn.Do("HMSET", redis.Args{}.Add("user:"+userInfo.Email).AddFlat(&userInfo)...)
}
