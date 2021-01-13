package types

import (
	uuid "github.com/satori/go.uuid"
)


// 交易信息
type Deal struct {
	ID     uuid.UUID // 交易ID
	UserID []byte //用户的加密公钥
	Data   []byte // 格式： sm4密钥长度(byte)(sm2加密) + sm2加密的sm4密钥 + sm4加密的data
}

// GetKey 获取实体键
func (deal *Deal) GetKey() []byte {
	return deal.ID.Bytes()
}

func (deal *Deal) getSignBytes() []byte {
	return deal.ID[:]
}
