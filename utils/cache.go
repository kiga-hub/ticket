package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"

	_ "github.com/astaxie/beego/cache/redis"
)

var cc cache.Cache

func InitCache() {
	host := beego.AppConfig.String("cache::redis_host")
	passWord := beego.AppConfig.String("cache::redis_password")
	var err error
	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()
	mstr := map[string]string{}
	mstr["key"] = "ticket"
	mstr["conn"] = host
	mstr["dbNum"] = "0"
	if passWord != "" {
		mstr["password"] = passWord
	}
	bytes, _ := json.Marshal(mstr)
	cc, err = cache.NewCache("redis", string(bytes))
	if err != nil {
		LogError("Connect to the redis host " + host + " failed")
		LogError(err)
	}
}

// SetCache
func SetCache(key string, value interface{}, timeout int) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			LogError(r)
			cc = nil
		}
	}()
	timeouts := time.Duration(timeout) * time.Second
	err = cc.Put(key, data, timeouts)
	if err != nil {
		LogError(err)
		LogError("Set cache failed, key:" + key)
		return err
	} else {
		return nil
	}
}

func GetCache(key string, to interface{}) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			LogError(r)
			cc = nil
		}
	}()

	data := cc.Get(key)
	if data == nil {
		return errors.New("Cache does not exist")
	}

	err := Decode(data.([]byte), to)
	if err != nil {
		LogError(err)
		LogError("Get Cache failed, key:" + key)
	}
	return nil
}

// DelCache
func DelCache(key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			//fmt.Println("get cache error caught: %v\n", r)
			cc = nil
		}
	}()
	err := cc.Delete(key)
	if err != nil {
		return errors.New("Cache delete failed")
	} else {
		return nil
	}
}

// Encode
// To use gob for data encoding
func Encode(data interface{}) (interface{}, error) {
	switch data.(type) {
	case int:
		return data, nil
	}
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode
// To use gob for data decoding
func Decode(data []byte, to interface{}) error {
	if count, err := strconv.Atoi(string(data)); err == nil {
		*to.(*int) = count
		return nil
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
