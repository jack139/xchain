package client

import (
	"xchain/types"

	"bytes"
	"math/big"
	"fmt"
	"io/ioutil"
	"context"
	"crypto/rand"
	"encoding/base64"
	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/os"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
)

// KEYFILENAME 私钥文件名
const KEYFILENAME string = "user.key"

var (
	cli *rpcclient.HTTP
	cdc = types.AminoCdc
	ctx = context.Background()
)

func init() {
	addr := cfg.DefaultRPCConfig().ListenAddress
	cli, _ = rpcclient.New(addr, "/websocket")
}

type cryptoPair struct {
	PrivKey *[]byte
	PubKey  *[]byte
}

type User struct {
	SignKey    sm2.PrivateKey `json:"sign_key"` // 节点私钥，用户签名
	CryptoPair cryptoPair     // 密钥协商使用
}

// 生成用户环境
func GetMe(path string) (*User, error) {
	keyFilePath := path + "/" + KEYFILENAME
	if cmn.FileExists(keyFilePath) {
		fmt.Printf("Found keyfile: %s\n", keyFilePath)
		uk, err := loadUserKey(keyFilePath)
		if err != nil {
			return nil, err
		}
		return uk, nil
	}

	return nil, fmt.Errorf("Keyfile does not exist!")
}

// 从文件装入key
func GenUserKey(path string) (*User, error) {
	keyFilePath := path + "/" + KEYFILENAME
	if cmn.FileExists(keyFilePath) {
		return nil, fmt.Errorf("Keyfile already exists!")
	}

	// 建目录
	if err := cmn.EnsureDir(path, 0700); err != nil {
		return nil, err
	}
	// 生成新的密钥文件
	fmt.Println("Make new key file: " + keyFilePath)	
	uk := new(User)
	signKey, err := sm2.GenerateKey(rand.Reader) // 生成密钥对
	if err != nil {
		return nil, err
	}
	uk.SignKey = *signKey
	pubKey := sm2.Compress(&uk.SignKey.PublicKey) // 33 bytes
	priKey := uk.SignKey.D.Bytes()

	uk.CryptoPair = cryptoPair{PrivKey: &priKey, PubKey: &pubKey}
	jsonBytes, err := cdc.MarshalJSON(uk.CryptoPair)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(keyFilePath, jsonBytes, 0644)
	if err != nil {
		return nil, err
	}
	return uk, nil
}

// 从 byte 恢复密钥对
func restoreKey(priv *[]byte) *sm2.PrivateKey {
	curve := sm2.P256Sm2()
	key := new(sm2.PrivateKey)
	key.PublicKey.Curve = curve
	key.D = new(big.Int).SetBytes(*priv)
	key.PublicKey.X, key.PublicKey.Y = curve.ScalarBaseMult(*priv)
	return key
}

// 从 byte 恢复公钥
func restorePublicKey(public []byte) *sm2.PublicKey {
	key := sm2.Decompress(public)
	return key
}

// 从文件导入用户密钥
func loadUserKey(keyFilePath string) (*User, error) {
	jsonBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	uk := new(User)
	err = cdc.UnmarshalJSON(jsonBytes, &uk.CryptoPair)
	if err != nil {
		return nil, fmt.Errorf("Error reading UserKey from %v: %v", keyFilePath, err)
	}
	// 恢复 privateKey
	uk.SignKey = *restoreKey(uk.CryptoPair.PrivKey)

	return uk, nil
}


// 交易结构转换为返回值的结构，能解密就解密
func txToResp(me *User, tx *types.Transx) *map[string]interface{} {
	var respMap = make(map[string]interface{})

	auth, ok := (*tx).Payload.(*types.Auth)	// 授权块
	if ok {
		//fmt.Printf("auth => %v\n", auth)
		var data string

		if auth.Action==0x05 { // 授权响应，则尝试解密 data
			// data 默认返回不解密的 base64
			data = base64.StdEncoding.EncodeToString(auth.Data) // 加密数据的 base64

			for {
				// 从请求的auth.Data里取到rb私钥
				// 获取 authID 对应的 授权请求 块
				addr, _ := cdc.MarshalJSON(auth.FromUserID) // trick: 应该是自己的，这里方便下面的参数只用 _
				authReqTx, err := queryTx(addr, "_", auth.ReqID.String())
				if err != nil {
					fmt.Println(err)
					break
				}
				if authReqTx==nil {
					fmt.Println("AuthID not found")
					break
				}
				authReq, ok := (*authReqTx).Payload.(*types.Auth)	// 授权块
				if !ok {
					fmt.Println("need a Auth Payload")
					break
				}

				// 解密出rb, 加密密钥使用私钥前16字节（128bit）
				rbBytesLen := int(authReq.Data[0])
				// 加密数据, rb私钥sm2加密
				rbPrivBytes, err := sm2.DecryptAsn1(&me.SignKey, authReq.Data[rbBytesLen+1:])
				if err != nil {
					fmt.Printf("sm2 decrypt fail: %s\n", err)
					break
				}
				rbPriv := restoreKey(&rbPrivBytes)

				// 私钥
				dbPriv := me.SignKey

				// auth.Data里取得密钥协商数据: ra.pub da.pub data
				raBytesLen := int(auth.Data[0])
				daBytes := auth.FromUserID
				raBytes := auth.Data[1:raBytesLen+1]

				daPub := restorePublicKey(daBytes)
				raPub := restorePublicKey(raBytes)

				// 生成 解密密钥
				decryptKey, _, _, err := sm2.KeyExchangeB(16, 
					auth.FromUserID, auth.ToUserID, &dbPriv, daPub, rbPriv, raPub)
				if err != nil {
					fmt.Printf("decryptKey error: %s", err)
					break
				}

				// 解密
				decrypted, err := sm4.Sm4CFB(decryptKey, auth.Data[raBytesLen+1:], false)
				if err==nil {
					data = string(decrypted)
				} 

				break
			}
		}

		userId, _ := cdc.MarshalJSON(auth.FromUserID)
		userId2, _ := cdc.MarshalJSON(auth.ToUserID)
		respMap["type"] = "AUTH"
		respMap["id"]  = auth.ID.String()
		respMap["from_user_id"]  = string(userId[1 : len(userId)-1]) // 去掉两边引号
		respMap["to_user_id"]  = string(userId2[1 : len(userId2)-1])
		respMap["data"]  = data
		respMap["deal_id"]  = auth.DealID.String() // 返回dealID
		respMap["action"]  = auth.Action
		respMap["send_time"]  = *(*tx).SendTime
		if auth.Action==5 {
			respMap["req_id"]  = auth.ReqID.String() // 授权请求的ID
		}
		return &respMap

	} else { // category == deal
		deal, ok := (*tx).Payload.(*types.Deal) // 交易块
		if ok {
			//fmt.Printf("deal => %v\n", deal)

			var data string
			
			// 尝试解密 data
			if bytes.Compare(deal.UserID, *me.CryptoPair.PubKey)==0 { // 是自己的交易, 进行解密
				// 从data中解密出sm4密钥
				// data格式： sm4密钥长度(byte)(sm2加密) + sm4加密的data
				sm4keyLen := int(deal.Data[0])

				decryptKey, err := sm2.DecryptAsn1(&me.SignKey, deal.Data[1:sm4keyLen+1])
				if err == nil {
					decrypted, err := sm4.Sm4CFB(decryptKey, deal.Data[sm4keyLen+1:], false)
					if err==nil {
						data = string(decrypted)
					} else {
						data = base64.StdEncoding.EncodeToString(deal.Data) // 加密数据的 base64
						fmt.Printf("sm4 decryption error: %s\n", err)
					}
				} else {
					fmt.Printf("sm2 decrypt fail: %s\n", err)
				}
			} else {
				data = base64.StdEncoding.EncodeToString(deal.Data) // 加密数据的 base64
			}

			userId, _ := cdc.MarshalJSON(deal.UserID)
			respMap["type"] = "DEAL"
			respMap["id"] = deal.ID.String()
			respMap["user_id"] = string(userId[1 : len(userId)-1]) // 去掉两边引号
			respMap["data"] = data
			respMap["send_time"] = *(*tx).SendTime
			return &respMap

		}
	}

	return nil
}