package service

import (
	"errors"
	"scarlet/common"
	"scarlet/tool"
	"time"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return new(TokenService)
}

// Token的生成之后可以替换成 某些算法
func (this *TokenService) Login(user common.User) string {
	token := this.GenToken(user)
	redisClient.Set(token, user.ID, time.Hour)
	return token
}

func (this *TokenService) Logout(token string) {
	redisClient.Del(token)
	return
}

func (this *TokenService) GenToken(user common.User) string {
	return tool.Md5(time.Now().String() + user.Name)
}

func (this *TokenService) Get(token string) (uint, error) {
	v := redisClient.Get(token)
	val, _ := v.Uint64()
	if val != 0 {
		return uint(val), nil
	} else {
		return 0, errors.New("User Not Found")
	}
}
