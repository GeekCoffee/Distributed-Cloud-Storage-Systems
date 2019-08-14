package redis

import(
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

const(
	HostRedis = "127.0.0.1:6379"  //redis-server的socket
	PassRedis = "abc5518988"     //redis的密码
)

var(
	connPool *redis.Pool
)

//NewRedisPool: 创建一个redis连接池
func NewRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,              //设置在连接池中最大闲置的可用连接数
		MaxActive:   30,              //在一段时间内，连接池中的最大的同时连接数
		IdleTimeout: 5 * time.Minute, //某条application到redis的连接，超过5分钟没有使用，就会被连接池回收
		Dial: func() (redis.Conn, error) { //dial是拨号的意思，也就是创建一个conn连接到连接池中
			//1.打开连接
			c, err := redis.Dial("tcp", HostRedis)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			//2.访问认证
			if _, err := c.Do("auth", PassRedis); err != nil {
				c.Close()
				log.Fatal(err)
				return nil, err
			}
			return c, nil
		},

		//TestOnBorrow: 用于在连接被应用层程序使用之前，做定时检查可用连接的健康状态
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute { //小于1分钟不做检测
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}


func init() {
	//生成连接池
	connPool = NewRedisPool()
}

func RedisPool() *redis.Pool {
	return connPool
}