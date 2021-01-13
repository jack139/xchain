package types

import (
	uuid "github.com/satori/go.uuid"
)

// 授权操作
// ToUserID 请求 FromUserID 授权，指定 DealID，进入 FromUserID 的链表
// FromUserID 加密数据 Data 后 返回 ToUserID，进入 ToUserID 的链表
type Auth struct {
	ID          uuid.UUID
	ReqID       uuid.UUID // 授权请求的ID， action==5时使用
	DealID      uuid.UUID // 交易ID
	FromUserID  []byte // 用户的加密公钥
	ToUserID    []byte // 被授权的用户的加密公钥
	Data        []byte //  action==4, 格式： rb.pub长度(byte) + rb.pub + sm2加密的rb.priv
			           //  action==5, 格式： ra.pub长度(byte) + ra.pub + K加密的数据
	Action      byte // 0x04 请求授权， 0x05 响应授权
}

// GetKey 获取实体键
func (auth *Auth) GetKey() []byte {
	return auth.ID.Bytes()
}

func (auth *Auth) getSignBytes() []byte {
	return auth.ID[:]
}