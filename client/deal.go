package client

import (
	"xchain/types"

	"fmt"
	"time"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
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

	// 加密密钥使用私钥前16字节（128bit）
	encryptKey := (*me.CryptoPair.PrivKey)[:16]

	// 加密数据
	encrypted, err := sm4.Sm4CFB(encryptKey, []byte(data), true)
	if err != nil {
		return nil, fmt.Errorf("sm4 encrypt error: %s", err)
	}
	//fmt.Printf("encrypted --> %x\n", encrypted)

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

	respMap := map[string]string{"id" : deal.ID.String()}

	// 返回结果转为json
	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
