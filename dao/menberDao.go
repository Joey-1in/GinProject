package dao

import (
	"ginProject/tool"
	"ginProject/model"
	"github.com/google/logger"
)

type MenberDao struct {
	*tool.Orm
}

// 保存短信登录验证码到数据库
func (md *MenberDao) InsertCode(sms model.SmsCode) int64 {
	result, err := md.InsertOne(&sms)
	if err != nil {
		logger.Error(err.Error())
	}
	return result
}

// 验证手机号和验证码是否存在
func (md *MenberDao) ValidateSmsCode(phone, code string) *model.SmsCode {
	var sms model.SmsCode 
	// 查询
	_, err := md.Where("phone = ? and code = ?", phone, code).Get(&sms)
	if err != nil {
		logger.Error(err.Error())
	}
	return &sms
}

func (md *MenberDao) QueryByPhone(phone string) *model.Menber {
	var member model.Menber
	if _, err:= md.Where("mobile = ?", phone).Get(&member); err != nil {
		logger.Error(err.Error())
	}
	return &member
}

func (md *MenberDao) InsertMember(member *model.Menber) int64 {
	result, err := md.InsertOne(&member)
	if err != nil {
		logger.Error(err.Error())
		return 0
	}
	return result
}

func (md *MenberDao) Query(name, password string) *model.Menber {
	var member model.Menber
	password = tool.EncoderSha256(password)
	_, err := md.Where("user_name = ? and password = ?", name, password).Get(&member)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}
	return &member
}

func (md *MenberDao) UpdateMemberAvatar(userId int64, fileName string) int64 {
	member := model.Menber{Avatar: fileName}
	result, err := md.Where("id = ?", userId).Update(&member)
	if err != nil {
		logger.Error(err.Error())
		return 0
	} 
	return result
}

func (md *MenberDao) QueryMemberById(id int64) *model.Menber {
	var member model.Menber
	if _, err := md.Orm.Where(" id = ? ", id).Get(&member); err != nil {
		return nil
	}
	return &member
}