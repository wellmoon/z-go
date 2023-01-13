package main

import (
	"fmt"
	"time"

	"github.com/wellmoon/z-go/zlock"
)

func main() {
	zl := zlock.New()
	zl.Lock(1)
	fmt.Println("lock 1 success first time")
	go func() {
		time.Sleep(time.Duration(2000) * time.Millisecond)
		zl.Unlock(1)
	}()
	zl.Lock(1)
	fmt.Println("lock 1 success second time")
}
