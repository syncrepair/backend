package redis

import (
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func New(url string) *Redis {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic("error parsing redis url: " + err.Error())
	}

	var rd Redis

	rd.Client = redis.NewClient(opt)

	return &rd
}

func (r *Redis) Close() {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			panic("error closing redis client: " + err.Error())
		}
	}
}
