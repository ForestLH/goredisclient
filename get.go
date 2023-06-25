//libadd.go
package main

import (
	"os"
	"fmt"
	"sync"
	"time"
	"github.com/go-redis/redis"
)

func do_worker(master_addr, master_pass string, worker_num int, sum int64) {
	client := redis.NewClient(&redis.Options{
			Addr:     master_addr,
			Password: master_pass})
		
	var wg sync.WaitGroup
	for num:=0; num < worker_num; num++{
		wg.Add(1)
		go func() {
			for {
				pipe := client.Pipeline()
				pipe.LRange("test", 0, sum-1)
				pipe.LTrim("test", sum, -1)
				cmders, err := pipe.Exec()
				if err != nil {
					fmt.Println("pip.Exec() err", err)
					time.Sleep(time.Duration(1)*time.Second)
					continue
				}

				values,_ := cmders[0].(*redis.StringSliceCmd).Result()
				length := len(values)
				if length==0 {
					break
				}

				for _, value := range values {
					if value == "" {
						continue
					}
					
					fout, err := os.OpenFile("out.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
					if err != nil {
						continue
					}
					_, err = fout.WriteString(value + "\n")
					if err != nil {
						continue
					}
					fout.Close()
				}
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func main(){
	master_addr := "127.0.0.1:6379"
	master_pass := "123456"

	do_worker(master_addr, master_pass, 20, 50)
}
