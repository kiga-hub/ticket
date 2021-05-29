package lock

import (
    "fmt"
	"time"
	"errors"
	"Two-Card/utils"

	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego"
)

var pool *redis.Pool

func init() {
	redis_host := beego.AppConfig.String("cache::redis_host")
	redis_password := beego.AppConfig.String("cache::redis_password")
    pool_size := 20
    pool = redis.NewPool(func() (redis.Conn, error) {
        c, err := redis.Dial("tcp", redis_host, redis.DialPassword(redis_password))
        if err != nil {
            return nil, err
        }
        return c, nil
    }, pool_size)
}
 
func RedisClient() redis.Conn {
    return pool.Get()
}

func Lock(lock_name string, acquire_time, time_out int) (interface{}, error){
	redis_conn := RedisClient()
    defer redis_conn.Close()
	lock := fmt.Sprintf("Two-Card:lock:%s", lock_name)
	end := time.Now().Unix() + int64(acquire_time)
	identifier := utils.Uuid()
	for cur := time.Now().Unix(); cur < end; cur = time.Now().Unix() {
		cnt, _ := redis.Int(redis_conn.Do("SETNX", lock, identifier))
        if cnt == 0 {
			// ttl, _ := redis.Int(redis_conn.Do("TTL", lock))
            time.Sleep(time.Duration(1) * time.Second)  
        } else {
			redis_conn.Do("EXPIRE", lock, time_out)
			return identifier, nil
        }
	}
	return nil, nil
}

func UnLock(lock_name string, identifier interface{}) error {
	redis_conn := RedisClient()
    defer redis_conn.Close()
	lock := fmt.Sprintf("Two-Card:lock:%s", lock_name)
	id, err := redis.String(redis_conn.Do("GET", lock))
    if err != nil {
        return err
	}
	if (id != identifier.(string)) {
		return errors.New("not self lock")
	}
	_, err = redis_conn.Do("Del", lock)
	return err
}

// SetCache
func SetnxCache(key string, value interface{}, timeout int) error {
	redis_conn := RedisClient()
    defer redis_conn.Close()
	redis_key := fmt.Sprintf("Two-Card:%s", key)
	cnt, err := redis.Int(redis_conn.Do("SETNX", redis_key, value))
	if err != nil {
		return err
	}
	if cnt != 0 {
		redis_conn.Do("EXPIRE", redis_key, timeout)
	}
	return nil
}

func DecrCache(key string, value int) error {
	redis_conn := RedisClient()
	defer redis_conn.Close()
	redis_key := fmt.Sprintf("Two-Card:%s", key)
	cnt, err := redis.Int(redis_conn.Do("EXISTS", redis_key))
	if err != nil {
		return err
	}
	if cnt == 0 {
		return nil
	}

	cnt, err = redis.Int(redis_conn.Do("DECRBY", redis_key, value))
	if err != nil {
		return err
	}
	if cnt <= 0 {
		redis_conn.Do("DEL", redis_key)
	}
	return nil
}

func IncrCache(key string, value int) error {
	redis_conn := RedisClient()
    defer redis_conn.Close()
	redis_key := fmt.Sprintf("Two-Card:%s", key)
	cnt, err := redis.Int(redis_conn.Do("EXISTS", redis_key))
	if err != nil {
		return err
	}
	if cnt == 0 {
		return nil
	}
	cnt, err = redis.Int(redis_conn.Do("INCRBY", redis_key, value))
	if err != nil {
		return err
	}
	return nil
}

func ClearCache(key string) error {
	redis_conn := RedisClient()
    defer redis_conn.Close()
	redis_key := fmt.Sprintf("Two-Card:*%s*", key)
	cachedKeys, err := redis.Strings(redis_conn.Do("KEYS", redis_key))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = redis_conn.Do("DEL", str); err != nil {
			return err
		}
	}
	return nil
}

func PushQueue(queue_name string) (string, error) {
	// 获取redis连接
	redis_conn := RedisClient()
	defer redis_conn.Close()

	// 加锁
	lockKey := fmt.Sprintf("Two-Card:lock:%s", queue_name)
	lockEnd := time.Now().Unix() + int64(10)
	lockId := utils.Uuid()
	for cur := time.Now().Unix(); cur < lockEnd; cur = time.Now().Unix() {
		cnt, _ := redis.Int(redis_conn.Do("SETNX", lockKey, lockId))
        if cnt == 0 {
            time.Sleep(time.Duration(1) * time.Second)  
        } else {
			redis_conn.Do("EXPIRE", lockId, 10)
			break
        }
	}

	// 入队
	queueKey := fmt.Sprintf("%s:%s", queue_name, lockId)
	_, err := redis_conn.Do("SETEX", queueKey, 60, 1)
	if err != nil {
		return "", err
	}

	// 取队尾
	queueTailKey := fmt.Sprintf("%s:tail", queue_name)
	queueTail, ret := redis.String(redis_conn.Do("GET", queueTailKey))

	// 修改队尾
	_, err = redis_conn.Do("SETEX", queueTailKey, 600, queueKey)
	if err != nil {
		return "", err
	}

	// 解锁
	id, err := redis.String(redis_conn.Do("GET", lockKey))
    if err == nil && id == lockId {
		redis_conn.Do("Del", lockKey)
	}

	if ret == nil {
		// 等待
		for {
			_, err := redis.String(redis_conn.Do("GET", queueTail))
			if err != nil {
				break
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}

	return queueKey, nil
}

func PopQueue(queueKey string) error {
	// 获取redis连接
	redis_conn := RedisClient()
	defer redis_conn.Close()

	_, err := redis_conn.Do("Del", queueKey)
	if err != nil {
		return err
	}
	return nil
}