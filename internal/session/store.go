package session

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
	"github.com/knadh/koanf"
	"time"
)

type sessionConf struct {
	Expiration     time.Duration `koanf:"expiration"`
	CookieName     string        `koanf:"cookie_name"`
	CookieDomain   string        `koanf:"cookie_domain"`
	CookiePath     string        `koanf:"cookie_path"`
	CookieSecure   bool          `koanf:"cookie_secure"`
	CookieHTTPOnly bool          `koanf:"cookie_httponly"`
	CookieSameSite string        `koanf:"cookie_same_site"`
}

type redisConf struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	URL      string `koanf:"url"`
	Database int    `koanf:"database"`
	Reset    bool   `koanf:"reset"`
}

func New(cfg *koanf.Koanf) (*session.Store, error) {
	var (
		sessionConf sessionConf
		redisConf   redisConf
	)
	if err := cfg.Unmarshal("session", &sessionConf); err != nil {
		return nil, fmt.Errorf("error loading session config: %v", err)
	}
	if err := cfg.Unmarshal("redis", &redisConf); err != nil {
		return nil, fmt.Errorf("error loading redis config: %v", err)
	}

	redisStore := redis.New(redis.Config{
		Host:     redisConf.Host,
		Port:     redisConf.Port,
		Username: redisConf.Username,
		Password: redisConf.Password,
		Database: redisConf.Database,
		URL:      redisConf.URL,
		Reset:    redisConf.Reset,
	})

	sessionStore := session.New(session.Config{
		Storage:        redisStore,
		Expiration:     sessionConf.Expiration,
		KeyLookup:      sessionConf.CookieName,
		CookieDomain:   sessionConf.CookieDomain,
		CookiePath:     sessionConf.CookiePath,
		CookieSecure:   sessionConf.CookieSecure,
		CookieHTTPOnly: sessionConf.CookieHTTPOnly,
		CookieSameSite: sessionConf.CookieSameSite,
	})

	return sessionStore, nil

}
