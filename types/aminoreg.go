package types

import (
	"crypto/elliptic"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
	"github.com/tjfoc/gmsm/sm2"
)

// AminoCdc amino编码类
var AminoCdc = amino.NewCodec()

func init() {
	AminoCdc.RegisterInterface((*IPayload)(nil), nil)
	AminoCdc.RegisterConcrete(&Deal{}, "deal", nil)
	AminoCdc.RegisterConcrete(&Auth{}, "auth", nil)
	AminoCdc.RegisterConcrete(sm2.PublicKey{}, "sm2/pubkey", nil)
	AminoCdc.RegisterConcrete(sm2.PrivateKey{}, "sm2/privkey", nil)
	AminoCdc.RegisterInterface((*elliptic.Curve)(nil), nil)
	// 不注册这个，对Block进行序列化时会报错：Unregistered interface types.Evidence
	AminoCdc.RegisterInterface((*types.Evidence)(nil), nil)
}
