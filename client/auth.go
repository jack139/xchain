package client

import (
	"xchain/types"

	"fmt"
	"bytes"
	"time"
	"crypto/rand"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
)


// 请求授权 上链
// xcli authRequest j9cIgmm17x0aLApf0i20UR7Pj34Ua/JwyWOuBGgYIFg= dcfe656c-6c65-45e7-9e94-f082a068a93d
func (me *User) AuthRequest(fromUserId, dealId string) ([]byte, error) {
	now := time.Now()

	// 检查 toUserId 合理性
	var fromUserIdBytes []byte
	err := cdc.UnmarshalJSON([]byte("\""+fromUserId+"\""), &fromUserIdBytes) // 反序列化时需要双引号，因为是字符串
	if err != nil {
		return nil, err
	}

	// 检查 dealID -->  UUID
	uuidDealId, err := uuid.FromString(dealId)
	if err != nil {
		return nil, err
	}


	// 生成 密钥交换的数据

	// 生成 rB
	rb, _ := sm2.GenerateKey(rand.Reader) // 生成密钥对
	rbPubBytes := sm2.Compress(&rb.PublicKey) // 33 bytes

	// 用sm2加密rb私钥
	encrypted, err := sm2.EncryptAsn1(&me.SignKey.PublicKey, rb.D.Bytes(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("sm2 encrypt fail: %s", err)
	}

	// data格式： rb.pub长度(byte) + rb.pub(33bytes?) + sm2加密的rb.priv
	cryptData := append([]byte{byte(len(rbPubBytes))}, rbPubBytes...)
	cryptData = append(cryptData, encrypted...)

	// 新建交易
	tx := new(types.Transx)
	tx.SendTime = &now

	auth := new(types.Auth)
	auth.ID = uuid.NewV4()
	auth.DealID = uuidDealId
	auth.FromUserID = fromUserIdBytes
	auth.ToUserID = *me.CryptoPair.PubKey
	auth.Data = cryptData
	auth.Action = 0x04

	tx.Payload = auth

	tx.Sign(me.SignKey)

	bz, err := cdc.MarshalJSON(&tx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ret, err := cli.BroadcastTxSync(ctx, bz)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("auth request => %+v\n", ret)

	// ret  *ctypes.ResultBroadcastTxCommit
	if ret.Code !=0 {
		fmt.Println(ret.Log)
		return nil, fmt.Errorf(ret.Log)
	}

	respMap := map[string]interface{}{"id" : auth.ID.String()}

	// 返回结果转为json
	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}


// 响应授权 上链
// xcli authRequest dcfe656c-6c65-45e7-9e94-f082a068a93d
func (me *User) AuthResponse(authId string) ([]byte, error) {
	addr, _ := cdc.MarshalJSON(*me.CryptoPair.PubKey)

	now := time.Now()

	// 获取 authID 对应的 授权请求 块
	authTx, err := queryTx(addr, "_", authId)
	if err != nil {
		return nil, err
	}
	if authTx==nil {
		return nil, fmt.Errorf("AuthID not found")
	}
	auth, ok := (*authTx).Payload.(*types.Auth)	// 授权块
	if !ok {
		return nil, fmt.Errorf("need a Auth Payload")
	}

	// 检查fromUserId 是否是自己, 说明有问题！ 可能被黑！
	if bytes.Compare(auth.FromUserID, *me.CryptoPair.PubKey)!=0 {
		return nil, fmt.Errorf("---> NOT MY AUTH <---")
	}

	// 检查是否已响应过，在toUserID的列表里找
	toUserId, _ := cdc.MarshalJSON(auth.ToUserID)
	isAuthorised, err := checkAuthResp(addr, string(toUserId), authId)
	if err != nil {
		return nil, err
	}
	if isAuthorised { // 已经授权过
		return nil, fmt.Errorf("Authorized")
	}

	// 获取 authID 对应的 交易块
	dealTx, err := queryTx(addr, "_", auth.DealID.String())
	if err != nil {
		return nil, err
	}
	if dealTx==nil {
		return nil, fmt.Errorf("DealID not found")
	}

	deal, ok := (*dealTx).Payload.(*types.Deal)	// 交易块
	if !ok {
		return nil, fmt.Errorf("need a Deal Payload")
	}

	// 解密

	// 从data中解密出sm4密钥
	// data格式： sm4密钥长度(byte)(sm2加密) + sm4加密的data
	sm4keyLen := int(deal.Data[0])

	decryptKey, err := sm2.DecryptAsn1(&me.SignKey, deal.Data[1:sm4keyLen+1])
	if err != nil {
		return nil, fmt.Errorf("sm2 decrypt fail: %s", err)
	}

	// 加密密钥使用私钥前16字节（128bit）
	//decryptKey := (*me.CryptoPair.PrivKey)[:16]

	decrypted, err := sm4.Sm4CFB(decryptKey, deal.Data[sm4keyLen+1:], false)
	if err!=nil {
		return nil, fmt.Errorf("sm4 decrypt error: %s", err)
	}
	//fmt.Printf("plain --> %s\n", decrypted)


	// 从 auth请求 里，获取密钥交换的数据
	rbBytesLen := int(auth.Data[0])
	dbBytes := auth.ToUserID
	rbBytes := auth.Data[1:rbBytesLen+1]

	dbPub := restorePublicKey(dbBytes)
	rbPub := restorePublicKey(rbBytes)

	// da 就是自己的密钥
	daPriv := me.SignKey
	// 生成 ra
	raPriv, _ := sm2.GenerateKey(rand.Reader) // 生成密钥对

	//  生成 key
	encryptKey, _, _, err := sm2.KeyExchangeA(16, 
		auth.FromUserID, auth.ToUserID, &daPriv, dbPub, raPriv, rbPub)
	if err != nil {
		return nil, err
	}

	// 重新加密
	encrypted, err := sm4.Sm4CFB(encryptKey, decrypted, true)
	if err != nil {
		return nil, fmt.Errorf("sm4 encrypt error: %s", err)
	}
	//fmt.Printf("%d %v\n", len(encrypted), encrypted)


	raPubBytes := sm2.Compress(&raPriv.PublicKey) // 33 bytes

	// data格式： ra.pub长度(byte) + ra.pub(33bytes?) + K加密的数据
	cryptData := append([]byte{byte(len(raPubBytes))}, raPubBytes...)
	cryptData = append(cryptData, encrypted...)

	// 新建交易
	tx := new(types.Transx)
	tx.SendTime = &now

	authResp := new(types.Auth)
	authResp.ID = uuid.NewV4()
	authResp.ReqID = auth.ID
	authResp.DealID = auth.DealID
	authResp.FromUserID = auth.FromUserID
	authResp.ToUserID = auth.ToUserID
	authResp.Data = cryptData
	authResp.Action = 0x05

	tx.Payload = authResp

	tx.Sign(me.SignKey)

	bz, err := cdc.MarshalJSON(&tx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ret, err := cli.BroadcastTxSync(ctx, bz)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("auth respose => %+v\n", ret)

	// ret  *ctypes.ResultBroadcastTxCommit
	if ret.Code !=0 {
		fmt.Println(ret.Log)
		return nil, fmt.Errorf(ret.Log)
	}

	respMap := map[string]interface{}{"id" : authResp.ID.String()}

	// 返回结果转为json
	respBytes, err := json.Marshal(respMap)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
