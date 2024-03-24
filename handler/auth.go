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
type AdminIdentity struct {
	Iat      time.Time       `json:"iat"` //签发时间
	Exp      time.Time       `json:"exp"` //过期时间
	ZjuId    string          `json:"zjuId"`
	At       uint32          `json:"at"`       //登录组织id
	Nickname string          `json:"nickname"` //组织内昵称
	Level    model.PermLevel `json:"level"`    //部门权限json
}

type voidInfo interface {
	needKeep(now time.Time) bool           //检查此时是否还需保留此失效信息
	needVoid(identity *AdminIdentity) bool //检查指定的identity是否因此失效，注意zjuid不需要检查
}

// 由于特定原因导致的zjuid-失效记录的map
var voidMap map[string][]voidInfo = make(map[string][]voidInfo)

// 主动退出某处登录，使特定签发时间的token失效
type voidOne struct {
	iat time.Time
}

func (info voidOne) needKeep(now time.Time) bool {
	const secGap = 5 * time.Second //保证此失效记录完全覆盖有效期的小间隙
	return info.iat.Add(utils.TokenDuration).Add(secGap).Compare(now) >= 0
}

func (info voidOne) needVoid(status *AdminIdentity) bool {
	return status.Iat == info.iat
}

// 退出所有登录，使签发时间小于某个点的token全部失效
type voidBefore struct {
	before time.Time
}

func (info voidBefore) needKeep(now time.Time) bool {
	const secGap = 5 * time.Second //保证此失效记录完全覆盖有效期的小间隙
	return info.before.Add(utils.TokenDuration).Add(secGap).Compare(now) >= 0
}

func (info voidBefore) needVoid(status *AdminIdentity) bool {
	return info.before.Compare(status.Iat) >= 0
}

// 中间件，要求用户必须登录才能访问API。
// 管理员信息(AdminIdentity类型)存至ctx.Keys["identity"]。
// 同时，如果有效token签发时间已经超过一个阙值，则在header提供一个新的token
func AuthWithRefresh(allowRefresh bool) gin.HandlerFunc {
	//subCode不小于0，不大于999
	code401 := func(message string, subCode int) (int, *utils.CodeMessageObj) {
		return utils.Message(message, 401, subCode)
	}
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("rop-token")
		parts := strings.Split(token, " ")
		if len(parts) != 2 {
			ctx.AbortWithStatusJSON(code401("token无法识别", 1))
			return
		}

		defer func() { //捕获base64编码错误
			if err := recover(); err != nil {
				switch err.(type) {
				case base64.CorruptInputError:
					ctx.AbortWithStatusJSON(code401("token编码无效", 2))
				default:
					panic(err)
				}
			}
		}()
		bArr := utils.MapArray(parts, func(part string, i int) []byte { return utils.Base64Decode(part) })
		jsonBytes, signBytes := bArr[0], bArr[1]

		now := time.Now()
		iden := &AdminIdentity{}
		if err := jsoniter.ConfigFastest.Unmarshal(jsonBytes, &iden); err != nil {
			ctx.AbortWithStatusJSON(code401("token反序列化失败", 3))
			return
		}
		if now.Compare(iden.Exp) >= 0 {
			ctx.AbortWithStatusJSON(code401("token已过期", 11))
			return
		}

		validSign := utils.HmacSha256(jsonBytes, utils.IdentityKey)
		if !bytes.Equal(validSign, signBytes) {
			ctx.AbortWithStatusJSON(code401("token验签失败", 21))
			return
		}

		if voidArray, exists := voidMap[iden.ZjuId]; exists {
			for _, v := range voidArray {
				if v.needVoid(iden) {
					ctx.AbortWithStatusJSON(code401("已退出登录", 31))
					return
				}
			}
		}

		//已确认token有效
		if allowRefresh && now.Compare(iden.Iat.Add(utils.TokenRefreshAfter)) >= 0 {
			//如果token已经过了一定时间，提供新的token
			ctx.Header("rop-refresh-token", newToken(&model.Admin{
				ZjuId:    iden.ZjuId,
				At:       iden.At,
				Nickname: iden.Nickname,
				Level:    iden.Level,
				//忽略CreateAt等信息
			}))
		}
		ctx.Set("identity", iden)
		ctx.Next()
	}
}

// 内部方法，生成一个新token
func newToken(user *model.Admin) string {
	iat := time.Now()
	iden := &AdminIdentity{
		Iat:      iat,
		Exp:      iat.Add(utils.TokenDuration),
		ZjuId:    user.ZjuId,
		At:       user.At,
		Nickname: user.Nickname,
		Level:    user.Level,
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

func adminLogin(ctx *gin.Context) {
	//TODO 测试用，直接登录
	type Arg struct {
		ZjuId string `json:"zju_id"`
		At    uint32 `json:"at"`
	}
	arg := &Arg{}
	if ctx.ShouldBindJSON(arg) != nil {
		ctx.AbortWithStatusJSON(utils.Message("参数绑定失败", 400, 0))
		return
	}

	admin := model.GetAdmin(arg.ZjuId, arg.At)
	if admin == nil {
		ctx.AbortWithStatusJSON(utils.Message("用户不存在", 400, 1))
		return
	}
	ctx.Header("rop-refresh-token", newToken(admin))
	ctx.PureJSON(utils.Success())
}

func logout(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*AdminIdentity)
	addVoidInfo(id.ZjuId, voidOne{iat: id.Iat})
	ctx.PureJSON(utils.Success())
}

func logoutAll(ctx *gin.Context) {
	id := ctx.MustGet("identity").(*AdminIdentity)
	addVoidInfo(id.ZjuId, voidBefore{before: time.Now()})
	ctx.PureJSON(utils.Success())
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

	routerGroup.POST("/login", adminLogin)
	routerGroup.GET("/logout", AuthWithRefresh(false), logout)
	routerGroup.GET("/logoutAll", AuthWithRefresh(false), logoutAll)
}
