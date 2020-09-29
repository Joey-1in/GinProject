package controller

import (
	"os"
	"time"
	"strconv"
	"ginProject/service"
	"ginProject/tool"
	"ginProject/param"
	"ginProject/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type MenberController struct {

}

func (mc *MenberController) Router(engine *gin.Engine) {
	engine.GET("/api/sendCode", mc.sendCode)
	engine.POST("/api/loginSMs", mc.smsLogin)
	engine.GET("/api/captcha", mc.captcha)
	// postman验证码测试
	engine.POST("/api/vertifycha", mc.vertifycha)
	engine.POST("/api/login", mc.login)
	// 头像上传
	engine.POST("/api/upload/avatar", mc.uploadAvatar)
	//用户信息查询
	engine.GET("/api/userinfo", mc.userInfo)
}

// 发送短信验证码
func (mc *MenberController) sendCode(ctx *gin.Context) {
	// 获取手机号码
	phone, exist := ctx.GetQuery("phone")
	if !exist {
		tool.Failed(ctx, "参数解析失败")
		return
	}
	ms := service.MenberService{}
	isSend := ms.SendCode(phone)
	if isSend {
		tool.Success(ctx, "发送成功")
		return
	}
	tool.Failed(ctx, "发送失败")
	return
}

// 登录方式：手机号+短信登录
func (mc *MenberController) smsLogin(ctx *gin.Context) {
	var smsLogin param.SmsLogin
	// 解析请求参数
	err := tool.Decode(ctx.Request.Body, &smsLogin)
	if err != nil {
		tool.Failed(ctx, "参数解析失败")
		return
	}

	// 完成手机号+验证码登录
	us := service.MenberService{}
	member := us.SmsLogin(smsLogin)
	if member != nil {
		// 保存session
		sess, _ := json.Marshal(member)
		err = tool.SetSess(ctx, "user_"+string(member.Id), sess)
		if err != nil {
			tool.Failed(ctx, "登录失败")
		}
		// 设置cookie
		ctx.SetCookie("cookie_id", "login_success", 10*60, "/", "localohost", true, true)

		tool.Success(ctx, member)
		return
	}
	tool.Failed(ctx, "登录失败")
}

// 生成图形验证码返回客户端
func (mc *MenberController) captcha(ctx *gin.Context) {
	tool.GenerateCaptcha(ctx)
}

// 验证验证码是否正确
func (mc *MenberController) vertifycha(ctx *gin.Context) {
	var captcha tool.CaptchaResult
	err := tool.Decode(ctx.Request.Body, &captcha)
	if err != nil {
		tool.Failed(ctx, " 参数解析失败 ")
		return
	}

	result := tool.VertifyCaptcha(captcha.Id, captcha.VertifyValue)
	if result {
		tool.Success(ctx, "验证通过")
	} else {
		tool.Failed(ctx, " 验证失败 ")
	}
}

// 账户密码登录
func (mc *MenberController) login(ctx *gin.Context) {
	// 解析用户账户密码
	var loginParam param.LoginParam
	err := tool.Decode(ctx.Request.Body, &loginParam)
	if err != nil {
		tool.Failed(ctx, "参数解析失败")
	}
	// 验证验证码
	validate := tool.VertifyCaptcha(loginParam.Id, loginParam.Value)
	if validate {
		tool.Failed(ctx, "验证码不正确，请重新输入")
		return
	}
	// 登录
	ms := service.MenberService{}
	member := ms.Login(loginParam.Name, loginParam.Password)
	if member.Id != 0 {
		//用户信息保存到session
		sess, _ := json.Marshal(member)
		err = tool.SetSess(ctx, "user_"+string(member.Id), sess)
		if err != nil {
			tool.Failed(ctx, "登录失败")
			return
		}
		tool.Success(ctx, &member)
		return
	}
	tool.Failed(ctx, "登录失败")
	return
}

// 头像上传
func (mc *MenberController) uploadAvatar(ctx *gin.Context) {
	// 解析参数 file、userID
	userId := ctx.PostForm("userId")
	// 获取图像内容
	file, err := ctx.FormFile("avatar")
	if err != nil || userId == "" {
		tool.Failed(ctx, "参数解析失败")
		return 
	}
	// 判断是否登陆
	sess := tool.GetSess(ctx, "user_"+userId)
	if sess ==nil {
		tool.Failed(ctx, "参数不合法")
		return 
	}
	var member model.Menber
	json.Unmarshal(sess.([]byte), &member)
	// file保存到本地
	fileName := "./uploadfile" + strconv.FormatInt(time.Now().Unix(), 10) + file.Filename
	err = ctx.SaveUploadedFile(file, fileName)
	if err != nil {
		tool.Failed(ctx, "头像更新失败")
		return 
	}

	// // 将头像路劲保存到数据库
	// memberService := service.MenberService{}
	
	// path := memberService.UploadAvatar(member.Id, fileName[1:])
	// if path != "" {
	// 	tool.Success(ctx, "http://localhost:8080" + path)
	// 	return
	// }
	// tool.Failed(ctx, "上传失败")

	// 将文件上传到fastDFS系统
	fileId := tool.UploadFile(fileName)
	if fileId != "" {
		//删除本地uploadfile下的文件
		os.Remove(fileName)

		// http://localhost:8080/static/.../davie.png
		//4、将保存后的文件本地路径 保存到用户表中的头像字段
		memberService := service.MenberService{}
		path := memberService.UploadAvatar(member.Id, fileId)
		if path != "" {
			tool.Success(ctx, tool.FileServerAddr()+"/"+path)
			return
		}
	}
	// 返回结果
	tool.Failed(ctx, "上传失败")
}

// 查询用户信息
 func (mc *MenberController) userInfo(context *gin.Context) {
	cookie, err := tool.CookieAuth(context)
	if err != nil {
		context.Abort()
		tool.Failed(context, "还未登录，请先登录")
		return
	}

	memberService := service.MenberService{}
	member := memberService.GetUserInfo(cookie.Value)
	if member != nil {
		//返回成功信息给客户端
		tool.Success(context, map[string]interface{}{
			"id":            member.Id,
			"user_name":     member.UserName,
			"mobile":        member.Mobile,
			"register_time": member.RegisterTime,
			"avatar":        member.Avatar,
			"balance":       member.Balance,
			"city":          member.City,
		})
		return
	}
	tool.Failed(context, "获取用户信息失败")
}