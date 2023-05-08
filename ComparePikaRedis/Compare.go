package ComparePikaRedis

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
)

func Hello() {

}

func init() {

}

type Comparator struct {
	pika          *redis.Client
	redis         *redis.Client
	pikaPipeline  redis.Pipeliner
	redisPipeline redis.Pipeliner
	context       *context.Context
}

func (comp *Comparator) Init(pika, redis *redis.Client, context *context.Context) {
	comp.pika = pika
	comp.redis = redis
	comp.context = context
}
func (comp *Comparator) Pipeline() {
	comp.pikaPipeline = comp.pika.Pipeline()
	comp.redisPipeline = comp.redis.Pipeline()
}

func (comp *Comparator) AddCmd(functioname string, key string, values ...interface{}) {
	switch functioname {
	case "LPush":
		comp.pikaPipeline.LPush(*comp.context, key, values[0].([]string))
		comp.redisPipeline.LPush(*comp.context, key, values[0].([]string))
	case "LTrim":
		comp.pikaPipeline.LTrim(*comp.context, key, values[0].(int64), values[1].(int64))
		comp.redisPipeline.LTrim(*comp.context, key, values[0].(int64), values[1].(int64))
	case "LRange":
		comp.pikaPipeline.LRange(*comp.context, key, values[0].(int64), values[1].(int64))
		comp.redisPipeline.LRange(*comp.context, key, values[0].(int64), values[1].(int64))
	}
}

func (comp *Comparator) ExecCompare() bool {
	pikaCmds, err := comp.pikaPipeline.Exec(*comp.context)
	if err != nil {
		panic(err)
	}
	redisCmds, err := comp.redisPipeline.Exec(*comp.context)
	if err != nil {
		panic(err)
	}
	if len(redisCmds) != len(pikaCmds) {
		panic("two redis client cmd number not equal ")
	}
	for i := range redisCmds {
		if IgnoreRedisResTime(redisCmds[i].String()) != IgnoreRedisResTime(pikaCmds[i].String()) {
			fmt.Println("redisCmds[i]:", redisCmds[i].FullName())
			fmt.Println("res:", redisCmds[i].String())
			fmt.Println("pikaCmds[i]:", pikaCmds[i].FullName())
			fmt.Println("res:", pikaCmds[i].String())
			return false
		}
	}
	return true
}

/**
忽略掉结果最后的时间
*/
func IgnoreRedisResTime(res string) string {
	lastColon := strings.LastIndex(res, ":")
	if lastColon < 0 {
		return res
	}
	return res[:lastColon]
}
