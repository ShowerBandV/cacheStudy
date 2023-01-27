package cache

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

var DB = map[string]string{
	"Tom":  "123",
	"Jack": "456",
	"Sam":  "7890",
}

func TestNewHTTPPool(t *testing.T) {
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := DB[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := NewHTTPPool(addr)
	log.Println("cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
func createGroup() *Group {
	return NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := DB[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, ca *Group) {
	peers := NewHTTPPool(addr)
	peers.Set(addrs...)
	ca.RegisterPeers(peers)
	log.Println("cache server is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, ca *Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := ca.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("api server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func TestDistribute(t *testing.T) {
	port := []int{
		8001, 8002, 8003,
	}
	ca := createGroup()
	apiAddr := "http://localhost:9999"
	go startAPIServer(apiAddr, ca)

	for _, v := range port {
		addrMap := map[int]string{
			8001: "http://localhost:8001",
			8002: "http://localhost:8002",
			8003: "http://localhost:8003",
		}

		var addrs []string
		for _, v := range addrMap {
			addrs = append(addrs, v)
		}

		startCacheServer(addrMap[v], addrs, ca)
	}
}

func TestConcurrency(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		TestDistribute(t)
	}()
	time.Sleep(1 * time.Second)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := http.Client{}
			_, err := client.Get("http://localhost:9999/api?key=Tom")
			time.Sleep(10 * time.Millisecond)
			if err != nil {
				return
			}
		}()
	}

	wg.Wait()
}
