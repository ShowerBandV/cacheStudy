package singleflight

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	var g Group
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}

func TestConcurrency(t *testing.T) {
	wg := sync.WaitGroup{}
	var g Group
	for i := 0; i < 10000; i++ {
		go g.Do("key1", func() (interface{}, error) {
			random := rand.Intn(10) + 1
			time.Sleep(time.Duration(random) * time.Millisecond)
			return "bar", nil
		})
		go g.Do("key2", func() (interface{}, error) {
			random := rand.Intn(10) + 1
			time.Sleep(time.Duration(random) * time.Millisecond)
			return "bar", nil
		})
	}
	wg.Wait()
}
