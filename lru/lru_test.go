package lru

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Id string
}

func (d TestStruct) Len() int {
	return len(d.Id)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", TestStruct{"1234"})
	if v, ok := lru.Get("key1"); !ok || v.(TestStruct).Id != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, TestStruct{Id: v1})
	lru.Add(k2, TestStruct{Id: v2})
	lru.Add(k3, TestStruct{Id: v3})

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", TestStruct{"123456"})
	lru.Add("k2", TestStruct{"k2"})
	lru.Add("k3", TestStruct{"k3"})
	lru.Add("k4", TestStruct{"k4"})

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

func TestAdd(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", TestStruct{"1"})
	lru.Add("key", TestStruct{"111"})

	if lru.nbytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nbytes)
	}
}
