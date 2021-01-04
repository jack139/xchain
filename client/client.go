package client

import (
	"xchain/types"

	"fmt"
	"io/ioutil"
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"

	cfg "github.com/tendermint/tendermint/config"
	cmn "github.com/tendermint/tendermint/libs/os"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"golang.org/x/crypto/nacl/box"
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
	PrivKey *[32]byte
	PubKey  *[32]byte
}

type User struct {
	SignKey    crypto.PrivKey `json:"sign_key"` // 节点私钥，用户签名
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
	uk.SignKey = ed25519.GenPrivKey()
	pubKey, priKey, err := box.GenerateKey(crypto_rand.Reader)
	if err != nil {
		return nil, err
	}
	uk.CryptoPair = cryptoPair{PrivKey: priKey, PubKey: pubKey}
	jsonBytes, err := cdc.MarshalJSON(uk)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(keyFilePath, jsonBytes, 0644)
	if err != nil {
		return nil, err
	}
	return uk, nil
}

func loadUserKey(keyFilePath string) (*User, error) {
	//copy(privKey[:], bz)
	jsonBytes, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	uk := new(User)
	err = cdc.UnmarshalJSON(jsonBytes, uk)
	if err != nil {
		return nil, fmt.Errorf("Error reading UserKey from %v: %v", keyFilePath, err)
	}
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
			var decryptKey, publicKey [32]byte

			publicKey = auth.FromUserID

			// 解密 data 数据
			box.Precompute(&decryptKey, &publicKey, me.CryptoPair.PrivKey)
			var decryptNonce [24]byte
			copy(decryptNonce[:], auth.Data[:24])
			//fmt.Printf("data=>%v,decryptNonce=>%v,decryptKey=>%v\n", deal.Data[24:], decryptNonce, decryptKey)
			decrypted, ok := box.OpenAfterPrecomputation(nil, auth.Data[24:], &decryptNonce, &decryptKey)
			if ok {
				data = string(decrypted)
			} else {
				data = base64.StdEncoding.EncodeToString(auth.Data) // 加密数据的 base64
				fmt.Println("decryption error")
			}
		}

		userId, _ := cdc.MarshalJSON(auth.FromUserID)
		userId2, _ := cdc.MarshalJSON(auth.ToUserID)
		respMap["type"] = "AUTH"
		respMap["id"]  = auth.ID.String()
		respMap["user_id"]  = string(userId[1 : len(userId)-1]) // 去掉两边引号
		respMap["auth_user_id"]  = string(userId2[1 : len(userId2)-1])
		respMap["data"]  = data
		respMap["deal_id"]  = auth.DealID.String() // 返回dealID
		respMap["action"]  = auth.Action
		respMap["send_time"]  = *(*tx).SendTime
		return &respMap

	} else { // category == deal
		deal, ok := (*tx).Payload.(*types.Deal) // 交易块
		if ok {
			//fmt.Printf("deal => %v\n", deal)

			var data string
			
			// 尝试解密 data
			var decryptKey, publicKey [32]byte

			if deal.UserID==*me.CryptoPair.PubKey { // 是自己的交易, 进行解密
				publicKey = deal.UserID

				// 解密 data 数据
				box.Precompute(&decryptKey, &publicKey, me.CryptoPair.PrivKey)
				var decryptNonce [24]byte
				copy(decryptNonce[:], deal.Data[:24])
				//fmt.Printf("data=>%v,decryptNonce=>%v,decryptKey=>%v\n", deal.Data[24:], decryptNonce, decryptKey)
				decrypted, ok := box.OpenAfterPrecomputation(nil, deal.Data[24:], &decryptNonce, &decryptKey)
				if ok {
					data = string(decrypted)
				} else {
					data = base64.StdEncoding.EncodeToString(deal.Data) // 加密数据的 base64
					fmt.Println("decryption error")
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