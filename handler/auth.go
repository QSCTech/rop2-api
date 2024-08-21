package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"rop2-api/model"
	"rop2-api/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 所有凭据(报名者、管理员)的抽象接口
type userIdentity interface {
	//检查是否可以刷新token，此方法不检查当前token有效期等
	//
	//返回值即为新的token，返回""表示不需要刷新
	canRefresh(now time.Time) string
	isValid(now time.Time) bool
	getId() string
	getIat() time.Time
}

// 凭据：类JWT的base64字符串
// 包含zjuid at nickname等信息
// 使用签名(不在结构内)保证权威性
type AdminIdentity struct {
	Iat      time.Time       `json:"iat"` //签发时间
	Exp      time.Time       `json:"exp"` //过期时间
	ZjuId    string          `json:"zjuId"`
	At       uint32          `json:"at"`       //登录组织id
	Nickname string          `json:"nickname"` //组织内昵称
	Level    model.PermLevel `json:"level"`    //权限级别
}

func (this AdminIdentity) canRefresh(now time.Time) string {
	if now.After(this.getIat().Add(utils.TokenRefreshAfter)) {
		copy := this //浅克隆this
		copy.Iat = now
		copy.Exp = copy.Iat.Add(utils.AdminTokenDuration)
		return newToken(copy)
	}
	return ""
}
func (this AdminIdentity) isValid(now time.Time) bool { return now.Before(this.Exp) }
func (this AdminIdentity) getId() string              { return this.ZjuId }
func (this AdminIdentity) getIat() time.Time          { return this.Iat }

type ApplicantIdentity struct {
	Iat   time.Time `json:"iat"` //签发时间
	Exp   time.Time `json:"exp"` //过期时间
	ZjuId string    `json:"zjuId"`
}

func (this ApplicantIdentity) canRefresh(now time.Time) string {
	if now.After(this.getIat().Add(utils.TokenRefreshAfter)) {
		copy := this //浅克隆this
		copy.Iat = now
		copy.Exp = copy.Iat.Add(utils.ApplicantTokenDuration)
		return newToken(copy)
	}
	return ""
}
func (this ApplicantIdentity) isValid(now time.Time) bool { return now.Before(this.Exp) }
func (this ApplicantIdentity) getId() string              { return this.ZjuId }
func (this ApplicantIdentity) getIat() time.Time          { return this.Iat }

type voidInfo interface {
	needKeep(now time.Time) bool         //检查此时是否还需保留此失效信息
	needVoid(identity userIdentity) bool //检查指定的identity是否因此失效，注意zjuid不需要检查
}

// 由于特定原因导致的zjuid-失效记录的map
var voidMap map[string][]voidInfo = make(map[string][]voidInfo)

// 主动退出某处登录，使特定签发时间的token失效
type voidOne struct {
	iat time.Time
}

func (info voidOne) needKeep(now time.Time) bool {
	const secGap = 5 * time.Second //保证此失效记录完全覆盖有效期的小间隙
	return info.iat.Add(utils.AdminTokenDuration).Add(secGap).After(now)
}
func (info voidOne) needVoid(status userIdentity) bool {
	return status.getIat().Sub(info.iat).Abs() <= 2*time.Second
}

// 退出所有登录，使签发时间小于某个点的token全部失效
type voidBefore struct {
	before time.Time
}

func (info voidBefore) needKeep(now time.Time) bool {
	const secGap = 5 * time.Second //保证此失效记录完全覆盖有效期的小间隙
	return info.before.Add(utils.AdminTokenDuration).Add(secGap).After(now)
}
func (info voidBefore) needVoid(status userIdentity) bool {
	return status.getIat().Compare(info.before) <= 0
}

// 从header读取token并转换，存储在resultPointer中，返回是否成功。
//
// golang默认json反序列化缺失字段不报错，必须另行是否是有效的AdminIdentity。
func parseToken[T userIdentity](ctx *gin.Context, resultPointer *T) bool {
	code401 := func(message string, subCode int) (int, *utils.CodeMessageObj) {
		return utils.Message(message, 401, subCode)
	}
	//token格式: base64encodedidentityjson base64sign
	token := ctx.GetHeader("rop-token")
	parts := strings.Split(token, " ")
	if len(parts) != 2 {
		ctx.AbortWithStatusJSON(code401("token无法识别", 1))
		return false
	}
	var err error
	bArr := utils.MapArray(parts, func(part string, i int) []byte {
		result, newErr := utils.Base64Decode(part)
		if newErr != nil {
			err = newErr
		}
		return result
	})
	if err != nil {
		ctx.AbortWithStatusJSON(code401("token编码无效", 2))
		return false
	}
	jsonBytes, signBytes := bArr[0], bArr[1]
	now := time.Now()
	if err := json.Unmarshal(jsonBytes, resultPointer); err != nil {
		ctx.AbortWithStatusJSON(code401("token反序列化失败", 3))
		return false
	}
	if !(*resultPointer).isValid(now) {
		ctx.AbortWithStatusJSON(code401("token已过期", 11))
		return false
	}
	if (*resultPointer).getIat().Before(utils.TokenValidSince) {
		ctx.AbortWithStatusJSON(code401("token签发过早", 12))
		return false
	}
	if validSign := utils.HmacSha256(jsonBytes, utils.IdentityKey); !bytes.Equal(validSign, signBytes) {
		ctx.AbortWithStatusJSON(code401("token验签失败", 21))
		return false
	}
	if voidArray, exists := voidMap[(*resultPointer).getId()]; exists {
		for _, v := range voidArray {
			if v.needVoid(*resultPointer) {
				ctx.AbortWithStatusJSON(code401("已退出登录", 31))
				return false
			}
		}
	}
	return true
}

// 中间件，要求用户必须进行管理员登录才能访问API。
// 管理员信息(AdminIdentity类型)存至ctx.Keys["identity"]。
// 同时，如果有效token签发时间已经超过一个阙值，则在header提供一个新的token
func RequireAdminWithRefresh(allowRefresh bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		iden := AdminIdentity{}
		if !parseToken(ctx, &iden) {
			return
		}

		if iden.Level <= model.Null || iden.At <= 0 {
			//未以管理员登录(候选人身份登录)
			ctx.AbortWithStatusJSON(utils.Message("暂无权限", 403, 11))
			return
		}

		//已确认token有效
		if allowRefresh {
			newToken := iden.canRefresh(time.Now())
			//不需要刷新token，对header设空字符串不会报错
			ctx.Header("rop-refresh-token", newToken)
		}
		ctx.Set("identity", iden)
		ctx.Next()
	}
}

// 中间件，要求用户必须进行登录才能访问API。适用于没有管理员权限的候选人登录。
// 候选人信息(ApplicantIdentity类型)存至ctx.Keys["identity"]。
// 同时，如果有效token签发时间已经超过一个阙值，则在header提供一个新的token
func RequireLoginWithRefresh(allowRefresh bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		iden := ApplicantIdentity{}
		if !parseToken(ctx, &iden) {
			return
		}

		//已确认token有效
		if allowRefresh {
			newToken := iden.canRefresh(time.Now())
			//不需要刷新token，对header设空字符串不会报错
			ctx.Header("rop-refresh-token", newToken)
		}
		ctx.Set("identity", iden)
		ctx.Next()
	}
}

// 中间件，要求用户在登录组织至少有指定的权限，否则403
func RequireLevel(requireLevel model.PermLevel) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.MustGet("identity").(AdminIdentity)
		if id.Level < requireLevel {
			ctx.AbortWithStatusJSON(utils.MessageForbidden())
			return
		}
		ctx.Next()
	}
}

// 内部方法，生成一个新token。不做任何检查。
// JSON序列化、base64编码、计算签名、拼接。
func newToken[T any](jsonInfo T) string {
	idenJson := utils.Stringify(jsonInfo)
	//直接获取string底层的byte[]
	idenBytes := utils.RawBytes(idenJson)
	//base64编码，不含padding(=)
	idenB64 := utils.Base64Encode(idenBytes)

	signBytes := utils.HmacSha256(idenBytes, utils.IdentityKey)
	signB64 := utils.Base64Encode(signBytes)
	return fmt.Sprintf("%s %s", idenB64, signB64)
}

func logout(ctx *gin.Context) {
	id := ctx.MustGet("identity").(userIdentity)
	addVoidInfo(id.getId(), voidOne{iat: id.getIat()})
	ctx.PureJSON(utils.Success())
}

func logoutAll(ctx *gin.Context) {
	id := ctx.MustGet("identity").(userIdentity)
	addVoidInfo(id.getId(), voidBefore{before: time.Now()})
	ctx.PureJSON(utils.Success())
}

// 内部方法（不是API），强制登出某个zjuId的所有登录
func ForceLogoutAll(zjuId string) {
	addVoidInfo(zjuId, voidBefore{before: time.Now()})
}

// 对指定的zjuId添加登录失效信息
func addVoidInfo(zjuId string, info voidInfo) {
	v, exists := voidMap[zjuId]
	if exists {
		voidMap[zjuId] = append(v, info)
	} else {
		voidMap[zjuId] = []voidInfo{info}
	}
}

func loginByToken(ctx *gin.Context) {
	type Arg struct {
		Token    string `form:"SESSION_TOKEN" binding:"required"`
		Continue string `form:"continue" binding:"required"`
		OrgId    uint32 `form:"orgId"`
	}
	arg := &Arg{}
	if ctx.ShouldBindQuery(arg) != nil {
		ctx.AbortWithStatusJSON(utils.MessageBindFail())
		return
	}
	token := arg.Token
	continueUrl := arg.Continue
	requestedOrgId := arg.OrgId
	//允许字符：字母、数字、下划线、短横线(减号)
	if matched, _ := regexp.MatchString((`^[\w-]+$`), token); !matched {
		ctx.AbortWithStatusJSON(utils.Message("token格式错误", 400, 1))
		return
	}
	if !utils.LoginCallbackRegex.MatchString(continueUrl) {
		ctx.AbortWithStatusJSON(utils.Message("回调地址无效", 400, 2))
		return
	}

	client := &http.Client{}
	//这个地址就硬编码在这里了。passport更新了的话 结果解析也要改
	req, _ := http.NewRequest("GET", "https://www.qsc.zju.edu.cn/passport/v4/profile", nil)
	req.AddCookie(&http.Cookie{
		Name:  "SESSION_TOKEN",
		Value: token,
	})
	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
		ctx.AbortWithStatusJSON(utils.MessageInternalError())
		return
	}
	defer resp.Body.Close()

	var result struct {
		Err  string `json:"err"`
		Code int    `json:"code"`
		Data *struct {
			Logined bool `json:"logined"`
			//注：以下json键名大写(AS IS)
			User *struct {
				Name  string `json:"Name"`
				ZjuId string `json:"ZjuId"`
			} `json:"User"`
		} `json:"data"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Err != "" || result.Code != 0 || !result.Data.Logined {
		ctx.AbortWithStatusJSON(utils.Message("token无效", 401, 51))
		return
	}
	zjuId := result.Data.User.ZjuId
	name := result.Data.User.Name
	if len(zjuId) <= 0 || len(name) <= 0 {
		ctx.AbortWithStatusJSON(utils.Message("无法获取个人信息", 401, 52))
		return
	}
	model.EnsurePerson(zjuId, name)

	// //TODO: DEBUG ONLY
	// //添加测试组织管理员权限
	// if model.TestOrgId > 0 {
	// 	model.SetAdmin(model.TestOrgId, zjuId, "", model.Maintainer)
	// }

	admin := model.GetAdmin(zjuId, requestedOrgId)
	now := time.Now()
	var ropToken string
	if len(admin) <= 0 {
		//无管理权限，候选人登录
		ropToken = newToken(ApplicantIdentity{
			Iat:   now,
			Exp:   now.Add(utils.ApplicantTokenDuration),
			ZjuId: zjuId,
		})
	} else {
		if len(admin) > 1 {
			orgProfiles := model.GetAvailableOrgs(zjuId)
			ctx.Redirect(302, utils.AddQuery(utils.MutipleChoicesRedirect,
				map[string]string{
					"choices":       utils.Stringify(orgProfiles),
					"SESSION_TOKEN": token,
					"continue":      continueUrl, //登录完成后跳转的地址，没有加token等参数
				}))
			return
		}
		exactAdmin := admin[0]
		ropToken = newToken(AdminIdentity{
			Iat:      now,
			Exp:      now.Add(utils.AdminTokenDuration),
			ZjuId:    exactAdmin.ZjuId,
			At:       exactAdmin.At,
			Nickname: exactAdmin.Nickname,
			Level:    exactAdmin.Level,
		})
	}
	ctx.Redirect(302, utils.AddQuery(continueUrl, map[string]string{"ropToken": ropToken}))
}

func authInit(routerGroup *gin.RouterGroup) {
	//定期清除不再需要的voidInfo
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

	routerGroup.GET("/loginByToken", loginByToken)
	routerGroup.GET("/logout", RequireLoginWithRefresh(false), logout)
	routerGroup.GET("/logoutAll", RequireLoginWithRefresh(false), logoutAll)
}
