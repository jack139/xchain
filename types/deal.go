package types

import (
	uuid "github.com/satori/go.uuid"
)

const (
	Buy byte = 0x01
	Sell byte = 0x02
	ChangeOwner byte = 0x03
)


// 交易信息
type Deal struct {
	ID         uuid.UUID // 交易ID
	UserID     []byte //用户的加密公钥
	Data       []byte // 加密交易数据（例如 ipfs hash）
}

// GetKey 获取实体键
func (deal *Deal) GetKey() []byte {
	return deal.ID.Bytes()
}

func (deal *Deal) getSignBytes() []byte {
	return deal.ID[:]
}
