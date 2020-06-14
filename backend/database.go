package backend

import (
	"context"
)

func CreatNewUser(email string, name string) {
	//Get new connection
	ctx := context.Background()
	rdb := newClient()
	rdb.Do(ctx, "HMSET", "user:"+email, "name", name, "email", email)
}


