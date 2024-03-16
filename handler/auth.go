package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"rop2-api/model"
	"rop2-api/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
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
	needKeep(now time.Time) bool          //检查此时是否还需保留此失效信息
	needVoid(identity *UserIdentity) bool //检查指定的identity是否因此失效，注意zjuid不需要检查
}

// 由于特定原因导致的zjuid-失效记录的map
var voidMap map[string][]voidInfo = make(map[string][]voidInfo)

// 主动退出某处登录，使特定签发时间的token失效
type voidOne struct {
	iat uint32
}

func (info voidOne) needKeep(now time.Time) bool {
	const secGap = 30 //保证此失效记录完全覆盖有效期的小间隙
	return utils.ToRelTimestamp(now) > info.iat+utils.TokenDuration+secGap
}

func (info voidOne) needVoid(status *UserIdentity) bool {
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

func (info voidBefore) needVoid(status *UserIdentity) bool {
	return status.Iat <= info.before
}

type RequireAuthOptions struct {
	AllowRefresh bool
}

// 中间件，要求用户必须登录才能访问API。
// 用户信息(UserIdentity类型)存至ctx.Keys["identity"]。
// 同时，如果有效token签发时间已经超过一个阙值，则在header提供一个新的token
func AuthWithRefresh(allowRefresh bool) gin.HandlerFunc {
	//subCode不小于0，不大于999
	code401 := func(subCode int32) *utils.ErrCodeObj {
		return utils.CodeObj(401, subCode)
	}
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("rop-token")
		parts := strings.Split(token, " ")
		if len(parts) != 2 {
			ctx.AbortWithStatusJSON(401, code401(1)) //格式无效
			return
		}

		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case base64.CorruptInputError:
					ctx.AbortWithStatusJSON(401, code401(2)) //base64编码无效
				default:
					panic(err)
				}
			}
		}() //捕获base64编码错误
		bArr := utils.MapArray(parts, func(part string, i int) []byte { return utils.Base64Decode(part) })
		jsonBytes, signBytes := bArr[0], bArr[1]

		now := utils.ToRelTimestamp(time.Now())
		iden := &UserIdentity{}
		if err := jsoniter.ConfigFastest.Unmarshal(jsonBytes, &iden); err != nil {
			ctx.AbortWithStatusJSON(401, code401(3)) //反序列化失败
			return
		}
		if iden.Exp <= now {
			ctx.AbortWithStatusJSON(401, code401(11)) //过期（还未验签）
			return
		}

		validSign := utils.HmacSha256(jsonBytes, utils.IdentityKey)
		if !bytes.Equal(validSign, signBytes) {
			ctx.AbortWithStatusJSON(401, code401(21)) //验签失败
			return
		}

		if voidArray, exists := voidMap[iden.ZjuId]; exists {
			for _, v := range voidArray {
				if v.needVoid(iden) {
					ctx.AbortWithStatusJSON(401, code401(31)) //已被登出
					return
				}
			}
		}

		//确认token有效
		if allowRefresh && now >= iden.Iat+utils.TokenRefreshAfter {
			//如果token已经过了一定时间，提供新的token
			ctx.Header("rop-refresh-token", newToken(&model.User{
				ZjuId:    iden.ZjuId,
				At:       iden.At,
				Nickname: iden.Nickname,
				Perm:     iden.Perm,
				//忽略CreateAt等信息
			}))
		}
		ctx.Set("identity", iden)
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
	idenBytes := utils.RawBytes(idenJson)
	//base64编码，不含padding(=)
	idenB64 := utils.Base64Encode(idenBytes)

	signBytes := utils.HmacSha256(idenBytes, utils.IdentityKey)
	signB64 := utils.Base64Encode(signBytes)
	return fmt.Sprintf("%s %s", idenB64, signB64)
}

func login(ctx *gin.Context) {
	//TODO 测试用，直接登录
	ctx.PureJSON(200, gin.H{
		"token": newToken(model.TestUser),
	})
}

func logout(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*UserIdentity)
	addVoidInfo(id.ZjuId, voidOne{iat: id.Iat})
	ctx.PureJSON(200, utils.CodeObj())
}

func logoutAll(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*UserIdentity)
	addVoidInfo(id.ZjuId, voidBefore{before: utils.ToRelTimestamp(time.Now())})
	ctx.PureJSON(200, utils.CodeObj())
}

func addVoidInfo(zjuId string, info voidInfo) {
	v, exists := voidMap[zjuId]
	if exists {
		voidMap[zjuId] = append(v, info)
	} else {
		voidMap[zjuId] = []voidInfo{info}
	}
}

func authInit(routerGroup *gin.RouterGroup) {
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

	routerGroup.POST("/login", login)
	routerGroup.GET("/logout", AuthWithRefresh(false), logout)
	routerGroup.GET("/logoutAll", AuthWithRefresh(false), logoutAll)
}
