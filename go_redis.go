package main

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"fmt"
	"time"
)

var (
	server string = "localhost:6379"
	password string = "123456"
)

var pool *redis.Pool

func test(i int ) {
	c := pool.Get()
	defer c.Close()

	t := strconv.Itoa(i)
	c.Do("SETEX","foo"+t,20,i)

	reply, err := redis.Int(c.Do("GET","foo"+t))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reply)
	time.Sleep(1*time.Second)

}

func poolInit() (*redis.Pool) {
	return &redis.Pool{
		MaxIdle:3,
		IdleTimeout:240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",server)
			if err != nil {
				fmt.Println(err)
				return nil,err
			}
			if _, err := c.Do("AUTH",password);err != nil {
				c.Close()
				return nil, err
			}
			return  c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_,err := c.Do("ping");
			return err
		},
	}
}

func main() {
	pool = poolInit()
	for i:=0;i<1000000;i++{
		test(i)
	}
}