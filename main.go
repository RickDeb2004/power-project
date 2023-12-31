package main

import (
	"fmt"
	"net/http"
	cacher "power-project/cache"
	load "power-project/loadbalancer"
	serve "power-project/server"
	"time"
)

func main() {

	servers := []serve.Server{
		serve.NewServer("http://www.facebook.com", "http://www.facebook.com/health", 5*time.Second, 2*time.Second),
		serve.NewServer("http://www.bing.com", "http://www.bing.com/health", 5*time.Second, 2*time.Second),
		serve.NewServer("http://www.duckduckgo.com", "http://www.duckduckgo.com/health", 5*time.Second, 2*time.Second),
	}
	cache := cacher.NewCache()
	lb := load.NewLoadBalancer("8000", servers, load.WeightedRoundRobin, &cache)
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(w, r)
	}
	http.HandleFunc("/", handleRedirect)
	fmt.Printf("serving request at 'local host : %s'\n", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)

	// fmt.Println("Start Cache")
	// expiration := 2 * time.Second //Set expiration time to 2 seconds
	// cache := NewCache()
	// for _, word := range []string{"parrot", "avocardo", "tree", "potato", "tree"} {
	// 	cache.Check(word, expiration)
	// 	cache.Display()
	// }
	// time.Sleep(3 * time.Second) //Sleep for 3 seconds to allow some entries to expire
	// cache.RemoveExpired()       //Remove expired entries from the cache
	// cache.Display()
	// fmt.Printf("HitRate:%.2f%%\n", cache.GithitRate())
	// fmt.Printf("MissRate:%.2f%%\n", 100-cache.GithitRate())
	// fmt.Printf("TotalRate:%.2f%%\n", cache.TotalCount)
}
