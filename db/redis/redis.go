package redis

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/log"
)

type HandleData func(d string)

type CRedis struct {
	Pool *redis.Pool
	Life int32
}

func (r *CRedis) InitDB() {
	r.Pool = Newfactory("")
	r.Life = conf.Config.Redisconfig.Life * 3600 * 24
}

func (r *CRedis) OnClose() {
	r.Pool.Close()
}

func Newfactory(name string) *redis.Pool {

	host := conf.Config.Redisconfig.Ip
	port := conf.Config.Redisconfig.Port
	password := conf.Config.Redisconfig.Password
	count := conf.Config.Redisconfig.MaxIdle
	pool := &redis.Pool{
		IdleTimeout: 180 * time.Second,
		MaxIdle:     int(count),
		MaxActive:   1024,
		Dial: func() (redis.Conn, error) {
			address := fmt.Sprintf("%s:%s", host, port)
			c, err := redis.Dial("tcp", address,
				redis.DialPassword(password),
			)
			if err != nil {
				log.Error("err:%v", err)
				return nil, err
			}

			return c, nil
		},
	}
	log.Debug("connnect redis success")
	return pool
}

func (r *CRedis) DelTable(table string) {
	c := r.Pool.Get()
	_, err := c.Do("del", table)
	if nil != err {
		log.Error("table:%v, error:%v", table, err)
	}
	c.Close()
}

func (r *CRedis) SetTableLifetime(table string, lifetime int32) error {
	c := r.Pool.Get()
	if lifetime > 0 {
		_, err := c.Do("EXPIRE", table, lifetime)
		if nil != err {
			log.Error("error table:%v", table)
		}
		c.Close()
		return err
	} else {
		_, err := c.Do("persist", table)
		if nil != err {
			log.Error("error table:%v", table)
		}
		c.Close()
		return err
	}
}

func (r *CRedis) SubChannel(name string, h HandleData) {
	c := r.Pool.Get()
	psc := redis.PubSubConn{c}
	psc.PSubscribe(name)
	for {
		switch v := psc.Receive().(type) {
		case redis.Subscription:
		//	fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case redis.Message: //单个订阅subscribe
			h(string(v.Data))
		case error:
			c.Close()
			return
		}

	}
	c.Close()
}

//string
func (r *CRedis) String_SetData(table, value string) error {
	c := r.Pool.Get()
	_, err := c.Do("set", table, value)
	if nil != err {
		log.Error("error table:%v, value:%v", table, value)
	}
	c.Close()
	return err
}

func (r *CRedis) String_SetDataMap(data map[string]string) error {
	c := r.Pool.Get()
	args := make([]interface{}, len(data)*2)
	i := 0
	for k, v := range data {
		args[i] = k
		args[i+1] = v
		i += 2
	}
	_, err := c.Do("mset", args...)
	if nil != err {
		log.Error("error data:%v, err:%v", data, err)
	}
	c.Close()
	return err
}

func (r *CRedis) String_GetData(table string) (string, error) {
	c := r.Pool.Get()
	value, err := redis.String(c.Do("get", table))
	if nil != err {
		log.Error("error table:%v， err:%v", table, err)
	}
	c.Close()
	return value, err
}

func (r *CRedis) String_GetData2(table interface{}) (string, error) {
	c := r.Pool.Get()
	value, err := redis.String(c.Do("get", table))
	if nil != err {
		log.Error("error table:%v， err:%v", table, err)
	}
	c.Close()
	return value, err
}

func (r *CRedis) String_GetDataMap(table []string) (map[string]string, error) {
	c := r.Pool.Get()
	args := make([]interface{}, len(table))
	for k, v := range table {
		args[k] = v
	}
	value, err := redis.Strings(c.Do("mget", args...))
	if nil != err {
		log.Error("error table:%v， err:%v", table, err)
	}
	values := make(map[string]string, len(table))
	for k, v := range value {
		if k < len(table) {
			values[table[k]] = v
		}
	}
	c.Close()
	return values, err
}

func (r *CRedis) String_GetDataMap1(table []int64) (map[int64]string, error) {
	c := r.Pool.Get()
	var args []interface{}
	var keys []int64
	for i := 0; i < len(table); {
		args = append(args, table[i])
		keys = append(keys, table[i])
		i += 2
	}

	value, err := redis.Strings(c.Do("mget", args...))
	if nil != err {
		log.Error("error table:%v", table)
	}
	values := make(map[int64]string, len(table)/2)

	for k, v := range value {
		if k < len(table)/2 {
			values[keys[k]] = v
		}
	}
	c.Close()
	return values, err
}

func (r *CRedis) String_GetDataMap2(table []string) ([]string, error) {
	c := r.Pool.Get()
	args := make([]interface{}, len(table))
	for k, v := range table {
		args[k] = v
	}
	value, err := redis.Strings(c.Do("mget", args...))
	c.Close()
	return value, err
}

func (r *CRedis) String_GetDataMap3(table []int64) ([]string, error) {
	c := r.Pool.Get()
	args := make([]interface{}, len(table))
	for k, v := range table {
		args[k] = v
	}
	value, err := redis.Strings(c.Do("mget", args...))
	c.Close()
	return value, err
}

func (r *CRedis) IncrData(table string) int64 {
	c := r.Pool.Get()
	value, err := redis.Int64(c.Do("INCR", table))
	c.Close()
	if nil != err {
		log.Error("error table:%v", table)
		return 0
	}
	return value
}

func (r *CRedis) DecrData(table string) int64 {
	c := r.Pool.Get()
	value, err := redis.Int64(c.Do("DECR", table))
	c.Close()
	if nil != err {
		log.Error("error table:%v", table)
		return 0
	}
	return value
}

//hash

func (r *CRedis) Hash_SetData(table string, name, value interface{}) error {
	c := r.Pool.Get()
	_, err := c.Do("hset", table, name, value)
	if nil != err {
		log.Error("error table:%v, name:%v", table, name)
	}
	c.Close()
	return err
}

func (r *CRedis) Hash_SetDataMap(table string, data []interface{}) error {
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(data))
	args[0] = table
	copy(args[1:], data)
	_, err := c.Do("hmset", args...)
	if nil != err {
		log.Error("error table:%v, data:%v", table, data)
	}
	c.Close()
	return err
}

func (r *CRedis) Hash_SetDataMap2(table string, data map[string]interface{}) error {
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(data)*2)
	args[0] = table
	i := 1
	for k, v := range data {
		args[i] = k
		args[i+1] = v
		i += 2
	}
	_, err := c.Do("hmset", args...)
	if nil != err {
		log.Error("error table:%v, data:%v", table, data)
	}
	c.Close()
	return err
}

func (r *CRedis) Hash_DelData(table string, name interface{}) error {
	c := r.Pool.Get()
	_, err := c.Do("hdel", table, name)
	if nil != err {
		log.Error("error table:%v, name:%v, err:%v", table, name, err)
	}
	c.Close()
	return err
}

func (r *CRedis) Hash_DelDataMap(table string, name []string) error {
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(name))
	args[0] = table
	i := 1
	for _, v := range name {
		args[i] = v
		i += 1
	}
	_, err := c.Do("hdel", args...)
	if nil != err {
		log.Error("error table:%v, name:%v, err:%v", table, name, err)
	}
	c.Close()
	return err
}

func (r *CRedis) Hash_GetDataString(table string, name interface{}) (string, error) {
	c := r.Pool.Get()
	value, err := redis.String(c.Do("HGET", table, name))
	if nil != err {
		log.Error("error table:%v, name:%v", table, name)
	}
	c.Close()
	return value, err
}

func (r *CRedis) Hash_GetDataInt64(table, name string) (int64, error) {
	c := r.Pool.Get()
	value, err := redis.Int64(c.Do("HGET", table, name))
	if nil != err {
		log.Error("error table:%v, name:%v", table, name)
	}
	c.Close()
	return value, err
}

func (r *CRedis) Hash_GetDataMap(table string, data []string) (map[string]string, error) {
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(data))
	args[0] = table
	i := 1
	for _, v := range data {
		args[i] = v
		i += 1
	}
	value, err := redis.Strings(c.Do("hmget", args...))
	if nil != err {
		log.Error("keys:%v, error:%v", args, err)
	}
	m := make(map[string]string)
	for k, v := range value {
		if k < len(data) {
			m[data[k]] = v
		}
	}
	c.Close()
	return m, err
}

func (r *CRedis) Hash_GetAllData(table string) (map[string]string, error) {
	c := r.Pool.Get()
	value, err := redis.Strings(c.Do("hgetall", table))
	if nil != err {
		log.Error("table:%v, error:%v", table, err)
	}
	m := make(map[string]string)
	for k, v := range value {
		if 0 != k%2 {
			m[value[k-1]] = v
		}
	}
	c.Close()
	return m, err
}

func (r *CRedis) Hash_GetData(table string, keys []string, data interface{}) bool {
	vt := reflect.ValueOf(data)
	if vt.IsNil() {
		log.Debug("err value:%v", data)
		return false
	}
	res, err := r.Hash_GetDataMap(table, keys)
	if nil != err {
		return false
	} else if len(res) > 0 {
		for k, v := range res {
			name := vt.Elem().FieldByName(k).Type().Name()
			var value interface{}
			switch name {
			case "int32":
				a, _ := strconv.Atoi(v)
				value = int32(a)
				break
			case "int64":
				value, _ = strconv.ParseInt(v, 10, 64)
				break
			case "string":
				value = v
				break
			case "float64":
				value, _ = strconv.ParseFloat(v, 64)
				break
			}

			vt.Elem().FieldByName(k).Set(reflect.ValueOf(value))
		}
	}
	return true
}

//sorted set
func (r *CRedis) ZSet_IsMember(table, member string) bool {
	c := r.Pool.Get()
	value, err := redis.Int64(c.Do("zrank", table, member))
	if nil != err {
		log.Error("table:%v, error:%v", table, err)
		c.Close()
		return false
	}
	log.Debug("value:%v", value)
	c.Close()
	return true
}

func (r *CRedis) ZSet_AddData(table string, value, score interface{}) error {
	c := r.Pool.Get()
	log.Debug("table:%v, value:%v, score:%v", table, value, score)
	_, err := c.Do("zadd", table, score, value)
	if nil != err {
		log.Error("table:%v, value:%v, score:%v, error:%v", table, value, score, err)
	}
	c.Close()
	return nil
}

func (r *CRedis) ZSet_AddDataMap(table string, data map[string]interface{}) error {
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(data)*2)
	args[0] = table
	i := 1
	for k, v := range data {
		args[i] = v
		args[i+1] = k
		i += 2
	}
	_, err := c.Do("zadd", args...)
	if nil != err {
		log.Error("table:%v, data:%v,  error:%v", table, data, err)
	}
	c.Close()
	return err

}

func (r *CRedis) ZSet_AddDataMap2(table string, data map[int64]int64) error {
	log.Debug("table:%v, data:%v", table, data)
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(data)*2)
	args[0] = table
	i := 1
	for k, v := range data {
		args[i] = v
		args[i+1] = k
		i += 2
	}
	_, err := c.Do("zadd", args...)
	if nil != err {
		log.Error("table:%v, data:%v,  error:%v", table, data, err)
	}
	c.Close()
	return err

}

func (r *CRedis) ZSet_DelData(table, value interface{}) error {
	c := r.Pool.Get()
	_, err := c.Do("zrem", table, value)
	if nil != err {
		log.Error("table:%v, value:%v error:%v", table, value, err)
	}
	c.Close()
	return err
}

func (r *CRedis) ZSet_DelDataByScore1(table, min, max interface{}) error {
	c := r.Pool.Get()
	_, err := c.Do("ZREMRANGEBYSCORE ", table, min, max)
	if nil != err {
		log.Error("table:%v, min:%v, max:%v, error:%v", table, min, max, err)
	}
	c.Close()
	return err
}

func (r *CRedis) ZSet_DelDataMap(table string, data []string) error {
	c := r.Pool.Get()
	args := make([]interface{}, 1+len(data))
	args[0] = table
	i := 1
	for _, v := range data {
		args[i] = v
		i += 1
	}
	_, err := c.Do("zrem", args...)
	if nil != err {
		log.Error("error table:%v, data:%v, err:%v", table, data, err)
	}
	c.Close()
	return err

}

func (r *CRedis) ZSet_GetDataByIndex(table string, start, end int32) ([]string, error) {
	c := r.Pool.Get()
	v, err := redis.Strings(c.Do("ZRANGE", table, start, end))
	if nil != err {
		log.Error("table:%v, start:%v, end:%v, error:%v", table, start, end, err)
	}
	log.Debug("v:%v", v)
	c.Close()
	return v, err
}

//table...    (newtable, number, table1,table2...)
func (r *CRedis) ZSet_InterStore(table ...interface{}) (int32, error) {
	c := r.Pool.Get()
	v, err := redis.Int(c.Do("ZINTERSTORE", table...))
	if nil != err {
		log.Error("v:%v, error:%v", v, err)
	}
	c.Close()
	return int32(v), err
}

func (r *CRedis) ZSet_GetDataWithScore(table string, start, end int32) (map[string]string, error) {
	c := r.Pool.Get()
	v, err := redis.StringMap(c.Do("ZREVRANGE", table, start, end, "WITHSCORES"))
	if nil != err {
		log.Error("v:%v, error:%v", v, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetDataWithScoreLimit(table string, start, end int32) ([]string, error) {
	c := r.Pool.Get()
	v, err := redis.Strings(c.Do("ZREVRANGE", table, start, end, "WITHSCORES"))
	if nil != err {
		log.Error("v:%v, error:%v", v, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetDataWithScoreLimit2(table string, start, end int32) ([]int64, error) {
	c := r.Pool.Get()
	v, err := redis.Int64s(c.Do("ZREVRANGE", table, start, end, "WITHSCORES"))
	if nil != err {
		log.Error("v:%v, error:%v", v, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetDataLimit(table string, start, end int32) ([]string, error) {
	c := r.Pool.Get()
	v, err := redis.Strings(c.Do("ZREVRANGE", table, start, end))
	if nil != err {
		log.Error("v:%v, error:%v", v, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetDataScore(table, name string) (int64, error) {
	c := r.Pool.Get()
	v, err := redis.Int64(c.Do("ZSCORE", table, name))
	if nil != err {
		log.Error("table:%v, name:%v, error:%v", table, name, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetAllDataByScore(table string) ([]string, error) {
	c := r.Pool.Get()
	v, err := redis.Strings(c.Do("ZRANGEBYSCORE", table, "-inf", "+inf"))
	if nil != err {
		log.Error("table:%v, error:%v", table, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetAllDataByScoreWithScore(table string) (map[string]int64, error) {
	c := r.Pool.Get()
	v, err := redis.Int64Map(c.Do("ZRANGEBYSCORE", table, "-inf", "+inf", "WITHSCORES"))
	if nil != err {
		log.Error("table:%v, error:%v", table, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetDataByScore(table, id interface{}) ([]string, error) {
	c := r.Pool.Get()
	v, err := redis.Strings(c.Do("ZRANGEBYSCORE", table, id, id))
	if nil != err {
		log.Error("v:%v, error:%v", v, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetDataByScoreLimitCount(table, id, count interface{}) ([]string, error) {
	c := r.Pool.Get()
	v, err := redis.Strings(c.Do("ZRANGEBYSCORE", table, id, id, "LIMIT", 0, count))
	if nil != err {
		log.Error("v:%v, error:%v", v, err)
	}
	c.Close()
	return v, err
}

func (r *CRedis) ZSet_GetDataByScoreOne(table, id, count interface{}) (string, error) {
	c := r.Pool.Get()
	v, err := redis.Strings(c.Do("ZRANGEBYSCORE", table, id, id, "LIMIT", 0, count))
	if nil != err {
		log.Error("table:%v,name:%v, error:%v", table, id, err)
	}
	c.Close()
	if len(v) > 0 {
		return v[0], err
	} else {
		return "", err
	}

}

func (r *CRedis) ZSet_DelDataByScore(table, id interface{}) error {
	c := r.Pool.Get()
	_, err := (c.Do("ZREMRANGEBYSCORE", table, id, id))
	if nil != err {
		log.Error("error:%v", err)
	}
	c.Close()
	return nil
}

func (r *CRedis) ZSet_GetDataCount(table string) int64 {
	c := r.Pool.Get()
	v, err := redis.Int64(c.Do("ZCARD", table))
	if nil != err {
		log.Error("error:%v", err)
	}
	c.Close()
	return v
}

func (r *CRedis) ZSet_GetDataRank(table string, uid int64) int64 {
	c := r.Pool.Get()
	v, err := redis.Int64(c.Do("ZREVRANK", table, uid))
	if nil != err {
		log.Error("error:%v", err)
		c.Close()
		return 0
	}
	c.Close()
	return v + 1
}
