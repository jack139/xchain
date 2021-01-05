package main

import (
	"crypto/rand"
	"fmt"
	//"io/ioutil"
	"math/big"
	//"os"
	"encoding/base64"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
)

/*

标准流程：
 A: 
 	生成 rA
 	rA.pub dA.pub --> B

 B: 
 	校验 rA
 	生成 rB
 	rB.pub dB.pub --> A

 A:
 	校验 rB
 	生成 key
 	加密数据 data
 	data --> B

 B: 
 	生成 key
 	解密数据 data

简化流程：
B：
 	生成 rB
 	rB.pub dB.pub --> A

A: 
  	校验 rB
 	生成 rA
 	生成 key
 	加密数据 data
 	rA.pub dA.pub data --> B

B: 
 	生成 key
 	解密数据 data

*/

var (
	// 用户 id
	ida = []byte("test1")
	idb = []byte("test2")

	// 用户密钥
	daPriStr = string("9nJzj0RrtDhYd3LzfbiIOL3LHryYYhJaxGozP001540=")
	dbPriStr = string("39WSsgO/Kp9QUiQj7STJ5z/X6Vg4L78blSi4sjQIkVQ=")

	// 临时密钥
	ra, rb *sm2.PrivateKey
)


// 从 base64私钥 恢复密钥对
func restoreKey(privStr string) *sm2.PrivateKey {
	priv, _  := base64.StdEncoding.DecodeString(privStr)

	curve := sm2.P256Sm2()
	key := new(sm2.PrivateKey)
	key.PublicKey.Curve = curve
	key.D = new(big.Int).SetBytes(priv)
	key.PublicKey.X, key.PublicKey.Y = curve.ScalarBaseMult(priv)
	return key
}

// 从 base64 恢复公钥
func restorePublicKey(pubStr string) *sm2.PublicKey {
	public, _  := base64.StdEncoding.DecodeString(pubStr)

	key := sm2.Decompress(public)
	return key
}

/*
B：
 	生成 rB
 	rB.pub dB.pub --> A
*/
func step_1_B() (string, string, error){
	// 恢复 dB
	db := restoreKey(dbPriStr)
	// 生成 rB
	rb, _ = sm2.GenerateKey(rand.Reader) // 生成密钥对

	return base64.StdEncoding.EncodeToString(sm2.Compress(&db.PublicKey)),
		base64.StdEncoding.EncodeToString(sm2.Compress(&rb.PublicKey)),
		nil
}

/*
A: 
  	校验 rB.pub
 	生成 rA
 	生成 key
 	加密数据 data
 	rA.pub dA.pub data --> B
*/
func step_2_A(dbPubStr string, rbPubStr string) (string, string, []byte, error){
	// 恢复 dA
	da := restoreKey(daPriStr)
	// 生成 rA
	ra, _ = sm2.GenerateKey(rand.Reader) // 生成密钥对


	// 验证 rb.pub
	dbPub := restorePublicKey(dbPubStr)
	rbPub := restorePublicKey(rbPubStr)

	//  生成 key
	key, _, _, err := sm2.KeyExchangeA(16, ida, idb, da, dbPub, ra, rbPub)
	if err != nil {
		return "", "", nil, err
	}

	fmt.Printf("%v\n", key)

	// 加密数据
	encrypted, err := sm4.Sm4CFB(key, []byte("hello world"), true)
	if err != nil {
		return "", "", nil, fmt.Errorf("sm4 encrypt error: %s", err)
	}
	fmt.Printf("encrypted --> %x\n", encrypted)

	return base64.StdEncoding.EncodeToString(sm2.Compress(&da.PublicKey)),
		base64.StdEncoding.EncodeToString(sm2.Compress(&ra.PublicKey)),
		encrypted, nil
}

/*
B: 
 	生成 key
 	解密数据 data

*/
func step_3_B(daPubStr string, raPubStr string, data []byte) error {
	// 恢复 dB
	db := restoreKey(dbPriStr)

	// 验证 ra.pub
	daPub := restorePublicKey(daPubStr)
	raPub := restorePublicKey(raPubStr)

	//  生成 key
	key, _, _, err := sm2.KeyExchangeB(16, ida, idb, db, daPub, rb, raPub)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", key)

	// 解密数据
	plain, err := sm4.Sm4CFB(key, data, false)
	if err != nil {
		return fmt.Errorf("sm4 decrypt error: %s", err)
	}
	fmt.Printf("plain --> %s\n", plain)


	return nil
}


func main() {

	/*
	priv, _ := sm2.GenerateKey(rand.Reader) // 生成密钥对
	pubBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	fmt.Printf("private: %s\npublic: %s\n",
		base64.StdEncoding.EncodeToString(priv.D.Bytes()),
		base64.StdEncoding.EncodeToString(pubBytes),
	)
	*/


	dbPub, rbPub, err := step_1_B()
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%s\n%s\n", dbPub, rbPub)	

	daPub, raPub, data, err := step_2_A(dbPub, rbPub)
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%s\n%s\n", daPub, raPub)

	err = step_3_B(daPub, raPub, data)
	if err!=nil{
		panic(err)
	}

	fmt.Println("pass")
}

