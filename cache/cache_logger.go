package cache

import (
	"fmt"
	"github.com/mgutz/logxi/v1"
	"github.com/thejackrabbit/aero/key"
	"github.com/thejackrabbit/aero/panik"
	"os"
	"strings"
	"time"
)

type CacheLogger struct {
	key.AsIsFormat
	Inner  Cacher
	Logger log.Logger
}

func NewCacheLogger(dir string, inner Cacher) Cacher {

	if !strings.HasSuffix(dir, "/") && !strings.HasSuffix(dir, "\\") {
		dir += "/"
	}

	fmt.Println(dir)
	f, err := os.OpenFile(dir+"cache.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	panik.On(err)

	l := log.NewLogger3(f, "cache", &log.JSONFormatter{})
	l.SetLevel(log.LevelInfo)

	return CacheLogger{
		Inner:  inner,
		Logger: l,
	}
}

func (c CacheLogger) Set(key string, data []byte, expireIn time.Duration) {
	k := c.Format(key)
	c.Inner.Set(key, data, expireIn)
	c.Logger.Info("cache.set", "key-orig", k, "key-final", c.Inner.Format(key), "data", data)
}

func (c CacheLogger) Get(key string) ([]byte, error) {
	data, err := c.Inner.Get(key)
	k := c.Format(key)

	if err != nil {
		c.Logger.Info("cache.get", "key-orig", k, "key-final", c.Inner.Format(key), "data", "<not-found>")
		return nil, err
	} else {
		c.Logger.Info("cache.get", "key-orig", k, "key-final", c.Inner.Format(key), "data", data)
		return data, nil
	}
}

func (c CacheLogger) Delete(key string) error {
	k := c.Format(key)
	c.Logger.Info("cache.delete", "key-orig", k, "key-final", c.Inner.Format(key))
	return c.Inner.Delete(key)
}

func (c CacheLogger) Close() {
	c.Inner.Close()
	c.Logger.Info("cache.close")
}