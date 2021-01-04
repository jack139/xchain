package types

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/types"
)

// AminoCdc amino编码类
var AminoCdc = amino.NewCodec()

func init() {
	AminoCdc.RegisterInterface((*IPayload)(nil), nil)
	AminoCdc.RegisterConcrete(&Deal{}, "deal", nil)
	AminoCdc.RegisterConcrete(&Auth{}, "auth", nil)
	AminoCdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	AminoCdc.RegisterConcrete(ed25519.PubKey{}, "ed25519/pubkey", nil)
	AminoCdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	AminoCdc.RegisterConcrete(ed25519.PrivKey{}, "ed25519/privkey", nil)
	// 不注册这个，对Block进行序列化时会报错：Unregistered interface types.Evidence
	AminoCdc.RegisterInterface((*types.Evidence)(nil), nil)
}
