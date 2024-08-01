package core

import (
    "fmt"
    "os"

    "github.com/bradfitz/gomemcache/memcache"
)

type Cache interface {
    Get(string) ([]byte, error)
    Set(string, []byte) error
}

type MemcachedCache struct {
    client *memcache.Client
}

func GetDefaultCacheServer() string {
    return os.Getenv("CACHE_SERVER")
}

func ConnectCache(server string) Cache {
    client := memcache.New(server)
    err := client.Ping()
    if err != nil {
        fmt.Println(err)
    }

    return MemcachedCache{client: client}
}

func (mc MemcachedCache) Get(key string) ([]byte, error) {
    item, err := mc.client.Get(key)
    if item == nil || err != nil {
        return nil, err
    }

    return item.Value, nil 
}

func (mc MemcachedCache) Set(key string, value []byte) error {
    mc.client.Set(&memcache.Item{Key: key, Value: value})
    return nil
}

