package handler

import (
	"crypto/md5"
	"encoding/hex"
	"time"

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
	C := &Cache{CacheHash: hex.EncodeToString(md5.New().Sum([]byte(CacheKey))), DeleteCache: time.Now(), Timeout: Timeout}
	if CLIENT.HSet(C.CacheHash, "Data", CacheData).Err() != nil {
		return nil
	}
	CACHED = append(CACHED, C)
	return C
}

func DelCache(CacheKey string) {
	for i := 0; i < len(CACHED); i++ {
		if CACHED[i].CacheHash == hex.EncodeToString(md5.New().Sum([]byte(CacheKey))) {
			CLIENT.HDel(hex.EncodeToString(md5.New().Sum([]byte(CacheKey))), "Data")
			copy(CACHED[i:], CACHED[i+1:])
			CACHED[len(CACHED)-1] = nil
			CACHED = CACHED[:len(CACHED)-1]
		}
	}
}

func GetCache(CacheKey string) ([]byte, error) {
	c := hex.EncodeToString(md5.New().Sum([]byte(CacheKey)))
	b, err := CLIENT.HExists(c, "Data").Result()
	if b && err == nil {
		return CLIENT.HGet(c, "Data").Bytes()
	}
	return []byte{}, err
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
