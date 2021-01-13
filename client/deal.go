package client

import (
	"xchain/types"

	"fmt"
	"time"
	"encoding/json"
	"crypto/rand"
	uuid "github.com/satori/go.uuid"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
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
func (me *User) Deal(data string) ([]byte, error) {

	// 用户id
	userId := *me.CryptoPair.PubKey

	// sm4密钥随机产生16字节（128bit）
	encryptKey := make([]byte, 16)
	_, err := rand.Reader.Read(encryptKey)
	if err != nil {
		return nil, fmt.Errorf("sm4 encryptKey error: %s", err)
	}

	// sm4加密数据
	encrypted, err := sm4.Sm4CFB(encryptKey, []byte(data), true)
	if err != nil {
		return nil, fmt.Errorf("sm4 encrypt error: %s", err)
	}
	//fmt.Printf("encrypted --> %x\n", encrypted)

	// 用sm2加密sm4密钥
	d0, err := sm2.EncryptAsn1(&me.SignKey.PublicKey, encryptKey, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("sm2 encrypt fail: %s", err)
	}

	// data格式： sm4密钥长度(byte)(sm2加密) + sm2加密的sm4密钥 + sm4加密的data
	cryptData := append([]byte{byte(len(d0))}, d0...)
	cryptData = append(cryptData, encrypted...)


	now := time.Now()
	tx := new(types.Transx)
	tx.SendTime = &now

	deal := types.Deal{
		ID:         uuid.NewV4(),
		UserID:     userId,
		Data:       cryptData,
	}

	tx.Payload = &deal

	tx.Sign(me.SignKey)

	// broadcast this tx
	bz, err := cdc.MarshalJSON(&tx)
	if err != nil {
		return nil, err
	}

	ret, err := cli.BroadcastTxSync(ctx, bz)
	if err != nil {
		return nil, err
	}

	fmt.Printf("deal => %+v\n", ret)

	if ret.Code !=0 {
		return nil, fmt.Errorf(ret.Log)
	}

	respMap := map[string]interface{}{"id" : deal.ID.String()}

	// 返回结果转为json
	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
