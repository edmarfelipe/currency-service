package cache

import (
	"github.com/rueian/rueidis"
)

// Open opens connection with the redis server
func Open(addr string, username string, password string) (rueidis.Client, error) {
	return rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{addr},
		Username:    username,
		Password:    password,
		SelectDB:    0,
	})
}
