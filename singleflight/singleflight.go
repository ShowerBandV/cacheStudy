package singleflight

import (
	"fmt"
	"sync"
	"time"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// 限制同类请求的重复调用
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		//测试时取消注释
		fmt.Println("本次请求", key, "被限制重复调用", time.Now())
		c.wg.Wait()
		return c.val, c.err
	}
	//测试时取消注释
	fmt.Println("本次请求", key, "未被限制", time.Now())
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
