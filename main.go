package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"goredisclient/ComparePikaRedis"
	"strconv"
	"sync"
)

func getStrings(prefix string, length int) []string {
	values := make([]string, 0)
	for i := 0; i < length; i++ {
		values = append(values, prefix+".val"+strconv.Itoa(i))
	}
	return values
}

var wg sync.WaitGroup
var mu sync.Mutex

func Parallel(comp *ComparePikaRedis.Comparator) {
	mu.Lock()
	defer mu.Unlock()
	cmdcot := 6

	for k := 1; k < cmdcot; k++ {
		key := fmt.Sprintf("list%d", k)
		vals := getStrings(key, 500)
		comp.AddCmd("LPush", key, vals)
		comp.AddCmd("LRange", key, int64(0), int64(200-1))
		comp.AddCmd("LTrim", key, int64(200), int64(-1))
		comp.ExecCompare()
	}
	wg.Done()
}
func main() {

	comp := ComparePikaRedis.Comparator{}
	redisClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	pikaClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:9221",
	})

	ctx := context.Background()
	comp.Init(pikaClient, redisClient, &ctx)
	threadnumber := 100
	wg.Add(threadnumber)
	for i := 0; i < threadnumber; i++ {
		go Parallel(&comp)
	}
	wg.Wait()
}
