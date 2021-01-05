package types

import (
	uuid "github.com/satori/go.uuid"
)

const (
	AuthQuery byte = 0x04
	DoQuery byte = 0x05
	DeAuthQuery byte = 0x06
)

// 授权操作
// ToUserID 请求 FromUserID 授权，指定 DealID，进入 FromUserID 的链表
// FromUserID 加密数据 Data 后 返回 ToUserID，进入 ToUserID 的链表
type Auth struct {
	ID             uuid.UUID
	DealID         uuid.UUID // 交易ID
	FromUserID     []byte // 用户的加密公钥
	ToUserID       []byte // 被授权的用户的加密公钥
	Data           []byte // FromUser加密数据，被授权者ToUserID可以解密
	Action         byte // 0x04 请求授权， 0x05 响应授权
}

// GetKey 获取实体键
func (auth *Auth) GetKey() []byte {
	return auth.ID.Bytes()
}

func (auth *Auth) getSignBytes() []byte {
	return auth.ID[:]
}