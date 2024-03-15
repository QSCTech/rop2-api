package handler

import (
	"encoding/base64"
	"fmt"
	"rop2-api/model"
	"rop2-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 凭据：类JWT的base64字符串
// 包含zjuid at nickname perm等信息
// 使用签名(不在结构内)保证权威性
type UserIdentity struct {
	Iat      uint32 `json:"iat"` //签发时间，相对时间戳
	Exp      uint32 `json:"exp"` //过期时间，相对时间戳
	ZjuId    string `json:"zjuId"`
	At       uint32 `json:"at"`       //登录组织id
	Nickname string `json:"nickname"` //组织内昵称
	Perm     string `json:"perm"`     //部门权限json
}

type voidInfo interface {
	needKeep(now time.Time) bool         //检查此时是否还需保留此失效信息
	needVoid(identity UserIdentity) bool //检查指定的identity是否因此失效，注意zjuid不需要检查
}

// 由于特定原因导致的zjuid-失效记录的map
var voidMap map[string][]voidInfo

// 主动退出某处登录，使特定签发时间的token失效
type voidOne struct {
	iat uint32
}

func (info voidOne) needKeep(now time.Time) bool {
	const secGap = 30 //保证此失效记录完全覆盖有效期的小间隙
	return utils.ToRelTimestamp(now) > info.iat+utils.TokenDuration+secGap
}

func (info voidOne) needVoid(status UserIdentity) bool {
	return status.Iat == info.iat
}

// 退出所有登录，使签发时间小于某个点的token全部失效
type voidBefore struct {
	before uint32
}

func (info voidBefore) needKeep(now time.Time) bool {
	const secGap = 30 //保证此失效记录完全覆盖有效期的小间隙
	return utils.ToRelTimestamp(now) > info.before+utils.TokenDuration+secGap
}

func (info voidBefore) needVoid(status UserIdentity) bool {
	return status.Iat <= info.before
}

// 中间件，要求用户必须登录才能访问API。
// 用户信息(UserIdentity类型)存至ctx.Keys["identity"]
func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.GetHeader("rop-identity")
		//TODO 检查token
		ctx.Next()
	}
}

// 内部方法，生成一个新token
func newToken(user *model.User) string {
	iat := utils.ToRelTimestamp(time.Now())
	iden := &UserIdentity{
		Iat:      iat,
		Exp:      iat + utils.TokenDuration,
		ZjuId:    user.ZjuId,
		At:       user.At,
		Nickname: user.Nickname,
		Perm:     user.Perm,
	}
	idenJson := utils.Stringify(iden)
	//直接获取string底层的byte[]
	idenBytes := utils.ToBytes(idenJson)
	//base64编码，不含padding(=)
	idenB64 := base64.RawStdEncoding.EncodeToString(idenBytes)

	signBytes := utils.HmacSha256(idenBytes, utils.IdentityKey)
	signB64 := base64.RawStdEncoding.EncodeToString(signBytes)
	//TODO
	return fmt.Sprintf("%s %s", idenB64, signB64)
}

func LogoutOnce(identity *UserIdentity) {
	addVoidInfo(identity.ZjuId, voidOne{iat: identity.Iat})
}

func LogoutAll(identity *UserIdentity) {
	addVoidInfo(identity.ZjuId, voidBefore{before: utils.ToRelTimestamp(time.Now())})
}

func addVoidInfo(zjuId string, info voidInfo) {
	v, exists := voidMap[zjuId]
	if exists {
		voidMap[zjuId] = append(v, info)
	} else {
		voidMap[zjuId] = []voidInfo{info}
	}
}

func authInit() {
	voidMapCleanupTicker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			time := <-voidMapCleanupTicker.C
			for key, array := range voidMap {
				n := 0
				for _, info := range array {
					if info.needKeep(time) {
						array[n] = info
						n++
					}
				}
				if n == 0 {
					delete(voidMap, key)
				} else {
					voidMap[key] = array[:n]
				}
			}
		}
	}()
}
