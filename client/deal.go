package client

import (
	"xchain/types"

	"fmt"
	"io"
	"time"
	crypto_rand "crypto/rand"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/nacl/box"
)

func isASCII(s string) bool {
    for i := 0; i < len(s); i++ {
        if s[i] <= 32 || s[i] >= 127 {
            return false
        }
    }
    return true
}

// 交易上链，数据加密
func (me *User) Deal(data string) error {

	// 用户id
	userId := *me.CryptoPair.PubKey

	sharedEncryptKey := new([32]byte)
	box.Precompute(sharedEncryptKey, &userId, me.CryptoPair.PrivKey)

	var nonce [24]byte
	if _, err := io.ReadFull(crypto_rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	//fmt.Printf("data=>%v,nonce=>%v,sharedEncryptKey=>%v\n", data, nonce, *sharedEncryptKey)
	encrypted := box.SealAfterPrecomputation(nonce[:], []byte(data), &nonce, sharedEncryptKey)

	now := time.Now()
	tx := new(types.Transx)
	tx.SendTime = &now

	deal := types.Deal{
		ID:         uuid.NewV4(),
		UserID:     userId,
		Data:       encrypted,
	}

	tx.Payload = &deal

	tx.Sign(me.SignKey)
	tx.SignPubKey = me.SignKey.PubKey()
	// broadcast this tx
	bz, err := cdc.MarshalJSON(&tx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ret, err := cli.BroadcastTxSync(ctx, bz)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("deal => %+v\n", ret)

	// ret  *ctypes.ResultBroadcastTxCommit
	if ret.Code !=0 {
		fmt.Println(ret.Log)
		return fmt.Errorf(ret.Log)
	}

	return nil
}
