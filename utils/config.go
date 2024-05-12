package utils

import (
	"fmt"
	"os"
	"time"
)

var (
	DSN string

	//自动刷新token距token签发需经过的时间
	TokenRefreshAfter      time.Duration = 300 * time.Second
	AdminTokenDuration     time.Duration = time.Hour * 24 * 2 //管理员不操作多久后token失效
	ApplicantTokenDuration time.Duration = time.Hour * 24 * 7 //候选人不操作多久后token失效

	IdentityKey []byte //加密凭据的私钥

	DoResetDb bool = false
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

	if readEnv("ResetDb", "false") == "true" {
		DoResetDb = true
	}

	//WARN: 生产环境请勿使用默认IDENTITY_KEY
	IdentityKey = Sha256(RawBytes(readEnv("IDENTITY_KEY", DSN)), 16)
}
