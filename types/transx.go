package types

import (
	"time"
	"crypto/rand"
	"github.com/tjfoc/gmsm/sm2"
)

// IPayload 接口
type IPayload interface {
	getSignBytes() []byte
	GetKey() []byte
}

// Transx 事务基类
type Transx struct {
	Signature  []byte //发送方对这个消息的私钥签名
	SendTime   *time.Time
	SignPubKey sm2.PublicKey
	Payload    IPayload
}

/*
	Sign 给消息签名
	privKey:发送方私钥
*/
func (cmu *Transx) Sign(privKey sm2.PrivateKey) error {
	bz := cmu.Payload.getSignBytes()
	sig, err := privKey.Sign(rand.Reader, bz, nil) 
	cmu.Signature = sig
	return err
}

/*
	Verify 校验消息是否未被篡改
*/
func (cmu *Transx) Verify() bool {
	if cmu.Payload==nil { // Transx 格式不对
		return false
	}
	data := cmu.Payload.getSignBytes()
	sig := cmu.Signature
	rslt := cmu.SignPubKey.Verify(sig, data)
	return rslt
}

