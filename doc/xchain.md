## 基于国密算法在区块链上进行授权数据交换



### 1. 概述

​		区块链具有分布式、去中心化、防篡改等特性。由于区块链构建的特性，每个节点都存储有全部的区块数据，理论上用户可以查看到所有的区块数据。当使用区块链对敏感数据进行数据交换或数据共享时，需要对用户查看数据的权限进行验证和授权，用户只能查看或共享被授权的数据。同时，因为区块链特有的放篡改性质，将数据授权和数据交换的过程记录在区块链上，可以实现对用户授权、数据交换、数据共享等操作的留痕和追溯。因此，如何在区块链上构建数据授权和交换的方法是具体应用区块链时必不可少的需求。
​		国密算法是我国自主研发创新的一套数据加密处理系列算法。从SM1-SM4分别实现了对称、非对称、摘要等算法功能。特别适合应用于嵌入式物联网等相关领域，完成身份认证和数据加解密等功能。2010年底，国家密码管理局公布了我国自主研制的“椭圆曲线公钥密码算法”（SM2算法）。为保障重要经济系统密码应用安全，国家密码管理局于2011年发布了《关于做好公钥密码算法升级工作的通知》，要求“自2011年3月1日起，在建和拟建公钥密码基础设施电子认证系统和密钥管理系统应使用国密算法。自2011年7月1日起，投入运行并使用公钥密码的信息系统，应使用SM2算法。”
​		现有典型区块链应用实现（例如以太坊、超级账本Fabric）比较完备，但是过于复杂，例如以太坊基于工作量证明实现共识（挖矿），交易吞吐量受限；Fabric基于PBFT共识算法，但是实现过于复杂，所有Fabric应用均需要基于docker容器（称为“链码”的模块），增加节点运维成本。
​		因此，在本次实现中，使用Tendermint作为区块链共识引擎实现轻量化区块链应用，便于维护。去中心化、分布式存储链上信息；使用LevelDB嵌入式数据库存储链表索引数据；使用国密算法进行密钥交换和数据加密。以授权数据交换和共享为场景，实现了：（1）用户数据链上存储，实现防篡改和可追溯；（2）使用SM2非对称密钥作为用户密钥进行用户识别和授权密钥交换；（3）使用SM4对称加密算法进行数据加密；（4）节点分布式去中心化部署，允许1/3节点故障仍可以正常工作。



### 2. 设计概要



### 3. 数据结构设计

区块交易结构基类

```go
// Transx 事务基类
type Transx struct {
	Signature  []byte //发送方对这个消息的私钥签名
	SendTime   *time.Time
	SignPubKey []byte // 签名公钥
	Payload    IPayload
}
```

交易数据payload

```go
// 交易信息
type Deal struct {
	ID     uuid.UUID // 交易ID
	UserID []byte //用户的加密公钥
	Data   []byte // 加密交易数据（例如 ipfs hash）
}
```

授权请求/授权操作payload

```go
// 授权操作
// ToUserID 请求 FromUserID 授权，指定 DealID
// FromUserID 加密数据 Data 后 返回 ToUserID
type Auth struct {
	ID          uuid.UUID
	ReqID       uuid.UUID // 授权请求的ID， action==5时使用
	DealID      uuid.UUID // 相关交易ID
	FromUserID  []byte // 数据所有者的用户ID（用户公钥）
	ToUserID    []byte // 被授权的用户ID（用户公钥）
	Data        []byte //  action==4, rb.pub (33 bytes) + rb.priv (B的私钥加密的)
			           //  action==5, ra.pub (33 butes) + data (协商的密钥加密)
	Action      byte // 0x04 请求授权， 0x05 响应授权
}
```



### 4. 使用示例

#### 4.1 本地节点测试

编译

```shell
$ make build
```


初始化

```shell
$ build/xchain init --home n1
```


启动节点

```shell
$ build/xchain node --home n1 --consensus.create_empty_blocks=false
```

初始化用户密钥
```shell
$ build/xcli init --home users/user1
$ build/xcli init --home users/user2
$ cat users/user1/user.key 
{"PrivKey":"UDD5X7pNUMgQs1XXxiqj91yteZkmcrQuiIux5RTUu90=","PubKey":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG"}
$ cat users/user2/user.key 
{"PrivKey":"8afZQKZXej31zU5CDpQZVFoi6+x///59ZuKh9wyJ2Ag=","PubKey":"A/45NkcCX4WIw+En0eQPgbp2oSaYbCAOiYqTfQTbZgih"}
```



#### 4.2 数据上链

##### user1 数据上链

> 返回值为区块id 

```shell
$ build/xcli deal --home users/user1 "hello world"
Deal ==> {"id":"ddd8c8c2-d625-46ad-ab6c-6c49aad45836"}
```



##### user1 查询链上数据

> 查询上述新添加的数据区块，因为user1是提交者，所以可以解密data数据，看到明文。

```shell
$ build/xcli queryTx --home users/user1 _ ddd8c8c2-d625-46ad-ab6c-6c49aad45836
Tx ==> {"data":"hello world","id":"ddd8c8c2-d625-46ad-ab6c-6c49aad45836","send_time":"2021-01-07T06:17:08.356336674Z","type":"DEAL","user_id":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG"}
```

> 下面是区块上的原始数据，Data字段是SM4加密存储的，使用user1的私钥进行加密。

```json
{
    "Signature":"MEUCIQCyDuZfg6aUGbFdaonXC81GLynPwoXQkGUbaX6rhdW02AIgcGlDVAhH0hh86CcXo3zrntC+n2NDMFNHuXIFAkG2TN0=",
    "SendTime":"2021-01-07T06:17:08.356336674Z",
    "SignPubKey":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG",
    "Payload":{
        "type":"deal",
        "value":{
            "ID":"3djIwtYlRq2rbGxJqtRYNg==",
            "UserID":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG",
            "Data":"NMiY8YnJVX7UUhLHTTgzsg=="
        }
    }
}
```



##### user2 查询链上数据

> 因为此区块不是user2的，所以当user2查询时，只能看到加密数据，不能解密。

```shell
$ build/xcli queryTx --home users/user2 "A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG" ddd8c8c2-d625-46ad-ab6c-6c49aad45836
Tx ==> {"data":"NMiY8YnJVX7UUhLHTTgzsg==","id":"ddd8c8c2-d625-46ad-ab6c-6c49aad45836","send_time":"2021-01-07T06:17:08.356336674Z","type":"DEAL","user_id":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG"}
```



#### 4.3 授权数据交换

##### user2 提出数据授权请求

> user2向user1发起请求，获取user1的数据。第一个参数是数据所有这的用户id，第二个参数是区块id。
> user1 公钥： ```A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG```
>
> 返回值为授权请求的区块id，授权请求也会上链，用于追溯。

```shell
$ build/xcli authReq --home users/user2 "A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG" ddd8c8c2-d625-46ad-ab6c-6c49aad45836
AuthReq ==> {"id":"a0dfc291-b70e-45d5-a360-67c3099ee775"}
```



##### user1 查询授权请求

> user1 可以查询谁提交的授权请求（action==4）
>
> from_user_id 是接收授权请求的用户，即数据的所有者，这里是 user1
>
> to_user_id 是发起授权请求的用户，这里是 user2  
>
> deal_id 是 to_user_id 请求授权的 区块id
>
> user2的公钥：```A/45NkcCX4WIw+En0eQPgbp2oSaYbCAOiYqTfQTbZgih```

```shell
$ build/xcli queryAuth --home users/user1
Auth ==> [{"action":4,"data":"","deal_id":"ddd8c8c2-d625-46ad-ab6c-6c49aad45836","from_user_id":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG","id":"a0dfc291-b70e-45d5-a360-67c3099ee775","send_time":"2021-01-07T06:19:52.274531072Z","to_user_id":"A/45NkcCX4WIw+En0eQPgbp2oSaYbCAOiYqTfQTbZgih","type":"AUTH"}]
```

> 下面是授权请求在区块上的原始数据（action==4），Data字段是用SM4加密存储的，使用user1和user2协商的密钥进行加密。此时Data字段存储用于密钥协商的密钥数据。

```json
{
    "Signature":"MEUCIQDx2eV34dEkXgjuMDja8UB14qJ2DkVsnFWrjePEF7kaNQIgWqrtGlrj03NY4cONYHSHCjNPOC7VF6IXTHbnFWdasdE=",
    "SendTime":"2021-01-07T06:19:52.274531072Z",
    "SignPubKey":"A/45NkcCX4WIw+En0eQPgbp2oSaYbCAOiYqTfQTbZgih",
    "Payload":{
        "type":"auth",
        "value":{
            "ID":"oN/CkbcORdWjYGfDCZ7ndQ==",
            "ReqID":"AAAAAAAAAAAAAAAAAAAAAA==",
            "DealID":"3djIwtYlRq2rbGxJqtRYNg==",
            "FromUserID":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG",
            "ToUserID":"A/45NkcCX4WIw+En0eQPgbp2oSaYbCAOiYqTfQTbZgih",
            "Data":"IQNMyfoDgrAH6XpPFMB+xhXsPx5cj3eir/JbWv9nbez9rcmvBjByL1Q1Y0CeV2kxFDpkN3lIg3v2/zWlKAnA2rT6h3LzD1NgEse1KFgK4h2qrA==",
            "Action":4
        }
    }
}
```



##### user1 对授权请求进行授权

> 参数为授权请求的区块id
>
> 返回值为响应授权请求的区块id，授权操作也会上链，用于追溯。

```shell
$ build/xcli authResp --home users/user1 a0dfc291-b70e-45d5-a360-67c3099ee775
AuthResp ==> {"id":"bea9a9d5-9cde-4dec-9ac9-8a23c2ca4903"}
```



##### user2 查询授权结果

> 当user1授权后，user2就可以查询到授权的记录，并可以加密data字段，获得明文。

```shell
$ build/xcli queryAuth --home users/user2
Auth ==> [{"action":5,"data":"hello world","deal_id":"ddd8c8c2-d625-46ad-ab6c-6c49aad45836","from_user_id":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG","id":"bea9a9d5-9cde-4dec-9ac9-8a23c2ca4903","req_id":"a0dfc291-b70e-45d5-a360-67c3099ee775","send_time":"2021-01-07T06:21:20.728255461Z","to_user_id":"A/45NkcCX4WIw+En0eQPgbp2oSaYbCAOiYqTfQTbZgih","type":"AUTH"}]
```

> 下面是授权操作在区块上的原始数据（action==5），Data字段是用SM4加密存储的，使用user1和user2协商的密钥进行加密。

```json
{
    "Signature":"MEQCIFQmtsprc9PvgzmHBSNRREF8emMBav5/N1s3o67KPQ10AiAt+q5tN8oH50S2uAGL0rQKXfzYdbRCYfObzAnUSIbguQ==",
    "SendTime":"2021-01-07T06:21:20.728255461Z",
    "SignPubKey":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG",
    "Payload":{
        "type":"auth",
        "value":{
            "ID":"vqmp1ZzeTeyayYojwspJAw==",
            "ReqID":"oN/CkbcORdWjYGfDCZ7ndQ==",
            "DealID":"3djIwtYlRq2rbGxJqtRYNg==",
            "FromUserID":"A2FCWvU0EUuqhZKL1KRRaIcxKNx/8HUw1Uz8ZfH/qEMG",
            "ToUserID":"A/45NkcCX4WIw+En0eQPgbp2oSaYbCAOiYqTfQTbZgih",
            "Data":
            "IQKrWsz5EOZAMLJUe1Dq6lzwPrp8ECoAU19IODmAhz3xVL+mWmJdkS6r5bsfmDj+BYc=",
            "Action":5
        }
    }
}
```



### 5. 代码工程

https://gitlab.ylzpay.com/guantao/xchain



### 6. 技术栈	

- 区块链中间件：Tendermint 0.34.0（https://github.com/tendermint/tendermint）

- 数据库：LevelDB 1.20（https://github.com/google/leveldb） 

- 开发语言：Go 1.15.6

