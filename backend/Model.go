package backend

type User struct {
	Email string `json: "email"`
	Password string `json: "password"`
}

type UserInfo struct {
	Name  string `redis:"name" json: "name"`
	Email string `redis:"email" json: "email"`
	Phone   string `redis:"phone" json: "phone"`
}
