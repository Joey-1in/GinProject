package service

import (
	"fmt"
	"time"
	"strconv"
	"math/rand"
	"encoding/json"
	"ginProject/dao"
	"ginProject/tool"
	"ginProject/param"
	"ginProject/model"
	"github.com/google/logger"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type MenberService struct {

}

func (ms *MenberService) SendCode(phone string) bool {
	// 生成验证码
	code := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	// 调用阿里云sdk
	config := tool.GetConfig().Sms
	client, err := dysmsapi.NewClientWithAccessKey(config.RegionId, config.AppKey, config.AppSecret)
	if err != nil {
		logger.Error(err.Error())
		return false
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = config.SignName
	request.TemplateCode = config.TemplateCode
	request.PhoneNumbers = phone
	par, err := json.Marshal(map[string]interface{}{
		"code": code,
	})

	request.TemplateParam = string(par)

	response, err := client.SendSms(request)

	if err != nil {
		logger.Error(err.Error())
		return false
	}
	// 接受返回结果判断发送状态
	if response.Code == "OK" {
		//将验证码保存到数据库中
		smsCode := model.SmsCode{
			Phone: phone, 
			Code: code, 
			BizId: response.BizId, 
			CreateTime: time.Now().Unix(),
		}
		memberDao := dao.MenberDao{ tool.DbEngine }
		result := memberDao.InsertCode(smsCode)
		return result > 0
	}
	return false
}

// 成功时就返回用户的结构体
func (ms *MenberService) SmsLogin(loginParam param.SmsLogin) *model.Menber {
	// 获取手机号和验证码
	// 验证手机号和验证码是否正确
	md := dao.MenberDao{ tool.DbEngine }
	sms := md.ValidateSmsCode(loginParam.Phone, loginParam.Code)
	if sms.Id == 0 {
		return nil
	}
	// 根据手机号member表中查询记录
	member := md.QueryByPhone(loginParam.Phone)
	if member.Id != 0 {
		return member
	}
	// 新创建一个member记录并保存
	user := model.Menber{}
	user.UserName = loginParam.Phone
	user.Mobile = loginParam.Phone
	user.RegisterTime = time.Now().Unix()
	user.Id = md.InsertMember(&user)
	return &user
}

// 账号密码登录
func (ms *MenberService) Login(name, password string) *model.Menber {
	// 查询用户信息，直接返回
	md := dao.MenberDao{ tool.DbEngine}
	member := md.Query(name, password)
	if member.Id != 0 {
		return member
	}
	// 不存在则作为新用户保存到数据库中
	user := model.Menber{}
	user.UserName = name
	user.Password = tool.EncoderSha256(password)
	user.Mobile = name
	user.RegisterTime = time.Now().Unix()

	result := md.InsertMember(&user)
	user.Id = result
	return &user
}

// 更新头像
func (ms *MenberService) UploadAvatar(userId int64, fileName string) string {
	memberDao := dao.MenberDao{tool.DbEngine}
	result := memberDao.UpdateMemberAvatar(userId, fileName)
	if result == 0 {
		return ""
	}
	return fileName
}

// 验证cookie是否登录,返回用户信息
func (ms *MenberService) GetUserInfo(userId string) *model.Menber {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil
	}
	memberDao := dao.MenberDao{ tool.DbEngine }
	return memberDao.QueryMemberById(int64(id))
}