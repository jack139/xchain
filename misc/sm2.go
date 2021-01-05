package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	//"io/ioutil"
	"math/big"
	//"os"
	"encoding/base64"
	"github.com/tjfoc/gmsm/sm2"
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
	daPubStr = string("gOh55IGeQ1dey1KTw5Gv01sb7n7uuZrmuBGg9oz34cM4wqsdAYTaJFg8Gl/toflrVe6mJswJtnvkIQcNe3GFXg==")
	dbPriStr = string("39WSsgO/Kp9QUiQj7STJ5z/X6Vg4L78blSi4sjQIkVQ=")
	dbPubStr = string("LQkeG2TpMARDnMu5HkSESTBq6oJMbGCEt0qenxvgwZXZzPl/zU5xycq+kcp5ZkXGIrO14lvVINMUakRYBUe/3w==")

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

	curve := sm2.P256Sm2()
	key := new(sm2.PublicKey)
	key.Curve = curve
	key.X = new(big.Int).SetBytes(public[:32])
	key.Y = new(big.Int).SetBytes(public[32:])
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

	dbPubBytes := append(db.PublicKey.X.Bytes(), db.PublicKey.Y.Bytes()...)
	rbPubBytes := append(rb.PublicKey.X.Bytes(), rb.PublicKey.Y.Bytes()...)
	return base64.StdEncoding.EncodeToString(dbPubBytes),
		base64.StdEncoding.EncodeToString(rbPubBytes),
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

	daPubBytes := append(da.PublicKey.X.Bytes(), da.PublicKey.Y.Bytes()...)
	raPubBytes := append(ra.PublicKey.X.Bytes(), ra.PublicKey.Y.Bytes()...)
	return base64.StdEncoding.EncodeToString(daPubBytes),
		base64.StdEncoding.EncodeToString(raPubBytes),
		key, nil
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

	if bytes.Compare(key, data) != 0 {
		return fmt.Errorf("key exchange differ")
	}

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
	fmt.Printf("%s\n%s\n", dbPub, rbPub)	

	daPub, raPub, data, err := step_2_A(dbPub, rbPub)
	fmt.Printf("%s\n%s\n", daPub, raPub)

	err = step_3_B(daPub, raPub, data)
	if err!=nil{
		panic(err)
	}

	fmt.Println("pass")
}

func exchangeTest() {
	ida := []byte("test1")
	idb := []byte("test2")
	daBuf := []byte{0x81, 0xEB, 0x26, 0xE9, 0x41, 0xBB, 0x5A, 0xF1,
		0x6D, 0xF1, 0x16, 0x49, 0x5F, 0x90, 0x69, 0x52,
		0x72, 0xAE, 0x2C, 0xD6, 0x3D, 0x6C, 0x4A, 0xE1,
		0x67, 0x84, 0x18, 0xBE, 0x48, 0x23, 0x00, 0x29}
	dbBuf := []byte{0x78, 0x51, 0x29, 0x91, 0x7D, 0x45, 0xA9, 0xEA,
		0x54, 0x37, 0xA5, 0x93, 0x56, 0xB8, 0x23, 0x38,
		0xEA, 0xAD, 0xDA, 0x6C, 0xEB, 0x19, 0x90, 0x88,
		0xF1, 0x4A, 0xE1, 0x0D, 0xEF, 0xA2, 0x29, 0xB5}
	raBuf := []byte{0xD4, 0xDE, 0x15, 0x47, 0x4D, 0xB7, 0x4D, 0x06,
		0x49, 0x1C, 0x44, 0x0D, 0x30, 0x5E, 0x01, 0x24,
		0x00, 0x99, 0x0F, 0x3E, 0x39, 0x0C, 0x7E, 0x87,
		0x15, 0x3C, 0x12, 0xDB, 0x2E, 0xA6, 0x0B, 0xB3}
	rbBuf := []byte{0x7E, 0x07, 0x12, 0x48, 0x14, 0xB3, 0x09, 0x48,
		0x91, 0x25, 0xEA, 0xED, 0x10, 0x11, 0x13, 0x16,
		0x4E, 0xBF, 0x0F, 0x34, 0x58, 0xC5, 0xBD, 0x88,
		0x33, 0x5C, 0x1F, 0x9D, 0x59, 0x62, 0x43, 0xD6}

	//expk := []byte{0x6C, 0x89, 0x34, 0x73, 0x54, 0xDE, 0x24, 0x84,
	//	0xC6, 0x0B, 0x4A, 0xB1, 0xFD, 0xE4, 0xC6, 0xE5}

	curve := sm2.P256Sm2()
	//curve.ScalarBaseMult(daBuf)
	da := new(sm2.PrivateKey)
	da.PublicKey.Curve = curve
	da.D = new(big.Int).SetBytes(daBuf)
	da.PublicKey.X, da.PublicKey.Y = curve.ScalarBaseMult(daBuf)

	db := new(sm2.PrivateKey)
	db.PublicKey.Curve = curve
	db.D = new(big.Int).SetBytes(dbBuf)
	db.PublicKey.X, db.PublicKey.Y = curve.ScalarBaseMult(dbBuf)

	ra := new(sm2.PrivateKey)
	ra.PublicKey.Curve = curve
	ra.D = new(big.Int).SetBytes(raBuf)
	ra.PublicKey.X, ra.PublicKey.Y = curve.ScalarBaseMult(raBuf)

	rb := new(sm2.PrivateKey)
	rb.PublicKey.Curve = curve
	rb.D = new(big.Int).SetBytes(rbBuf)
	rb.PublicKey.X, rb.PublicKey.Y = curve.ScalarBaseMult(rbBuf)

	k1, Sb, S2, err := sm2.KeyExchangeB(16, ida, idb, db, &da.PublicKey, rb, &ra.PublicKey)
	if err != nil {
		panic(err)
	}
	k2, S1, Sa, err := sm2.KeyExchangeA(16, ida, idb, da, &db.PublicKey, ra, &rb.PublicKey)
	if err != nil {
		panic(err)
	}
	if bytes.Compare(k1, k2) != 0 {
		panic("key exchange differ")
	}
	//if bytes.Compare(k1, expk) != 0 {
	//	fmt.Printf("expected %x, found %x\n", expk, k1)
	//	panic("")
	//}
	if bytes.Compare(S1, Sb) != 0 {
		panic("hash verfication failed")
	}
	if bytes.Compare(Sa, S2) != 0 {
		panic("hash verfication failed")
	}

	fmt.Println("pass")
}
