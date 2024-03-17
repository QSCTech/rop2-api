package utils

import (
	"fmt"
	"os"
	"time"
)

var (
	DSN string

	TokenDuration     time.Duration = 60 * 60 * 8 * time.Second
	TokenRefreshAfter time.Duration = min(TokenDuration/25, 60*5) * time.Second //自动刷新token需经过的时间

	IdentityKey []byte
)

func readEnv(envKey, defaultValue string) string {
	if v, ok := os.LookupEnv(fmt.Sprintf("ROP2_%s", envKey)); ok {
		return v
	}
	return defaultValue
}

// 读取配置
func Init() {
	DSN = readEnv("DSN", "root:root@tcp(localhost:3306)/rop2?charset=utf8mb4&parseTime=true")
	//默认值可以考虑改成机器唯一id
	IdentityKey = RawBytes(readEnv("IDENTITY_KEY", time.Now().String()))
}
