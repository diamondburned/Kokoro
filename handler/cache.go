package handler

import (
	"time"

	"github.com/Gigamons/common/helpers"

	"github.com/go-redis/redis"
)

type Cache struct {
	CacheHash   string
	DeleteCache time.Time
	Timeout     int64
}

var CLIENT *redis.Client
var CACHED []*Cache

func SetCache(CacheKey string, CacheData []byte, Timeout int64) *Cache {
	k, err := helpers.MD5String(CacheKey)
	if err != nil {
		return nil
	}
	C := &Cache{CacheHash: string(k), DeleteCache: time.Now(), Timeout: Timeout}
	if CLIENT.HSet(C.CacheHash, "Data", CacheData).Err() != nil {
		return nil
	}
	CACHED = append(CACHED, C)
	return C
}

func DelCache(CacheKey string) {
	k, err := helpers.MD5String(CacheKey)
	if err != nil {
		return
	}
	for i := 0; i < len(CACHED); i++ {
		if CACHED[i].CacheHash == string(k) {
			CLIENT.HDel(string(k), "Data")
			copy(CACHED[i:], CACHED[i+1:])
			CACHED[len(CACHED)-1] = nil
			CACHED = CACHED[:len(CACHED)-1]
		}
	}
}

func GetCache(CacheKey string) ([]byte, error) {
	k, err := helpers.MD5String(CacheKey)
	if err != nil {
		return nil, err
	}
	b, err := CLIENT.HExists(string(k), "Data").Result()
	if b && err == nil {
		return CLIENT.HGet(string(k), "Data").Bytes()
	}
	return nil, err
}

func init() {
	StartCacheCheck()
}

func StartCacheCheck() {
	go func() {
		for {
			for i := 0; i < len(CACHED); i++ {
				if CACHED[i].DeleteCache.Unix() < time.Now().Unix()-CACHED[i].Timeout {
					CLIENT.HDel(CACHED[i].CacheHash, "Data")
					copy(CACHED[i:], CACHED[i+1:])
					CACHED[len(CACHED)-1] = nil
					CACHED = CACHED[:len(CACHED)-1]
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
