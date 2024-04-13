package client

import "time"

type Client struct {
	Id         int64
	Name       string
	ApiKey     string
	LastUpdate time.Time
}
