package utils

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr                   string `yaml:"Addr"`
	DSN                    string `yaml:"DSN"`
	MutipleChoicesRedirect string `yaml:"MutipleChoicesRedirect"`

	LoginCallbackRegex string `yaml:"LoginCallbackRegex"`

	// TokenDomain string `yaml:"TokenDomain"`
	// TokenPath   string `yaml:"TokenPath"`
	// TokenSecure bool   `yaml:"TokenSecure"`

	//加密凭据的密钥
	IdentityKey string `yaml:"IdentityKey"`
	//自动刷新token距token签发需经过的时间
	TokenRefreshAfter time.Duration `yaml:"TokenRefreshAfter"`
	//token有效期
	TokenDuration time.Duration `yaml:"TokenDuration"`

	CORSAllowOrigins []string `yaml:"CORSAllowOrigins"`
}

var Cfg = Config{
	Addr:                   "127.0.0.1:8080",
	DSN:                    "root:root@tcp(localhost:3306)/rop2?charset=utf8mb4&loc=Local&parseTime=true",
	MutipleChoicesRedirect: "http://localhost:5173/login/choice",
	LoginCallbackRegex:     "^http://localhost:5173(/.*)?$",
	TokenRefreshAfter:      time.Minute * 2,
	TokenDuration:          time.Hour * 24 * 1,
	CORSAllowOrigins:       []string{"http://localhost:5173"},
}

var (
	LoginCallbackRegex *regexp.Regexp
	IdentityKey        []byte
)

func argContains(str string) bool {
	for _, v := range os.Args {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

// 读取配置
func Init() {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	LoginCallbackRegex = regexp.MustCompile(Cfg.LoginCallbackRegex)
	if Cfg.IdentityKey == "" {
		log.Fatal("IdentityKey is empty")
		return
	}
	IdentityKey = Sha256([]byte(Cfg.IdentityKey), 16)
}
