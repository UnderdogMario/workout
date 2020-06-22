package backend

type User struct {
	Email string `json: "email"`
	Password string `json: "password"`
}

type UserInfo struct {
	Name  string `redis:"name"`
	Email string `redis:"email"`
	Phone   string `redis:"phone"`
}