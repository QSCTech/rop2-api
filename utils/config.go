package utils

import (
	"fmt"
	"os"
	"time"
)

var (
	DSN string

	TimeOffset    int64  = 1704038400  //所有相对时间戳的起点，相对Unix时间戳的起点，所偏移的秒数
	TokenDuration uint32 = 60 * 60 * 8 //token的有效秒数，无操作相应时间后失效

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
	IdentityKey = ToBytes(readEnv("IDENTITY_KEY", time.Now().String()))
}
