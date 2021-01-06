## Xchain 
在区块链上共享、授权加密数据并进行追溯


### 功能描述

1. 链上存储数据
2. 数据实体为字节流，内容由应用确定，不做限制
3. 交易数据、授权数据全部上链
4. 以链用户区分数据所有权，链用户提交的数据对自己是开放的
5. 交易数据使用链用户密钥加密，其他用户是否可见，需要获取数据所有者的授权
6. 交易节点提供功能：数据提交，交易查询，查询授权请求，授权查询
7. 交易区块的类型：交易数据；查询授权
8. 密钥使用SM2算法，加密算法使用SM4
9. 密钥交换过程依照《GMT00003.3-2012 SM2椭圆曲线公钥密码算法 第3部分：密钥交换协议》


### 交易链请求

1. 交易内容
```json
{
	"Signature":"SXKaVqfAe5ypHpz1qM3tQTYa42F9JQoq4zwMYQLKN0E0s+nViVk2Z3b98mFXvTHnqRCFousPVCYdR7b+d21jCg==",
	"SendTime":"2020-12-18T05:36:00.281914675Z",
	"SignPubKey":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG",
	"Payload":{
		// 详见下
	}
}
```

2. 资产交易 payload
```json
{
	"type":"deal", // 交易
	"value":{
		"ID":"", // 资产交易ID
		"User":"P1ABCkAph4DQlnEMahGW6I2mfOOtfZKYyssOZ4L8MTc=", // 用户ID（公钥）
		"Data":"", // 加密交易数据 （IPFS HASH）
	}
}
```

3. 查询授权 payload
```json
{
	"type":"auth", // 查询授权（授权其他用户查看某资产），查询记录（只记录被授权方的查询动作）
	"value":{
		"ID":"", // 授权操作ID
		"DealID":"", // 交易ID
		"FromUser":"P1ABCkAph4DQlnEMahGW6I2mfOOtfZKYyssOZ4L8MTc=", // 授权用户ID（公钥）
		"ToUser":"P1ABCkAph4DQlnEMahGW6I2mfOOtfZKYyssOZ4L8MTc=", // 被授权用户ID（公钥）
		"Data":"", // FromUser加密数据，被授权者ToUserID可以解密
		"Action":4, // 0x04 请求授权， 0x05 响应授权
	}
}
```



### 区块例子

```json
{
	"header":{
		"version":{
			"block":"11",
			"app":"1"
		},
		"chain_id":"test-chain-FEeTGF",
		"height":"8",
		"time":"2020-12-24T05:24:01.760181367Z",
		"last_block_id":{
			"hash":"3D326CA03E1D0E6D9C80FB6B788AD1A72BB12E10B9DE617B13AF311E5258ABA8",
			"parts":{
				"total":1,
				"hash":"0A675856141DB019BB245E87896DEA5A3C7BAE8CC2C3C1A7666DF96236B802CB"
			}
		},
		"last_commit_hash":"63AA80A1CE7261E001383A3754DAF09C52DA25881EFDD6E3E1F1541A937C5AAE",
		"data_hash":"D14D445C897F2E7A2518FB1EEE69A969F178F6D819ADEBB8D91B5B528CCC01C7",
		"validators_hash":"82F872A1F21F7C05578D5397DA499A2D656E61D0DD7F9EFE6531F433F72306EB",
		"next_validators_hash":"82F872A1F21F7C05578D5397DA499A2D656E61D0DD7F9EFE6531F433F72306EB",
		"consensus_hash":"048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
		"app_hash":"0000000000000000",
		"last_results_hash":"6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D",
		"evidence_hash":"E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
		"proposer_address":"89581502243C3D3401C38EEF4C1A1145AB514B11"
	},
	"data":{
		"txs":[
			"eyJTaWduYXR1cmUiOiIwbDFvTFhoayt4YXF2UC96MXJMZ1U5Mzh3d01wZzExYlMwcko3V1lrNmNlNzg1Z ..."
		]
	},
	"evidence":{
		"evidence":[]
	},
	"last_commit":{
		"height":"7",
		"round":0,
		"block_id":{
			"hash":"3D326CA03E1D0E6D9C80FB6B788AD1A72BB12E10B9DE617B13AF311E5258ABA8",
			"parts":{
				"total":1,
				"hash":"0A675856141DB019BB245E87896DEA5A3C7BAE8CC2C3C1A7666DF96236B802CB"
			}
		},
		"signatures":[
			{
				"block_id_flag":2,
				"validator_address":"89581502243C3D3401C38EEF4C1A1145AB514B11",
				"timestamp":"2020-12-24T05:24:01.760181367Z",
				"signature":"ZfgneVPY/pOEjygwmEQnMIu4iQT8QgRf/AdjHptkbpqT57dCMFa4V+7bxAKIzoCUcJgtLtrg1bJtJdXQNk7gCA=="
			}
		]
	}
}
```



### leveldb 逻辑分表

| 前缀       | key             | value     |
| ---------- | --------------- | --------- |
| blockLink: | 区块高度        | 区块高度  |
| userLink:  | 用户id        | 区块高度  |



### 技术栈

1. 区块链 Tendermint 0.34.0
2. 节点数据库 LevelDB 1.20
3. 开发语言 Go 1.15.6
