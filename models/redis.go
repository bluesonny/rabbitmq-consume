package models

import (
	"github.com/gomodule/redigo/redis"
	"log"
	. "rabbitmq-consume/config"
	"time"
)

var RedisClient *redis.Pool

func init() {

	//初始化redis连接池

	RedisClient = &redis.Pool{
		MaxIdle:     50,
		MaxActive:   500,
		IdleTimeout: 30 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ViperConfig.Redis.Host)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", ViperConfig.Redis.Password); err != nil {
				_ = c.Close()
				log.Printf("redis连接%v", err.Error())
				return nil, err
			}

			log.Printf("redis连接成功了")

			return c, err
		},
	}
}
