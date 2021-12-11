package internal

import (
	"net"
	"sync"
	"time"
)

var redisConnPool = &sync.Pool{
	New: func() interface{} {
		conn, err := (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 3 * time.Minute,
		}).Dial("tcp", "localhost:6379")
		if err != nil {
			panic(err)
		}
		return conn
	},
}

func Get() net.Conn {
	return redisConnPool.Get().(net.Conn)
}

func Put(conn net.Conn) {
	redisConnPool.Put(conn)
}
